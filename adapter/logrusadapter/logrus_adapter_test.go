package logrusadapter_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/logrusadapter"
)

// TestAdapter_Debug verifies that the logrus adapter correctly logs debug messages.
func TestAdapter_Debug(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.DebugLevel))
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
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

// TestAdapter_Info verifies that the logrus adapter correctly logs info messages.
func TestAdapter_Info(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.InfoLevel))
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
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

// TestAdapter_Warn verifies that the logrus adapter correctly logs warn messages.
func TestAdapter_Warn(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.WarnLevel))
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
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

// TestAdapter_Error verifies that the logrus adapter correctly logs error messages.
func TestAdapter_Error(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.ErrorLevel))
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
	}

	logger.Error(context.Background(), errors.New("error message"), "code", 90)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	got := string(output)
	for _, want := range []string{
		`level=error`,
		`msg="error message"`,
		`error="error message"`,
		`code=90`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected log output to contain %q, got %s", want, got)
		}
	}
}

// TestAdapter_ErrorNil verifies that the logrus adapter logs a nil error without panicking.
func TestAdapter_ErrorNil(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.ErrorLevel))
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
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

// TestAdapter_MapArgs verifies that the logrus adapter expands a map into individual fields.
func TestAdapter_MapArgs(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.InfoLevel))
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
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

// TestAdapter_NoArgs verifies that the logrus adapter logs messages without extra arguments.
func TestAdapter_NoArgs(t *testing.T) {
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
			expected: `level=error msg="error message" error="error message"` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}

			logger, err := logrusadapter.New(newTestLogrus(buf, tc.level))
			if err != nil {
				t.Fatalf("failed to create logrus adapter: %v", err)
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

// TestNew_NilLogger verifies that New returns ErrNilLogrusLogger for a nil
// FieldLogger, including typed nil *logrus.Logger and *logrus.Entry values.
func TestNew_NilLogger(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		logger logrus.FieldLogger
	}{
		{name: "untyped", logger: nil},
		{name: "typed_logger", logger: (*logrus.Logger)(nil)},
		{name: "typed_entry", logger: (*logrus.Entry)(nil)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := logrusadapter.New(tc.logger)
			if err == nil {
				t.Fatal("expected error when creating logrus adapter with nil logger, but got nil")
			}

			if !errors.Is(err, logrusadapter.ErrNilLogrusLogger) {
				t.Fatalf(
					"unexpected error when creating logrus adapter with nil logger:\nGot:  %v\nWant: %v",
					err,
					logrusadapter.ErrNilLogrusLogger,
				)
			}
		})
	}
}

// TestAdapter_Entry verifies that New accepts a *logrus.Entry and keeps
// fields already attached to that entry.
func TestAdapter_Entry(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	entry := newTestLogrus(buf, logrus.InfoLevel).WithField("service", "api")

	logger, err := logrusadapter.New(entry)
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
	}

	logger.Info(context.Background(), "hello", "status", 200)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	got := string(output)
	for _, want := range []string{
		`level=info`,
		`msg=hello`,
		`service=api`,
		`status=200`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected log output to contain %q, got %s", want, got)
		}
	}
}

// TestAdapter_DuplicateKeys verifies that logrus keeps the last value for a repeated key.
func TestAdapter_DuplicateKeys(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	logger, err := logrusadapter.New(newTestLogrus(buf, logrus.InfoLevel))
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
	}

	logger.Info(context.Background(), "dup", "status", "ok", "status", "fail")

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=info msg=dup status=fail` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}

	if strings.Contains(string(output), "status=ok") {
		t.Errorf("expected last-wins status=fail only, got %s", string(output))
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
