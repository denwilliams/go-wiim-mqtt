package logging_test

import (
	"bytes"
	"testing"

	"github.com/denwilliams/go-wiim-mqtt/internal/logging"
)

func TestLogging(t *testing.T) {
	var buf bytes.Buffer

	logging.Init(&buf, 0)

	logging.Debug("debug message")
	if buf.String() != "" {
		t.Errorf("Debug message should not be logged")
	}

	logging.Info("info message")
	if got, want := buf.String(), "INFO: info message\n"; got != want {
		t.Errorf("Info message was not logged correctly. Got %q, want %q", got, want)
	}
	buf.Reset()

	logging.Warn("warning message")
	if got, want := buf.String(), "WARN: warning message\n"; got != want {
		t.Errorf("Warning message was not logged correctly. Got %q, want %q", got, want)
	}
	buf.Reset()

	logging.Error("error message")
	if got, want := buf.String(), "ERRO: error message\n"; got != want {
		t.Errorf("Error message was not logged correctly. Got %q, want %q", got, want)
	}
	buf.Reset()
}
