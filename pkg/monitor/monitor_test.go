package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"github.com/mdwn/ghstatus/pkg/notifier"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// channelNotifier is for testing and notifies a channel with the notification message.
type channelNotifier struct {
	ch chan notifier.Message
}

func (c *channelNotifier) Name() string { return "channel" }
func (c *channelNotifier) Cleanup() error {
	close(c.ch)
	return nil
}
func (c *channelNotifier) Notify(_ context.Context, msg notifier.Message) error {
	c.ch <- msg
	return nil
}

var _ notifier.Notifier = &channelNotifier{}

func TestMonitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	log := zap.NewNop()
	clock := clockwork.NewFakeClock()

	server, client := ghstatus.NewTestServerAndClient(t)
	m := New(log, clock, client, true)

	// Make this channel for the dummy notifier.
	ch := make(chan notifier.Message, 1)
	require.NoError(t, m.RegisterNotifier(&channelNotifier{
		ch: ch,
	}))

	// Start with an empty response.
	server.SetSummary(t, ghstatus.SummaryResponse{})

	go m.MonitorAndNotify(ctx, time.Minute)

	status := ghstatus.Status{
		Indicator:   ghstatus.Major,
		Description: "something happened",
	}

	// The first response should be ignored. Let's update the summary with a new status.
	server.SetSummary(t, ghstatus.SummaryResponse{
		Page: ghstatus.Page{
			UpdatedAt: clock.Now().UTC(),
		},
		Status: status,
	})
	clock.Advance(time.Minute)

	msg := waitForNotification(t, ch)

	require.Equal(t, notifier.Message{ChangedStatus: &status}, msg)

	// Let's add in a component.
	component := ghstatus.Component{
		Name:      "component",
		UpdatedAt: clock.Now().UTC(),
	}
	server.SetSummary(t, ghstatus.SummaryResponse{
		Page: ghstatus.Page{
			UpdatedAt: clock.Now().UTC(),
		},
		Status:     status,
		Components: []ghstatus.Component{component},
	})

	clock.Advance(time.Minute)

	msg = waitForNotification(t, ch)

	require.Equal(t, notifier.Message{ChangedComponents: []ghstatus.Component{component}}, msg)

	// Let's update the component.
	component = ghstatus.Component{
		Name:      "component",
		UpdatedAt: clock.Now().UTC(),
		Status:    ghstatus.DegradedPerformance,
	}
	server.SetSummary(t, ghstatus.SummaryResponse{
		Page: ghstatus.Page{
			UpdatedAt: clock.Now().UTC(),
		},
		Status:     status,
		Components: []ghstatus.Component{component},
	})

	clock.Advance(time.Minute)

	msg = waitForNotification(t, ch)

	require.Equal(t, notifier.Message{ChangedComponents: []ghstatus.Component{component}}, msg)

	// Let's keep the everything the same and update with a new incident.
	incident1 := ghstatus.Incident{
		ID:        "Incident 1",
		UpdatedAt: clock.Now().UTC(),
		IncidentUpdates: []ghstatus.IncidentUpdate{
			{
				UpdatedAt: clock.Now().UTC(),
				Body:      "Update 1",
			},
		},
	}
	server.SetSummary(t, ghstatus.SummaryResponse{
		Page: ghstatus.Page{
			UpdatedAt: clock.Now().UTC(),
		},
		Status:     status,
		Components: []ghstatus.Component{component},
		Incidents: []ghstatus.Incident{
			incident1,
		},
	})
	clock.Advance(time.Minute)

	msg = waitForNotification(t, ch)

	require.Equal(t, notifier.Message{
		ChangedIncidents: []ghstatus.Incident{incident1},
	}, msg)

	// Let's add another new incident
	incident2 := ghstatus.Incident{
		ID:        "Incident 2",
		UpdatedAt: clock.Now().UTC(),
		IncidentUpdates: []ghstatus.IncidentUpdate{
			{
				UpdatedAt: clock.Now().UTC(),
				Body:      "Update 2",
			},
		},
	}
	server.SetSummary(t, ghstatus.SummaryResponse{
		Page: ghstatus.Page{
			UpdatedAt: clock.Now().UTC(),
		},
		Status:     status,
		Components: []ghstatus.Component{component},
		Incidents: []ghstatus.Incident{
			incident2,
			incident1,
		},
	})
	clock.Advance(time.Minute)

	msg = waitForNotification(t, ch)

	require.Equal(t, notifier.Message{
		ChangedIncidents: []ghstatus.Incident{incident2},
	}, msg)

	// Let's add in a scheduled maintenance
	maintenance := ghstatus.ScheduledMaintenance{
		UpdatedAt: clock.Now().UTC(),
		Name:      "Maintenance",
	}
	server.SetSummary(t, ghstatus.SummaryResponse{
		Page: ghstatus.Page{
			UpdatedAt: clock.Now().UTC(),
		},
		Status:     status,
		Components: []ghstatus.Component{component},
		Incidents: []ghstatus.Incident{
			incident2,
			incident1,
		},
		ScheduledMaintenances: []ghstatus.ScheduledMaintenance{
			maintenance,
		},
	})
	clock.Advance(time.Minute)

	msg = waitForNotification(t, ch)

	require.Equal(t, notifier.Message{
		ChangedScheduledMaintenances: []ghstatus.ScheduledMaintenance{
			maintenance,
		},
	}, msg)
}

func waitForNotification(t *testing.T, ch chan notifier.Message) notifier.Message {
	select {
	case msg := <-ch:
		return msg
	case <-time.After(5 * time.Second):
		require.Fail(t, "timeout waiting for channel")
	}
	return notifier.Message{}
}
