package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/mdwn/ghstatus/pkg/ghstatus"
	"github.com/mdwn/ghstatus/pkg/logging"
	"github.com/mdwn/ghstatus/pkg/monitor"
	"github.com/mdwn/ghstatus/pkg/notifiers"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	monitorNotifiers        []string
	monitorNotifyOnFirstRun bool

	monitorCmd = &cobra.Command{
		Use:   "monitor",
		Short: "Monitors the Github status",
		Long: func() string {
			builder := strings.Builder{}

			builder.WriteString("Monitor will monitor the Github Status and report changes to the configued notifiers.\n\n")
			builder.WriteString("Available notifiers:\n")
			for _, name := range notifiers.ListNotifiers() {
				builder.WriteString(fmt.Sprintf(" - %s\n", name))
			}

			return builder.String()
		}(),

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			log, err := logging.NewLogger()
			if err != nil {
				return fmt.Errorf("error creating logger: %w", err)
			}

			client := ghstatus.NewClient(log)
			clock := clockwork.NewRealClock()

			if len(monitorNotifiers) == 0 {
				return errors.New("no notifiers configured")
			}

			monitor := monitor.New(log, clock, client, monitorNotifyOnFirstRun)

			for _, name := range monitorNotifiers {
				notifier, err := notifiers.GetNotifier(log, name)
				if err != nil {
					return fmt.Errorf("error creating notifier %s: %w", name, err)
				}
				defer func() {
					if err := notifier.Cleanup(); err != nil {
						log.With(zap.Error(err), zap.String("notifier", notifier.Name())).Error("error cleaning up")
					}
				}()
				if err := monitor.RegisterNotifier(notifier); err != nil {
					return fmt.Errorf("error registering notifier %s: %w", name, err)
				}
			}

			monitor.MonitorAndNotify(ctx, time.Minute)

			return nil
		},
	}
)

func init() {
	monitorCmd.Flags().StringSliceVarP(&monitorNotifiers, "notifiers", "n", []string{notifiers.Stdout}, "The notifiers to use for the monitor.")
	monitorCmd.Flags().BoolVarP(&monitorNotifyOnFirstRun, "notify-on-first-run", "f", false, "Whether the monitor should send notifications on the first run.")
	notifiers.RegisterCommandFlags(monitorCmd)
}
