package render

import (
	"io"

	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"github.com/olekukonko/tablewriter"
)

// statusTable writes a rendered table representing the status to the writer.
func statusTable(w io.Writer, status ghstatus.Status) {
	table := newTable(w, "Indicator", "Description")
	table.Append([]string{string(status.Indicator), status.Description})
	table.Render()
}

// componentsTable writes a rendered table of components to the writer.
func componentsTable(w io.Writer, components []ghstatus.Component) {
	table := newTable(w, "Name", "Description", "Status", "Updated")
	for _, c := range components {
		table.Append([]string{string(c.Name), c.Description, string(c.Status), c.UpdatedAt.String()})
	}
	table.Render()
}

// incidentsTable writes a rendered table of incidents to the writer.
func incidentsTable(w io.Writer, incidents []ghstatus.Incident) {
	table := newTable(w, "Name", "Status", "Updated", "Latest Update")
	for _, i := range incidents {
		lastUpdate := ""
		if len(i.IncidentUpdates) > 0 {
			lastUpdate = i.IncidentUpdates[0].Body
		}
		table.Append([]string{string(i.Name), string(i.Status), i.UpdatedAt.String(), lastUpdate})
	}
	table.Render()
}

// scheduledMaintenancesTable writes a rendered table of scheduled maintenances to the writer.
func scheduledMaintenancesTable(w io.Writer, scheduledMaintenances []ghstatus.ScheduledMaintenance) {
	table := newTable(w, "Name", "Impact", "Status", "Scheduled For", "Scheduled Until")
	for _, s := range scheduledMaintenances {
		table.Append([]string{string(s.Name), string(s.Impact), string(s.Impact), s.ScheduledFor.String(), s.ScheduledUntil.String()})
	}
	table.Render()
}

// newTable returns a markdown compatible table generator.
func newTable(w io.Writer, headers ...string) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)
	return table
}
