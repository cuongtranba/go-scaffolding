package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	log := New("debug", &buf)

	log.Info().Msg("test message")

	var logEntry map[string]interface{}
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "test message", logEntry["message"])
	assert.Equal(t, "info", logEntry["level"])
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		logFunc  string
		expected bool
	}{
		{"debug enabled", "debug", "debug", true},
		{"info enabled", "info", "info", true},
		{"debug filtered", "info", "debug", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log := New(tt.level, &buf)

			switch tt.logFunc {
			case "debug":
				log.Debug().Msg("test")
			case "info":
				log.Info().Msg("test")
			}

			if tt.expected {
				assert.Greater(t, buf.Len(), 0)
			} else {
				assert.Equal(t, 0, buf.Len())
			}
		})
	}
}
