package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Debug(component string, format string, a ...any)
	Info(component string, format string, a ...any)
	Warn(component string, format string, a ...any)
	Error(component string, format string, a ...any)
	Fatal(component string, format string, a ...any)
}

type Log struct {
}

func NewLog() Logger {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log.Logger = log.Output(output)
	return &Log{}
}

func (l *Log) Debug(component string, format string, a ...any) {
	log.Debug().Msgf(fmt.Sprintf("%-6s | %s", component, format), a...)
}

func (l *Log) Info(component string, format string, a ...any) {
	log.Info().Msgf(fmt.Sprintf("%-6s | %s", component, format), a...)
}

func (l *Log) Warn(component string, format string, a ...any) {
	log.Warn().Msgf(fmt.Sprintf("%-6s | %s", component, format), a...)
}

func (l *Log) Error(component string, format string, a ...any) {
	log.Error().Msgf(fmt.Sprintf("%-6s | %s", component, format), a...)
}

func (l *Log) Fatal(component string, format string, a ...any) {
	log.Fatal().Msgf(fmt.Sprintf("| %-6s |%s", component, format), a...)
}
