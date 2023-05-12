package render

import (
	"bytes"

	"github.com/mdwn/ghstatus/pkg/ghstatus"
)

// summaryResponseWithTables will render the summary response with tables.
func summaryResponseWithTables(s ghstatus.SummaryResponse) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("# Status\n\n")
	statusTable(buf, s.Status)

	if len(s.Components) != 0 {
		buf.WriteString("\n\n")
		buf.WriteString("# Components\n\n")
		componentsTable(buf, s.Components)
	}

	if len(s.Incidents) != 0 {
		buf.WriteString("\n\n")
		buf.WriteString("# Incidents\n\n")
		incidentsTable(buf, s.Incidents)
	}

	if len(s.ScheduledMaintenances) != 0 {
		buf.WriteString("\n\n")
		buf.WriteString("# Scheduled Maintenances\n\n")
		scheduledMaintenancesTable(buf, s.ScheduledMaintenances)
	}

	return buf.String()
}

// statusResponseWithTables will render the status response with tables.
func statusResponseWithTables(s ghstatus.StatusResponse) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("# Status\n\n")
	statusTable(buf, s.Status)

	return buf.String()
}

// componentsResponseWithTables will render the components response with tables.
func componentsResponseWithTables(c ghstatus.ComponentsResponse) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("# Components\n\n")
	componentsTable(buf, c.Components)

	return buf.String()
}

// incidentsResponseWithTables will render an incidents response with tables.
func incidentsResponseWithTables(i ghstatus.IncidentsResponse) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("# Incidents\n\n")
	incidentsTable(buf, i.Incidents)

	return buf.String()
}

// scheduledMaintenancesResponseWithTables will render a scheduled maintenances response with tables.
func scheduledMaintenancesResponseWithTables(s ghstatus.ScheduledMaintenancesResponse) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("# Scheduled Maintenances\n\n")
	scheduledMaintenancesTable(buf, s.ScheduledMaintenances)

	return buf.String()
}
