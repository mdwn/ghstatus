package ghstatus

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestServer is a test server with fixed responses that can be modified by
// callers.
type TestServer struct {
	Summary                       []byte
	Status                        []byte
	Components                    []byte
	UnresolvedIncidents           []byte
	AllIncidents                  []byte
	UpcomingScheduledMaintenances []byte
	ActiveScheduledMaintenances   []byte
	AllScheduledMaintenances      []byte
}

// SetSummary will set the summary response by encoding it into bytes.
func (ts *TestServer) SetSummary(t *testing.T, resp SummaryResponse) {
	ts.Summary = jsonEncode(t, resp)
}

// SetSummaryRaw will set the summary response to be the given raw bytes.
func (ts *TestServer) SetSummaryRaw(resp []byte) { ts.Summary = resp }

func (ts *TestServer) summary() []byte             { return ts.Summary }
func (ts *TestServer) status() []byte              { return ts.Status }
func (ts *TestServer) components() []byte          { return ts.Components }
func (ts *TestServer) unresolvedIncidents() []byte { return ts.UnresolvedIncidents }
func (ts *TestServer) allIncidents() []byte        { return ts.AllIncidents }
func (ts *TestServer) upcomingScheduledMaintenances() []byte {
	return ts.UpcomingScheduledMaintenances
}
func (ts *TestServer) activeScheduledMaintenances() []byte {
	return ts.ActiveScheduledMaintenances
}
func (ts *TestServer) allScheduledMaintenances() []byte {
	return ts.AllScheduledMaintenances
}

// NewTestServerAndClient creates a new test server with the client pointing to it.
func NewTestServerAndClient(t *testing.T) (*TestServer, Client) {
	testServer := &TestServer{}

	mux := http.NewServeMux()
	mux.HandleFunc(summaryEndpoint, response(t, testServer.summary))
	mux.HandleFunc(statusEndpoint, response(t, testServer.status))
	mux.HandleFunc(componentsEndpoint, response(t, testServer.components))
	mux.HandleFunc(unresolvedIncidentsEndpoint, response(t, testServer.unresolvedIncidents))
	mux.HandleFunc(allIncidentsEndpoint, response(t, testServer.allIncidents))
	mux.HandleFunc(upcomingScheduledMaintenancesEndpoint, response(t, testServer.upcomingScheduledMaintenances))
	mux.HandleFunc(activeScheduledMaintenancesEndpoint, response(t, testServer.activeScheduledMaintenances))
	mux.HandleFunc(allScheduledMaintenancesEndpoint, response(t, testServer.allScheduledMaintenances))

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return testServer, newClient(zap.NewNop(), server.URL)
}

func response(t *testing.T, response func() []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(response())
		require.NoError(t, err)
	}
}

func jsonEncode(t *testing.T, input any) []byte {
	data, err := json.Marshal(input)
	require.NoError(t, err)
	return data
}
