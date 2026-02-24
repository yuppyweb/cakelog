package adapter_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog/adapter"
)

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

	logger := adapter.NewLogrusLogger(log)

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

	logger := adapter.NewLogrusLogger(log)

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

	logger := adapter.NewLogrusLogger(log)

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

	logger := adapter.NewLogrusLogger(log)

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

	logger := adapter.NewLogrusLogger(log)
	logger.ArgsKey = "custom_args1"

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

	logger := adapter.NewLogrusLogger(log)
	logger.ArgsKey = "custom_args2"

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

	logger := adapter.NewLogrusLogger(log)
	logger.ArgsKey = "custom_args3"

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

	logger := adapter.NewLogrusLogger(log)
	logger.ArgsKey = "custom_args4"

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
