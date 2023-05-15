package notifiers

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mdwn/ghstatus/pkg/notifier"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var (
	registeredNotifiersMu sync.RWMutex
	registeredNotifiers   = map[string]NotifierCreator{}
	notifierFlags         = pflag.NewFlagSet("notifiers", pflag.ContinueOnError)
)

func init() {
}

// RegisterNotifier will register the given notifier and the given notifier creation function.
func RegisterNotifier(name string, creator NotifierCreator) error {
	registeredNotifiersMu.Lock()
	defer registeredNotifiersMu.Unlock()

	_, ok := registeredNotifiers[name]
	if ok {
		return fmt.Errorf("duplicate notifier %s", name)
	}

	registeredNotifiers[name] = creator
	return nil
}

// RegisterCommandFlags will register CLI flags for all commands with the given command.
func RegisterCommandFlags(cmd *cobra.Command) {
	cmd.Flags().AddFlagSet(notifierFlags)
}

// NotifierCreator creates a notifier.
type NotifierCreator func(log *zap.Logger) (notifier.Notifier, error)

// ListNotifiers will list all notifier names.
func ListNotifiers() []string {
	registeredNotifiersMu.RLock()
	notifiers := make([]string, 0, len(registeredNotifiers))
	for name := range registeredNotifiers {
		notifiers = append(notifiers, name)
	}
	registeredNotifiersMu.RUnlock()

	sort.Strings(notifiers)
	return notifiers
}

// GetNotifier will retrieve the notifier with the given name.
func GetNotifier(log *zap.Logger, name string) (notifier.Notifier, error) {
	registeredNotifiersMu.RLock()
	defer registeredNotifiersMu.RUnlock()

	notifierCreator, ok := registeredNotifiers[name]
	if !ok {
		return nil, fmt.Errorf("no notifier named %s", name)
	}

	notifier, err := notifierCreator(log)
	if err != nil {
		return nil, fmt.Errorf("error creating notifier %s: %w", name, err)
	}

	return notifier, nil
}
