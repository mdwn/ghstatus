package notifiers

import (
	"os"

	"github.com/mdwn/ghstatus/pkg/notifier"
)

const (
	Stdout = "stdout"
)

func init() {
	registeredNotifiers[Stdout] = NewStdoutNotifier
}

// StdoutNotifier writes the output to stdout.
type StdoutNotifier struct {
	*WriterNotifier
}

// NewStdoutNotifier will return an stdout notifier.
func NewStdoutNotifier() (notifier.Notifier, error) {
	return &StdoutNotifier{
		WriterNotifier: NewWriterNotifier(os.Stdout),
	}, nil
}

// Name is the name of the notifier.
func (*StdoutNotifier) Name() string {
	return Stdout
}
