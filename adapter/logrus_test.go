package adapter_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog/adapter"
)

// TestLogrusLogger_DebugWithDefaultArgsKey verifies that LogrusLogger correctly
// logs debug messages with default args key.
func TestLogrusLogger_DebugWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.DebugLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log)
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=debug msg="debug message" context="[debug 42]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_InfoWithDefaultArgsKey verifies that LogrusLogger correctly
// logs info messages with default args key.
func TestLogrusLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.InfoLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log)
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 75)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=info msg="info message" context="[info 75]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_WarnWithDefaultArgsKey verifies that LogrusLogger correctly
// logs warn messages with default args key.
func TestLogrusLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.WarnLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log)
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 85)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=warning msg="warn message" context="[warn 85]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_ErrorWithDefaultArgsKey verifies that LogrusLogger correctly
// logs error messages with default args key.
func TestLogrusLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.ErrorLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log)
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	expectedErr := errors.New("error message")
	logger.Error(context.Background(), expectedErr, "error", 90)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=error msg="error message" context="[error 90]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_DebugWithCustomArgsKey verifies that LogrusLogger correctly
// logs debug messages with custom args key.
func TestLogrusLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.DebugLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log, adapter.WithArgsKey("custom_args1"))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=debug msg="debug message" custom_args1="[debug 42]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_InfoWithCustomArgsKey verifies that LogrusLogger correctly logs info messages with custom args key.
func TestLogrusLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.InfoLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log, adapter.WithArgsKey("custom_args2"))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 75)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=info msg="info message" custom_args2="[info 75]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_WarnWithCustomArgsKey verifies that LogrusLogger correctly logs warn messages with custom args key.
func TestLogrusLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.WarnLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log, adapter.WithArgsKey("custom_args3"))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 85)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=warning msg="warn message" custom_args3="[warn 85]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestLogrusLogger_ErrorWithCustomArgsKey verifies that LogrusLogger correctly
// logs error messages with custom args key.
func TestLogrusLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.ErrorLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	logger, err := adapter.NewLogrusLogger(log, adapter.WithArgsKey("custom_args4"))
	if err != nil {
		t.Fatalf("failed to create LogrusLogger: %v", err)
	}

	expectedErr := errors.New("error message")
	logger.Error(context.Background(), expectedErr, "error", 90)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := `level=error msg="error message" custom_args4="[error 90]"` + "\n"

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestNewLogrusLogger_WithNilLogger verifies that NewLogrusLogger returns an error when provided with a nil logger.
func TestNewLogrusLogger_WithNilLogger(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewLogrusLogger(nil)
	if err == nil {
		t.Fatal("expected error when creating LogrusLogger with nil logger, but got nil")
	}

	if !errors.Is(err, adapter.ErrNilLogrusLogger) {
		t.Fatalf(
			"unexpected error when creating LogrusLogger with nil logger:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrNilLogrusLogger,
		)
	}
}

// TestNewLogrusLogger_WithNilOption verifies that NewLogrusLogger returns an error when provided with a nil option.
func TestNewLogrusLogger_WithNilOption(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewLogrusLogger(logrus.New(), nil)
	if err == nil {
		t.Fatal("expected error when creating LogrusLogger with nil option, but got nil")
	}

	if !errors.Is(err, adapter.ErrNilLogrusOption) {
		t.Fatalf(
			"unexpected error when creating LogrusLogger with nil option:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrNilLogrusOption,
		)
	}
}

// TestNewLogrusLogger_WithInvalidArgsKey verifies that NewLogrusLogger returns
// an error when provided with an empty args key.
func TestNewLogrusLogger_WithInvalidArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := logrus.New()

	log.SetOutput(buf)
	log.Level = logrus.DebugLevel
	log.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		DisableColors:    true,
	})

	_, err := adapter.NewLogrusLogger(log, adapter.WithArgsKey(""))
	if err == nil {
		t.Fatal("expected error when creating LogrusLogger with empty args key, but got nil")
	}

	if !errors.Is(err, adapter.ErrEmptyArgsKey) {
		t.Fatalf(
			"unexpected error when creating LogrusLogger with empty args key:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrEmptyArgsKey,
		)
	}

	if !strings.Contains(err.Error(), "failed to apply option:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"failed to apply option:",
		)
	}
}
