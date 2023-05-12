package ghstatus

import (
	"time"
)

// Page describes metadata for the status page.
type Page struct {
	// ID is the ID of the page.
	ID string `json:"id" yaml:"id"`

	// Name is the name of the page.
	Name string `json:"name" yaml:"name"`

	// URL is the URL of the page retrieved.
	URL string `json:"url" yaml:"url"`

	// UpdatedAt is when the page was last updated.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// Indicator is an indicator of the status of an incident or the overall status.
type Indicator string

const (
	None     Indicator = "none"
	Minor    Indicator = "minor"
	Major    Indicator = "major"
	Critical Indicator = "critical"
)

// Status is an overall description of the current github status.
type Status struct {
	// Description is the description of the overall status.
	Description string `json:"description" yaml:"description"`

	// Indicator is the indicator of the overall status.
	Indicator Indicator `json:"indicator" yaml:"indicator"`
}

// ComponentStatus is the current status of a github component.
type ComponentStatus string

const (
	Operational         ComponentStatus = "operational"
	DegradedPerformance ComponentStatus = "degraded_performance"
	PartialOutage       ComponentStatus = "partial_outage"
	MajorOutage         ComponentStatus = "major_outage"
)

// Component is a github component along with its current status.
type Component struct {
	// CreatedAt is when the component was created.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// Description is the description of the component.
	Description string `json:"description" yaml:"description"`

	// Group is is whether this component is in a group.
	Group bool `json:"group" yaml:"group"`

	// GroupID is the ID of the group this component belongs to.
	GroupID string `json:"group_id" yaml:"group_id"`

	// Name is the name of the component.
	Name string `json:"name" yaml:"name"`

	// OnlyShowIfDegraded designates whether this should only be shown if the component is degraded.
	OnlyShowIfDegraded bool `json:"only_show_if_degraded" yaml:"only_show_if_degraded"`

	// PageID is the page ID this component belongs to.
	PageID string `json:"page_id" yaml:"page_id"`

	// Position is the position of the component on the page.
	Position int `json:"position" yaml:"position"`

	// Showcase is whether the component should be showcased.
	Showcase bool `json:"showcase" yaml:"showcase"`

	// StartDate is the start date for the component.
	StartDate string `json:"start_date" yaml:"start_date"`

	// Status is the status of the component.
	Status ComponentStatus `json:"status" yaml:"status"`

	// UpdatedAt is when the component was last updated.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// IncidentStatus is the status of an incident.
type IncidentStatus string

const (
	Investigating IncidentStatus = "investigating"
	Identified    IncidentStatus = "identified"
	Monitoring    IncidentStatus = "monitoring"
	Resolved      IncidentStatus = "resolved"
	Postmorten    IncidentStatus = "postmortem"
)

// Incident is an ongoing incident.
type Incident struct {
	// CreatedAt is when the incident was created.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// ID is the identifier of the incident.
	ID string `json:"id" yaml:"id"`

	// Impact is the impact of the incident.
	Impact Indicator `json:"impact" yaml:"impact"`

	// IncidentUpdates are the updates for an incident.
	IncidentUpdates []IncidentUpdate `json:"incident_updates" yaml:"incident_updates"`

	// MonitoringAt is the time when the incident is being monitored.
	MonitoringAt time.Time `json:"monitoring_at" yaml:"monitoring_at"`

	// Name is the name of the incident.
	Name string `json:"name" yaml:"name"`

	// PageID is the ID of this page.
	PageID string `json:"page_id" yaml:"page_id"`

	// ResolvedAt is when the incident was resolved.
	ResolvedAt time.Time `json:"resolved_at" yaml:"resolved_at"`

	// ShortLink is the shortlink to this incident.
	Shortlink string `json:"shortlink" yaml:"shortlink"`

	// Status is the status of the incident.
	Status IncidentStatus `json:"status" yaml:"status"`

	// UpdatedAt is when the incident was last updated.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// IncidentUpdate is an update to an incident.
type IncidentUpdate struct {
	// Body is the plaintext description of the update.
	Body string `json:"body" yaml:"body"`

	// CreatedAt is when the update was created.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// DisplayAt is when the update should be displayed.
	DisplayAt time.Time `json:"display_at" yaml:"display_at"`

	// ID is the identifier of the update.
	ID string `json:"id" yaml:"id"`

	// IncidentID is the identifier of the associated incident.
	IncidentID string `json:"incident_id" yaml:"incident_id"`

	// Status is the status of the incident associated with this update.
	Status IncidentStatus `json:"status" yaml:"status"`

	// UpdatedAt is when this update entry was last updated.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// ScheduledMaintenanceStatus is the status of a scheduled maintenance.
type ScheduledMaintenanceStatus string

const (
	Scheduled  ScheduledMaintenanceStatus = "scheduled"
	InProgress ScheduledMaintenanceStatus = "in_progress"
	Verifying  ScheduledMaintenanceStatus = "verifying"
	Completed  ScheduledMaintenanceStatus = "completed"
)

// ScheduledMaintenance is a scheduled maintenance.
type ScheduledMaintenance struct {
	// CreatedAt is when the scheduled maintenance was created.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// ID is the identifier of the scheduled maintenance.
	ID string `json:"id" yaml:"id"`

	// Impact is the impact of the scheduled maintenance.
	Impact Indicator `json:"impact" yaml:"impact"`

	// IncidentUpdates are updates to the scheduled maintenance.
	IncidentUpdates []IncidentUpdate `json:"incident_updates" yaml:"incident_updates"`

	// MonitoringAt is the time when the scheduled maintenance is being monitored.
	MonitoringAt time.Time `json:"monitoring_at" yaml:"monitoring_at"`

	// Name is the name of the scheduled maintenance.
	Name string `json:"name" yaml:"name"`

	// PageID is the ID of this page.
	PageID string `json:"page_id" yaml:"page_id"`

	// ResolvedAt is the time when the scheduled maintenance is resolved.
	ResolvedAt time.Time `json:"resolved_at" yaml:"resolved_at"`

	// ScheduledFor is when the scheduled maintenance is scheduled.
	ScheduledFor time.Time `json:"scheduled_for" yaml:"scheduled_for"`

	// ScheduledUntil is when the scheduled maintenance is supposed to end.
	ScheduledUntil time.Time `json:"scheduled_until" yaml:"scheduled_until"`

	// Status is the status of the scheduled maintenance.
	Status ScheduledMaintenanceStatus `json:"status" yaml:"status"`

	// UpdatedAt is the time when the scheduled maintenance is updated.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}
