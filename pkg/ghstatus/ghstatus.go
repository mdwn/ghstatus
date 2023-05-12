package ghstatus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

const (
	githubStatusURL = "https://www.githubstatus.com"

	// The number of times to retry fetching from the Github Status API on failure.
	maxRetries = 5

	// The minimum and maximum amount of time to wait when retrying.
	retryMin = 1 * time.Second
	retryMax = 5 * time.Second

	summaryEndpoint    = "/api/v2/summary.json"
	statusEndpoint     = "/api/v2/status.json"
	componentsEndpoint = "/api/v2/components.json"

	// Incidents endpoints
	unresolvedIncidentsEndpoint = "/api/v2/incidents/unresolved.json"
	allIncidentsEndpoint        = "/api/v2/incidents.json"

	// Scheduled maintenance endpoints
	upcomingScheduledMaintenancesEndpoint = "/api/v2/scheduled-maintenances/upcoming.json"
	activeScheduledMaintenancesEndpoint   = "/api/v2/scheduled-maintenances/active.json"
	allScheduledMaintenancesEndpoint      = "/api/v2/scheduled-maintenances.json"
)

// Client is the github status client.
type Client interface {
	// Summary returns the summary.
	Summary(ctx context.Context) (SummaryResponse, error)

	// Status returns the status.
	Status(ctx context.Context) (StatusResponse, error)

	// Components returns the components.
	Components(ctx context.Context) (ComponentsResponse, error)

	// UnresolvedIncidents returns the unresolved incidents.
	UnresolvedIncidents(ctx context.Context) (IncidentsResponse, error)

	// AllIncidents returns all incidents.
	AllIncidents(ctx context.Context) (IncidentsResponse, error)

	// UpcomingScheduledMaintenances returns all upcoming scheduled maintenances.
	UpcomingScheduledMaintenances(ctx context.Context) (ScheduledMaintenancesResponse, error)

	// ActiveScheduledMaintenances returns all active scheduled maintenances.
	ActiveScheduledMaintenances(ctx context.Context) (ScheduledMaintenancesResponse, error)

	// AllScheduledMaintenances returns all scheduled maintenances.
	AllScheduledMaintenances(ctx context.Context) (ScheduledMaintenancesResponse, error)
}

type client struct {
	endpoint   string
	httpClient *retryablehttp.Client
}

// NewClient creates a new github status client.
func NewClient(log *zap.Logger) Client {
	return newClient(log, githubStatusURL)
}

func newClient(log *zap.Logger, endpointURL string) *client {
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = maxRetries
	httpClient.RetryWaitMin = retryMin
	httpClient.RetryWaitMax = retryMax
	stdLog, err := zap.NewStdLogAt(log, zap.DebugLevel)
	if err != nil {
		panic(fmt.Sprintf("panic creating standard log: %v", err))
	}
	httpClient.Logger = stdLog

	return &client{
		endpoint:   endpointURL,
		httpClient: httpClient,
	}
}

// Summary returns the summary.
func (c *client) Summary(ctx context.Context) (SummaryResponse, error) {
	var resp SummaryResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, summaryEndpoint, &resp); err != nil {
		return SummaryResponse{}, err
	}
	return resp, nil
}

// Status returns the status.
func (c *client) Status(ctx context.Context) (StatusResponse, error) {
	var resp StatusResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, statusEndpoint, &resp); err != nil {
		return StatusResponse{}, err
	}
	return resp, nil
}

// Components returns the components.
func (c *client) Components(ctx context.Context) (ComponentsResponse, error) {
	var resp ComponentsResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, componentsEndpoint, &resp); err != nil {
		return ComponentsResponse{}, err
	}
	return resp, nil
}

// UnresolvedIncidents returns the unresolved incidents.
func (c *client) UnresolvedIncidents(ctx context.Context) (IncidentsResponse, error) {
	var resp IncidentsResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, unresolvedIncidentsEndpoint, &resp); err != nil {
		return IncidentsResponse{}, err
	}
	return resp, nil
}

// AllIncidents returns all incidents.
func (c *client) AllIncidents(ctx context.Context) (IncidentsResponse, error) {
	var resp IncidentsResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, allIncidentsEndpoint, &resp); err != nil {
		return IncidentsResponse{}, err
	}
	return resp, nil
}

// UpcomingScheduledMaintenances returns all upcoming scheduled maintenances.
func (c *client) UpcomingScheduledMaintenances(ctx context.Context) (ScheduledMaintenancesResponse, error) {
	var resp ScheduledMaintenancesResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, upcomingScheduledMaintenancesEndpoint, &resp); err != nil {
		return ScheduledMaintenancesResponse{}, err
	}
	return resp, nil
}

// ActiveScheduledMaintenances returns all active scheduled maintenances.
func (c *client) ActiveScheduledMaintenances(ctx context.Context) (ScheduledMaintenancesResponse, error) {
	var resp ScheduledMaintenancesResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, activeScheduledMaintenancesEndpoint, &resp); err != nil {
		return ScheduledMaintenancesResponse{}, err
	}
	return resp, nil
}

// AllScheduledMaintenances returns all scheduled maintenances.
func (c *client) AllScheduledMaintenances(ctx context.Context) (ScheduledMaintenancesResponse, error) {
	var resp ScheduledMaintenancesResponse
	if err := getAndUnmarshal(ctx, c.httpClient, c.endpoint, allScheduledMaintenancesEndpoint, &resp); err != nil {
		return ScheduledMaintenancesResponse{}, err
	}
	return resp, nil
}

func getAndUnmarshal[T any](ctx context.Context, client *retryablehttp.Client, prefix, suffix string, target T) error {
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", prefix, suffix), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
}
