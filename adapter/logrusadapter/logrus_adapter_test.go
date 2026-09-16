package logrusadapter_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/logrusadapter"
)

// TestLogrusLogger_Debug verifies that LogrusLogger correctly logs debug messages.
func TestLogrusLogger_Debug(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.DebugLevel))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=debug msg="debug message" debug=42` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_Info verifies that LogrusLogger correctly logs info messages.
func TestLogrusLogger_Info(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.InfoLevel))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 75)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=info msg="info message" info=75` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_Warn verifies that LogrusLogger correctly logs warn messages.
func TestLogrusLogger_Warn(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.WarnLevel))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 85)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=warning msg="warn message" warn=85` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_Error verifies that LogrusLogger correctly logs error messages.
func TestLogrusLogger_Error(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.ErrorLevel))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Error(context.Background(), errors.New("error message"), "error", 90)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=error msg="error message" error=90` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_ErrorNil verifies that LogrusLogger logs a nil error without panicking.
func TestLogrusLogger_ErrorNil(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.ErrorLevel))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Error(context.Background(), nil, "error", 90)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=error msg="<nil>" error=90` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_MapArgs verifies that LogrusLogger expands a map into individual fields.
func TestLogrusLogger_MapArgs(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.InfoLevel))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Info(context.Background(), "mapped", map[string]any{"method": "GET"})

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=info msg=mapped method=GET` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_NoArgs verifies that LogrusLogger logs messages without extra arguments.
func TestLogrusLogger_NoArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		level    logrus.Level
		call     func(cakelog.Logger)
		expected string
	}{
		{
			name:  "debug",
			level: logrus.DebugLevel,
			call: func(logger cakelog.Logger) {
				logger.Debug(context.Background(), "debug message")
			},
			expected: `level=debug msg="debug message"` + "\n",
		},
		{
			name:  "info",
			level: logrus.InfoLevel,
			call: func(logger cakelog.Logger) {
				logger.Info(context.Background(), "info message")
			},
			expected: `level=info msg="info message"` + "\n",
		},
		{
			name:  "warn",
			level: logrus.WarnLevel,
			call: func(logger cakelog.Logger) {
				logger.Warn(context.Background(), "warn message")
			},
			expected: `level=warning msg="warn message"` + "\n",
		},
		{
			name:  "error",
			level: logrus.ErrorLevel,
			call: func(logger cakelog.Logger) {
				logger.Error(context.Background(), errors.New("error message"))
			},
			expected: `level=error msg="error message"` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}

			logger, err := logrusadapter.New(newTestLogrus(buf, tc.level))
			if err != nil {
				t.Fatalf("failed to create LogrusLogger: %v", err)
			}

			tc.call(logger)

			output, err := io.ReadAll(buf)
			if err != nil {
				t.Fatalf("failed to read log output: %v", err)
			}

			if string(output) != tc.expected {
				t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), tc.expected)
			}
		})
	}
}

// TestNewLogrusLogger_WithNilLogger verifies that NewLogrusLogger returns an error when provided with a nil logger.
func TestNewLogrusLogger_WithNilLogger(t *testing.T) {
	t.Parallel()

	_, err := logrusadapter.New(nil)
	if err == nil {
		t.Fatal("expected error when creating LogrusLogger with nil logger, but got nil")
	}

	if !errors.Is(err, logrusadapter.ErrNilLogrusLogger) {
		t.Fatalf(
			"unexpected error when creating LogrusLogger with nil logger:\nGot:  %v\nWant: %v",
			err,
			logrusadapter.ErrNilLogrusLogger,
		)
	}
}

func newTestLogrus(buf *bytes.Buffer, level logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.SetOutput(buf)
	log.Level = level
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	return log
}
