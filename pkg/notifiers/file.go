package notifiers

import (
	"errors"
	"fmt"
	"os"

	"github.com/mdwn/ghstatus/pkg/notifier"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	File = "file"
)

var (
	fileNotifierFilepath string
)

func init() {
	if err := RegisterNotifier(File, NewFileNotifier); err != nil {
		panic(err.Error())
	}

	flags := pflag.NewFlagSet("file-notifier", pflag.ContinueOnError)
	flags.StringVar(&fileNotifierFilepath, "fn-file-path", "", "The file to use for the file notifier.")
	flags.FlagUsages()
	notifierFlags.AddFlagSet(flags)
}

// FileNotifier writes the output to the given file.
type FileNotifier struct {
	*WriterNotifier

	file os.File
}

// NewFileNotifier will return a file notifier.
func NewFileNotifier(_ *zap.Logger) (notifier.Notifier, error) {
	if fileNotifierFilepath == "" {
		return nil, errors.New("file notifier needs the file path to be set")
	}

	file, err := os.OpenFile(fileNotifierFilepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("error opening file for writing: %w", err)
	}

	return &FileNotifier{
		WriterNotifier: NewWriterNotifier(file),
	}, nil
}

// Name is the name of the notifier.
func (*FileNotifier) Name() string {
	return File
}

// Cleanup performs any cleanup steps.
func (f *FileNotifier) Cleanup() error {
	return f.file.Close()
}
