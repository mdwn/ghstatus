package render

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"gopkg.in/yaml.v3"
)

// Format is the format to render.
type Format int

const (
	YAML Format = iota
	JSON
	Table
)

// FormatFromString returns a Format from a string descriptor.
func FormatFromString(format string) (Format, error) {
	switch format {
	case "yaml":
		return YAML, nil
	case "json":
		return JSON, nil
	case "table":
		return Table, nil
	}

	return 0, fmt.Errorf("unrecognized format string %T", format)
}

// Render will render the target with the given output type and return it as a string.
func Render(target any, format Format) (string, error) {
	switch format {
	case YAML:
		buf := bytes.NewBuffer(nil)
		encoder := yaml.NewEncoder(buf)
		if err := encoder.Encode(target); err != nil {
			return "", fmt.Errorf("error encoding as YAML: %w", err)
		}
		return buf.String(), nil
	case JSON:
		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		if err := encoder.Encode(target); err != nil {
			return "", fmt.Errorf("error encoding as JSON: %w", err)
		}
		return buf.String(), nil
	case Table:
		switch t := target.(type) {
		case ghstatus.SummaryResponse:
			return summaryResponseWithTables(t), nil
		case ghstatus.StatusResponse:
			return statusResponseWithTables(t), nil
		case ghstatus.ComponentsResponse:
			return componentsResponseWithTables(t), nil
		case ghstatus.IncidentsResponse:
			return incidentsResponseWithTables(t), nil
		case ghstatus.ScheduledMaintenancesResponse:
			return scheduledMaintenancesResponseWithTables(t), nil
		default:
			return "", fmt.Errorf("type %T does not support table rendering", target)
		}
	default:
		return "", fmt.Errorf("unrecognized format: %d", format)
	}
}
