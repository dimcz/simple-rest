package util

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

var e *logrus.Entry

func GetLogger() *Logger {
	return &Logger{e}
}

func init() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	if err := os.MkdirAll("logs", os.FileMode(0o700)); err != nil {
		log.Fatalf("logging.Setup(Mkdir) error: %v", err)
	}

	logFile, err := os.OpenFile("logs/all.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0o600))
	if err != nil {
		log.Fatalf("logging.Setup(OpenFile) error: %v", err)
	}

	l.SetOutput(logFile)

	e = logrus.NewEntry(l)
}
