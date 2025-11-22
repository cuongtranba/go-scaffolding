package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger
type Logger struct {
	*zerolog.Logger
}

// New creates a new logger with the specified level
func New(level string, writer io.Writer) *Logger {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	// Configure output
	var output io.Writer
	if writer != nil {
		output = writer
	} else {
		output = os.Stdout
		// Use console writer for development when writing to stdout
		if level == "debug" {
			output = zerolog.ConsoleWriter{
				Out:        output,
				TimeFormat: time.RFC3339,
			}
		}
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{Logger: &logger}
}

// With creates a child logger with additional context
func (l *Logger) With() zerolog.Context {
	return l.Logger.With()
}
