package notifiers

import (
	"os"

	"github.com/mdwn/ghstatus/pkg/notifier"
	"go.uber.org/zap"
)

const (
	Stdout = "stdout"
)

func init() {
	if err := RegisterNotifier(Stdout, NewStdoutNotifier); err != nil {
		panic(err.Error())
	}
}

// StdoutNotifier writes the output to stdout.
type StdoutNotifier struct {
	*WriterNotifier
}

// NewStdoutNotifier will return an stdout notifier.
func NewStdoutNotifier(_ *zap.Logger) (notifier.Notifier, error) {
	return &StdoutNotifier{
		WriterNotifier: NewWriterNotifier(os.Stdout),
	}, nil
}

// Name is the name of the notifier.
func (*StdoutNotifier) Name() string {
	return Stdout
}
