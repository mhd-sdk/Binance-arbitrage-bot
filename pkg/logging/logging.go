package logging

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type Logger struct {
	Slog *slog.Logger
}

func NewLogger() (*Logger, error) {

	// if /logs folder does not exist, create it
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
	}
	logger := slog.New(tint.NewHandler(os.Stdout, nil))

	return &Logger{Slog: logger}, nil
}

func (l Logger) FileLog(filePath string, content string) {
	// if file does not exist, create it
	// else, append to it
	file := &os.File{}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// create the file inside the logs folder
		file, err = os.Create(filePath)
		if err != nil {
			l.Slog.Error(err.Error())
			return
		}
	} else {
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			l.Slog.Error(err.Error())
			return
		}
		defer file.Close()
	}

	// write to the file
	_, err := file.WriteString(content)
	if err != nil {
		l.Slog.Error(err.Error())
		return
	}
}
