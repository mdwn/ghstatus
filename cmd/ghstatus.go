package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"github.com/mdwn/ghstatus/pkg/ghstatus/render"
	"github.com/mdwn/ghstatus/pkg/logging"
	"github.com/spf13/cobra"
)

// Below are commands that wrap the various Github client methods
// and render them in a human readable manner.

var (
	format string
)

var (
	summaryCmd = &cobra.Command{
		Use:   "summary",
		Short: "Provides the summary",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.Summary); err != nil {
				return fmt.Errorf("error getting summary: %w", err)
			}

			return nil
		},
	}

	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Provides the status",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.Status); err != nil {
				return fmt.Errorf("error getting status: %w", err)
			}

			return nil
		},
	}

	componentsCmd = &cobra.Command{
		Use:   "components",
		Short: "Provides the list of components",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.Components); err != nil {
				return fmt.Errorf("error getting components: %w", err)
			}

			return nil
		},
	}

	incidentsCmd = &cobra.Command{
		Use:   "incidents",
		Short: "Retrieve different types of incidents",
	}

	unresolvedIncidentsCmd = &cobra.Command{
		Use:   "unresolved",
		Short: "Provides the list of unresolved incidents",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.UnresolvedIncidents); err != nil {
				return fmt.Errorf("error getting unresolved: %w", err)
			}

			return nil
		},
	}

	allIncidentsCmd = &cobra.Command{
		Use:   "all",
		Short: "Provides the list of all Github Status incidents",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.AllIncidents); err != nil {
				return fmt.Errorf("error getting all incidents: %w", err)
			}

			return nil
		},
	}

	scheduledMaintenancesCmd = &cobra.Command{
		Use:   "scheduled-maintenances",
		Short: "Retrieve different types of scheduled maintenances",
	}

	upcomingScheduledMaintenancesCmd = &cobra.Command{
		Use:   "upcoming",
		Short: "Provides the list of upcoming scheduled maintenances",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.UpcomingScheduledMaintenances); err != nil {
				return fmt.Errorf("error getting upcoming scheduled maintenances: %w", err)
			}

			return nil
		},
	}

	activeScheduledMaintenancesCmd = &cobra.Command{
		Use:   "active",
		Short: "Provides the list of active scheduled maintenances",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.UpcomingScheduledMaintenances); err != nil {
				return fmt.Errorf("error getting active scheduled maintenances: %w", err)
			}

			return nil
		},
	}

	allScheduledMaintenancesCmd = &cobra.Command{
		Use:   "all",
		Short: "Provides the list of all scheduled maintenances",

		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}
			client := ghstatus.NewClient(log)

			if err := printResponse(cmd.Context(), client.AllScheduledMaintenances); err != nil {
				return fmt.Errorf("error getting all scheduled maintenances: %w", err)
			}

			return nil
		},
	}
)

func init() {
	incidentsCmd.AddCommand(unresolvedIncidentsCmd)
	incidentsCmd.AddCommand(allIncidentsCmd)

	scheduledMaintenancesCmd.AddCommand(upcomingScheduledMaintenancesCmd)
	scheduledMaintenancesCmd.AddCommand(activeScheduledMaintenancesCmd)
	scheduledMaintenancesCmd.AddCommand(allScheduledMaintenancesCmd)
}

func init() {
	addOutputFlag(summaryCmd)
	addOutputFlag(statusCmd)
	addOutputFlag(componentsCmd)
	addOutputFlag(unresolvedIncidentsCmd)
	addOutputFlag(allIncidentsCmd)
	addOutputFlag(upcomingScheduledMaintenancesCmd)
	addOutputFlag(activeScheduledMaintenancesCmd)
	addOutputFlag(allScheduledMaintenancesCmd)
}

// addOutputFlag will add the output flag to the given command.
func addOutputFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&format, "format", "f", "table", "Output format (valid values are [yaml, json, table])")
}

// printResponse will get an object and then print the response for the object.
func printResponse[T any](ctx context.Context, getFn func(context.Context) (T, error)) error {
	format, err := render.FormatFromString(format)
	if err != nil {
		return fmt.Errorf("error getting format from string: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := getFn(ctx)
	if err != nil {
		return fmt.Errorf("error getting response: %w", err)
	}

	out, err := render.Render(resp, format)
	if err != nil {
		return fmt.Errorf("error rendering response: %w", err)
	}

	fmt.Print(out)

	return nil
}
