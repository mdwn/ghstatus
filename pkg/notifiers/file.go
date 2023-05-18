package notifiers

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/mdwn/ghstatus/pkg/notifier"
	"github.com/ory/viper"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	File = "file"

	fileNotifierFilepathCfg  = "file.filepath"
	fileNotifierFilepathFlag = "fn-filepath"
	fileNotifierFilepathEnv  = "FN_FILEPATH"
)

func init() {
	if err := RegisterNotifier(File, NewFileNotifier); err != nil {
		panic(err.Error())
	}

	flags := pflag.NewFlagSet("file-notifier", pflag.ContinueOnError)
	flags.String(fileNotifierFilepathFlag, "", "The file to use for the file notifier.")

	notifierFlags.AddFlagSet(flags)

	err := multierror.Append(nil,
		viper.BindPFlag(fileNotifierFilepathCfg, flags.Lookup(fileNotifierFilepathFlag)),
		viper.BindEnv(fileNotifierFilepathCfg, fileNotifierFilepathEnv),
	)

	if err.ErrorOrNil() != nil {
		panic(fmt.Sprintf("error binding file notifier configs: %v", err))
	}
}

// FileNotifier writes the output to the given file.
type FileNotifier struct {
	*WriterNotifier

	file os.File
}

// NewFileNotifier will return a file notifier.
func NewFileNotifier(_ *zap.Logger) (notifier.Notifier, error) {
	fileNotifierFilepath := viper.GetString(fileNotifierFilepathCfg)

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
