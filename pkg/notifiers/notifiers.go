package notifiers

import (
	"fmt"
	"sort"

	"github.com/mdwn/ghstatus/pkg/notifier"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	registeredNotifiers = map[string]NotifierCreator{}
	notifierFlags       = pflag.NewFlagSet("notifiers", pflag.ContinueOnError)
)

func init() {
}

// RegisterCommandFlags will register CLI flags for all commands with the given command.
func RegisterCommandFlags(cmd *cobra.Command) {
	cmd.Flags().AddFlagSet(notifierFlags)
}

// NotifierCreator creates a notifier.
type NotifierCreator func() (notifier.Notifier, error)

// ListNotifiers will list all notifier names.
func ListNotifiers() []string {
	notifiers := make([]string, 0, len(registeredNotifiers))
	for name := range registeredNotifiers {
		notifiers = append(notifiers, name)
	}
	sort.Strings(notifiers)
	return notifiers
}

// GetNotifier will retrieve the notifier with the given name.
func GetNotifier(name string) (notifier.Notifier, error) {
	notifierCreator, ok := registeredNotifiers[name]
	if !ok {
		return nil, fmt.Errorf("no notifier named %s", name)
	}

	notifier, err := notifierCreator()
	if err != nil {
		return nil, fmt.Errorf("error creating notifier %s: %w", name, err)
	}

	return notifier, nil
}
