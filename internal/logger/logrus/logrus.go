package logrus

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}
type Fields map[string]interface{}

func New(logLevel string, output ...io.Writer) (*Logger, error) {

	log := logrus.New()
	logrusLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	logrusOutput := io.MultiWriter(output...)
	log.SetOutput(logrusOutput)
	log.SetLevel(logrusLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	logger := &Logger{
		log: log,
	}
	return logger, nil
}
func (l *Logger) Tracef(message string, arg ...interface{}) {
	l.log.Tracef(message, arg...)
}
func (l *Logger) Debugf(message string, arg ...interface{}) {
	l.log.Debugf(message, arg...)
}
func (l *Logger) Infof(message string, arg ...interface{}) {
	l.log.Infof(message, arg...)
}
func (l *Logger) Warnf(message string, arg ...interface{}) {
	l.log.Warnf(message, arg...)
}
func (l *Logger) Errorf(message string, arg ...interface{}) {
	l.log.Errorf(message, arg...)
}
