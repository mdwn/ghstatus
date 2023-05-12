package ghstatus

// SummaryResponse is a response from the summary endpoint.
type SummaryResponse struct {
	// Page is page metadata.
	Page Page `json:"page"`

	// Status is the current github status.
	Status Status `json:"status"`

	// Components are a list of components.
	Components []Component `json:"components"`

	// Incidents is a list of incidents.
	Incidents []Incident `json:"incidents"`

	// ScheduledMaintenances is a list of scheduled maintenances.
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
}

// StatusResponse is the response from the status endpoint.
type StatusResponse struct {
	// Page is page metadata.
	Page Page `json:"page"`

	// Status is the current github status.
	Status Status `json:"status"`
}

// ComponentsResponse is the response from the components endpoint.
type ComponentsResponse struct {
	// Page is page metadata.
	Page Page `json:"page"`

	// Components are a list of components.
	Components []Component `json:"components"`
}

// IncidentsResponse is the response from one of the incident endpoints.
type IncidentsResponse struct {
	// Page is page metadata.
	Page Page `json:"page"`

	// Incidents is a list of incidents.
	Incidents []Incident `json:"incidents"`
}

// ScheduledMaintenancesResponse is the response from one of the scheduled maintenance endpoints.
type ScheduledMaintenancesResponse struct {
	// Page is page metadata.
	Page Page `json:"page"`

	// ScheduledMaintenances is a list of scheduled maintenances.
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
}
