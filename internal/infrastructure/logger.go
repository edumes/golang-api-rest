package infrastructure

import (
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	Level      string
	Format     string
	OutputPath string
	Colors     bool
}

type ColoredFormatter struct {
	TimestampFormat string
}

func (f *ColoredFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(f.TimestampFormat)
	level := entry.Level.String()
	message := entry.Message

	var levelColor *color.Color
	switch entry.Level {
	case logrus.PanicLevel, logrus.FatalLevel:
		levelColor = color.New(color.FgRed, color.Bold)
	case logrus.ErrorLevel:
		levelColor = color.New(color.FgRed)
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow, color.Bold)
	case logrus.InfoLevel:
		levelColor = color.New(color.FgGreen)
	case logrus.DebugLevel:
		levelColor = color.New(color.FgCyan)
	case logrus.TraceLevel:
		levelColor = color.New(color.FgMagenta)
	default:
		levelColor = color.New(color.FgWhite)
	}

	timestampColor := color.New(color.FgBlue)
	coloredTimestamp := timestampColor.Sprint(timestamp)

	coloredLevel := levelColor.Sprint(level)

	var messageColor *color.Color
	switch {
	case strings.Contains(strings.ToLower(message), "error") || strings.Contains(strings.ToLower(message), "failed"):
		messageColor = color.New(color.FgRed)
	case strings.Contains(strings.ToLower(message), "success"):
		messageColor = color.New(color.FgGreen)
	case strings.Contains(strings.ToLower(message), "warning") || strings.Contains(strings.ToLower(message), "warn"):
		messageColor = color.New(color.FgYellow)
	case strings.Contains(strings.ToLower(message), "debug"):
		messageColor = color.New(color.FgCyan)
	case strings.Contains(strings.ToLower(message), "info"):
		messageColor = color.New(color.FgBlue)
	default:
		messageColor = color.New(color.FgWhite)
	}

	coloredMessage := messageColor.Sprint(message)

	fieldsStr := ""
	if len(entry.Data) > 0 {
		fieldColor := color.New(color.FgMagenta)
		for key, value := range entry.Data {
			if fieldsStr != "" {
				fieldsStr += " "
			}
			fieldsStr += fieldColor.Sprintf("%s=%v", key, value)
		}
	}

	var result string
	if fieldsStr != "" {
		result = coloredTimestamp + " " + coloredLevel + " " + coloredMessage + " " + fieldsStr + "\n"
	} else {
		result = coloredTimestamp + " " + coloredLevel + " " + coloredMessage + "\n"
	}

	return []byte(result), nil
}

func NewLogger(config LoggerConfig) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "colored":
		logger.SetFormatter(&ColoredFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "text":
		fallthrough
	default:
		if config.Colors {
			logger.SetFormatter(&ColoredFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			})
		} else {
			logger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			})
		}
	}

	if config.OutputPath != "" {
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		}
	} else {
		logger.SetOutput(os.Stdout)
	}

	return logger
}

func GetDefaultLogger() *logrus.Logger {
	config := LoggerConfig{
		Level:  "info",
		Format: "text",
		Colors: true,
	}
	return NewLogger(config)
}

func GetDebugLogger() *logrus.Logger {
	config := LoggerConfig{
		Level:  "debug",
		Format: "colored",
		Colors: true,
	}
	return NewLogger(config)
}

func GetProductionLogger() *logrus.Logger {
	config := LoggerConfig{
		Level:  "info",
		Format: "json",
		Colors: false,
	}
	return NewLogger(config)
}

func GetColoredLogger() *logrus.Logger {
	config := LoggerConfig{
		Level:  "debug",
		Format: "colored",
		Colors: true,
	}
	return NewLogger(config)
}
