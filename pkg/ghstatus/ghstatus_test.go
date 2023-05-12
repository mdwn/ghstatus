package ghstatus

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/summary.json
	summaryResponse []byte
)

func TestUnmarshalActualResponse(t *testing.T) {
	server, client := NewTestServerAndClient(t)
	server.SetSummaryRaw(summaryResponse)

	ctx := context.Background()
	summary, err := client.Summary(ctx)
	require.NoError(t, err)

	require.Len(t, summary.Components, 11)
	require.Len(t, summary.Incidents, 1)
	require.Len(t, summary.Incidents[0].IncidentUpdates, 15)
	require.Empty(t, summary.ScheduledMaintenances)
}
