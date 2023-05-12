package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghstatus",
	Short: "A tool for querying and monitoring the Github Status API",
	Long: "ghstatus provides utilities for manually querying and " +
		"monitoring Github's status using the Github Status API.",
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(componentsCmd)
	rootCmd.AddCommand(incidentsCmd)
	rootCmd.AddCommand(scheduledMaintenancesCmd)
	rootCmd.AddCommand(monitorCmd)
}

func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
