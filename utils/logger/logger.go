package logger

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func init() {
	logger.Level = logrus.InfoLevel
	logger.Formatter = &formatter{}

	logger.SetReportCaller(true)
}

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

type Fields logrus.Fields

func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		logger.Debug(fmt.Sprint(args...))
	}
}

func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		logger.Info(fmt.Sprint(args...))
	}
}

func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		logger.Warn(fmt.Sprint(args...))
	}
}

func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		logger.Error(fmt.Sprint(args...))
	}
}

func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		logger.Fatal(fmt.Sprint(args...))
	}
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

// Formatter implements logrus.Formatter interface.
type formatter struct{}

// Format building log message.
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var sb bytes.Buffer
	var color string
	switch entry.Level {
	case logrus.InfoLevel:
		color = colorGreen
	case logrus.WarnLevel:
		color = colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel:
		color = colorRed
	default:
		color = colorReset
	}

	timestamp := entry.Time.Format(time.RFC3339)
	level := strings.ToUpper(entry.Level.String())

	logLine := fmt.Sprintf("%s[%s] %s %s%s\n", color, level, timestamp, entry.Message, colorReset)
	sb.WriteString(logLine)

	return sb.Bytes(), nil
}
