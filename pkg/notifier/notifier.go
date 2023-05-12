package notifier

import (
	"context"

	"github.com/mdwn/ghstatus/pkg/ghstatus"
)

// Notifier is used by the monitor to pass the given message to another system.
type Notifier interface {
	// Name is the name of the notifier.
	Name() string

	// Notify will notify an underlying system with the given message.
	Notify(context.Context, Message) error

	// Cleanup performs any cleanup steps.
	Cleanup() error
}

// Message is a notification message.
type Message struct {
	// ChangedStatus is populated if the status has changed.
	ChangedStatus *ghstatus.Status

	// ChangedComponents is populated if the components have changed.
	ChangedComponents []ghstatus.Component

	// ChangedIncidents is populated if the incidents have changed.
	ChangedIncidents []ghstatus.Incident

	// ChangedScheduledMaintenances is populated of the scheduled maintenances have changed.
	ChangedScheduledMaintenances []ghstatus.ScheduledMaintenance
}
