package logrus

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log  *logrus.Logger
	file *os.File
}

func New(logLevel string, output ...io.Writer) (*Logger, error) {
	log := logrus.New()
	logrusLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	logrusOutput := io.MultiWriter(output...)
	log.SetOutput(logrusOutput)
	log.SetLevel(logrusLevel)
	log.SetFormatter(&logrus.TextFormatter{})

	logger := &Logger{
		log: log,
	}
	return logger, nil
}
func (l *Logger) Debugf(message string, arg ...interface{}) {
	l.log.Debugf(message, arg...)
}
func (l *Logger) Infof(message string, arg ...interface{}) {
	l.log.Infof(message, arg...)
}
func (l *Logger) Errorf(message string, arg ...interface{}) {
	l.log.Errorf(message, arg...)
}

func (l *Logger) Close() {
	l.file.Close()
}
