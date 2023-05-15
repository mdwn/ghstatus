package monitor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"github.com/mdwn/ghstatus/pkg/logging"
	"github.com/mdwn/ghstatus/pkg/notifier"
	"go.uber.org/zap"
)

const (
	// this is a faux-component that shows up in the Github API. We'll filter it out when displaying changes to the user.
	fauxComponentName = "Visit www.githubstatus.com for more information"
)

// Monitor will periodically poll the Github Status and issues updates
// to the given callback.
type Monitor struct {
	log              *zap.Logger
	clock            clockwork.Clock
	client           ghstatus.Client
	notifyOnFirstRun bool

	notifiersMu sync.RWMutex
	notifiers   map[string]notifier.Notifier
}

// New creates a new Github Status monitor.
func New(log *zap.Logger, clock clockwork.Clock, client ghstatus.Client, notifyOnFirstRun bool) *Monitor {
	return &Monitor{
		log:              logging.WithComponent(log, "monitor"),
		clock:            clock,
		client:           client,
		notifyOnFirstRun: notifyOnFirstRun,
		notifiers:        map[string]notifier.Notifier{},
	}
}

// RegisterNotifier will register a notifier with the monitor.
func (m *Monitor) RegisterNotifier(notifier notifier.Notifier) error {
	m.notifiersMu.Lock()
	defer m.notifiersMu.Unlock()

	if _, ok := m.notifiers[notifier.Name()]; ok {
		return fmt.Errorf("duplicate notifier %s", notifier.Name())
	}

	m.notifiers[notifier.Name()] = notifier
	return nil
}

// MonitorAndNotify will monitor the Github Status and notify subscribers upon relevant changes.
func (m *Monitor) MonitorAndNotify(ctx context.Context, timeBetweenPolls time.Duration) {
	ticker := m.clock.NewTicker(timeBetweenPolls)
	defer ticker.Stop()

	var lastSummary ghstatus.SummaryResponse
	var err error

	for {
		lastSummary, err = m.detectChangesAndNotify(ctx, lastSummary)
		if err != nil {
			m.log.With(zap.Error(err)).Error("error during monitoring")
		}

		select {
		case <-ticker.Chan():
		case <-ctx.Done():
			return
		}
	}
}

// detectChangesAndNotify will detect any changes and send notifications based on the differences.
func (m *Monitor) detectChangesAndNotify(ctx context.Context, lastSummary ghstatus.SummaryResponse) (ghstatus.SummaryResponse, error) {
	summary, err := m.client.Summary(ctx)
	if err != nil {
		return lastSummary, fmt.Errorf("error getting summary: %w", err)
	}

	// Skip the first notification if notifyOnFirstRun is disabled.
	if lastSummary.Page.UpdatedAt.IsZero() && !m.notifyOnFirstRun {
		m.log.Debug("Notify on first run is disabled, skipping the notification.")
		return summary, nil
	}

	// If the summary page hasn't updated, no need to continue.
	if summary.Page.UpdatedAt.Equal(lastSummary.Page.UpdatedAt) {
		m.log.Debug("Current summary is equal to the old one, no updates.")
		return lastSummary, err
	}

	// Check to see if the status has changed.
	var changedStatus *ghstatus.Status
	if lastSummary.Status.Description != summary.Status.Description || lastSummary.Status.Indicator != summary.Status.Indicator {
		changedStatus = &summary.Status
	}

	changedComponents := findChangedComponents(lastSummary.Components, summary.Components)
	changedIncidents := findChangedIncidents(lastSummary.Incidents, summary.Incidents)
	changedScheduledMaintenances := findChangedScheduledMaintenances(lastSummary.ScheduledMaintenances, summary.ScheduledMaintenances)

	if changedStatus != nil || len(changedComponents) > 0 || len(changedIncidents) > 0 || len(changedScheduledMaintenances) > 0 {
		m.log.Debug("A change was found, running through the notifiers.")
		var errs []error
		notifierMsg := notifier.Message{
			ChangedStatus:                changedStatus,
			ChangedComponents:            changedComponents,
			ChangedIncidents:             changedIncidents,
			ChangedScheduledMaintenances: changedScheduledMaintenances,
		}
		for _, notifier := range m.notifiers {
			if err := notifier.Notify(ctx, notifierMsg); err != nil {
				errs = append(errs, err)
			}
		}
		return summary, errors.Join(errs...)
	}

	return summary, nil
}

// these functions will be used for finding changed resources generically.
func getComponentID(component ghstatus.Component) string           { return component.Name }
func getComponentUpdatedAt(component ghstatus.Component) time.Time { return component.UpdatedAt }
func getIncidentID(incident ghstatus.Incident) string              { return incident.ID }
func getIncidentUpdatedAt(incident ghstatus.Incident) time.Time    { return incident.UpdatedAt }
func getScheduledMaintenanceID(scheduledMaintenance ghstatus.ScheduledMaintenance) string {
	return scheduledMaintenance.ID
}
func getScheduledMaintenanceUpdatedAt(scheduledMaintenance ghstatus.ScheduledMaintenance) time.Time {
	return scheduledMaintenance.UpdatedAt
}

// idGetter will get the ID from a resource.
type idGetter[T any] func(T) string

// updatedAtGetter will get the updated at time from a resource.
type updatedAtGetter[T any] func(T) time.Time

// findChangedComponents will return any components which have changed from the last known state.
func findChangedComponents(last []ghstatus.Component, current []ghstatus.Component) []ghstatus.Component {
	return findChangedResources(last, current, getComponentID, getComponentUpdatedAt)
}

// findChangedIncidents will return any incidents which have changed from the last known state.
func findChangedIncidents(last []ghstatus.Incident, current []ghstatus.Incident) []ghstatus.Incident {
	return findChangedResources(last, current, getIncidentID, getIncidentUpdatedAt)
}

// findChangedScheduledMaintenances will return any scheduled maintenances which have changed from the last known state.
func findChangedScheduledMaintenances(last []ghstatus.ScheduledMaintenance, current []ghstatus.ScheduledMaintenance) []ghstatus.ScheduledMaintenance {
	return findChangedResources(last, current, getScheduledMaintenanceID, getScheduledMaintenanceUpdatedAt)
}

// findChangedResources will return any resources which have changed from the last known state.
func findChangedResources[T any](last []T, current []T, idGetter idGetter[T], updatedAtGetter updatedAtGetter[T]) []T {
	lastMap := map[string]T{}
	for _, resource := range last {
		lastMap[idGetter(resource)] = resource
	}

	var changedResources []T
	// Compare the current list of resources to see if any of them are updated.
	// We'll intentionally disregard any disappearing resources as it likely means the
	// resource has been cleaned up or is resolved.
	for _, resource := range current {
		resourceID := idGetter(resource)

		// Disregard the faux component if that's what we're looking at.
		if resourceID == fauxComponentName {
			continue
		}

		// If this resource isn't present in the last map, then this is new.
		lastResource, ok := lastMap[resourceID]
		if !ok {
			changedResources = append(changedResources, resource)
			continue
		}

		// If the updated field has changed, then the resource has changed.
		if !updatedAtGetter(lastResource).Equal(updatedAtGetter(resource)) {
			changedResources = append(changedResources, resource)
		}
	}

	return changedResources
}
