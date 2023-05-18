package notifiers

import (
	"context"
	"fmt"
	"io"

	"github.com/mdwn/ghstatus/pkg/notifier"
)

const (
	writer = "writer-notifier"
)

// WriterNotifier writes the output to the given io.Writer. This is
// meant to be used by other notifiers and not directly, so it is
// not registered with the notifier registry.
type WriterNotifier struct {
	writer io.Writer
}

// NewWriterNotifier will return a writer notifier.
func NewWriterNotifier(writer io.Writer) *WriterNotifier {
	return &WriterNotifier{
		writer: writer,
	}
}

// Name is the name of the notifier.
func (*WriterNotifier) Name() string {
	return writer
}

// Notify will notify an underlying system with the given message.
func (w *WriterNotifier) Notify(_ context.Context, msg notifier.Message) error {
	if msg.ChangedStatus != nil {
		_, err := fmt.Fprintf(w.writer, "Status: %s (%s)\n", msg.ChangedStatus.Indicator, msg.ChangedStatus.Description)
		if err != nil {
			return fmt.Errorf("error while writing status: %w", err)
		}
	}

	if len(msg.ChangedComponents) > 0 {
		for _, component := range msg.ChangedComponents {
			_, err := fmt.Fprintf(w.writer, "Component %s: %s, updated at: %s\n", component.Name, component.Status, component.UpdatedAt)
			if err != nil {
				return fmt.Errorf("error while writing component: %w", err)
			}
		}
	}

	if len(msg.ChangedIncidents) > 0 {
		for _, incident := range msg.ChangedIncidents {
			lastUpdate := ""
			if len(incident.IncidentUpdates) > 0 {
				lastUpdate = incident.IncidentUpdates[0].Body
			}
			_, err := fmt.Fprintf(w.writer, "Incident %s: %s, updated at: %s%s\n", incident.Name, incident.Status, incident.UpdatedAt, lastUpdate)
			if err != nil {
				return fmt.Errorf("error while writing status: %w", err)
			}
		}
	}

	if len(msg.ChangedScheduledMaintenances) > 0 {
		for _, scheduledMaintenance := range msg.ChangedScheduledMaintenances {
			_, err := fmt.Fprintf(w.writer, "Scheduled maintenance %s: %s, updated at: %s\n",
				scheduledMaintenance.Name, scheduledMaintenance.Status, scheduledMaintenance.UpdatedAt)
			if err != nil {
				return fmt.Errorf("error while writing status: %w", err)
			}
		}
	}

	return nil
}

// Cleanup performs any cleanup steps.
func (w *WriterNotifier) Cleanup() error { return nil }
