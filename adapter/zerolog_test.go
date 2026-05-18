package adapter_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/yuppyweb/cakelog/adapter"
)

// TestZerologLogger_DebugWithDefaultArgsKey verifies that ZerologLogger correctly
// logs debug messages with default args key.
func TestZerologLogger_DebugWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log)
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"debug","%s":[%q,%d],"message":"debug message"}`+"\n",
		"context", "debug", 42,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestZerologLogger_InfoWithDefaultArgsKey verifies that ZerologLogger correctly
// logs info messages with default args key.
func TestZerologLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log)
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 65)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"info","%s":[%q,%d],"message":"info message"}`+"\n",
		"context", "info", 65,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestZerologLogger_WarnWithDefaultArgsKey verifies that ZerologLogger correctly
// logs warn messages with default args key.
func TestZerologLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log)
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 80)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"warn","%s":[%q,%d],"message":"warn message"}`+"\n",
		"context", "warn", 80,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestZerologLogger_ErrorWithDefaultArgsKey verifies that ZerologLogger correctly
// logs error messages with default args key.
func TestZerologLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log)
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logErr := errors.New("test error")
	logger.Error(context.Background(), logErr, "error", 99)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"error","%s":[%q,%d],"message":"test error"}`+"\n",
		"context", "error", 99,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestZerologLogger_DebugWithCustomArgsKey verifies that ZerologLogger correctly
// logs debug messages with custom args key.
func TestZerologLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log, adapter.WithArgsKey("customArgs1"))
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"debug","%s":[%q,%d],"message":"debug message"}`+"\n",
		"customArgs1", "debug", 42,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestZerologLogger_InfoWithCustomArgsKey verifies that ZerologLogger correctly
// logs info messages with custom args key.
func TestZerologLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log, adapter.WithArgsKey("customArgs2"))
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 65)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"info","%s":[%q,%d],"message":"info message"}`+"\n",
		"customArgs2", "info", 65,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestZerologLogger_WarnWithCustomArgsKey verifies that ZerologLogger correctly
// logs warn messages with custom args key.
func TestZerologLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log, adapter.WithArgsKey("customArgs3"))
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 80)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"warn","%s":[%q,%d],"message":"warn message"}`+"\n",
		"customArgs3", "warn", 80,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestZerologLogger_ErrorWithCustomArgsKey verifies that ZerologLogger correctly
// logs error messages with custom args key.
func TestZerologLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := adapter.NewZerologLogger(&log, adapter.WithArgsKey("customArgs4"))
	if err != nil {
		t.Fatalf("failed to create ZerologLogger: %v", err)
	}

	logErr := errors.New("test error")
	logger.Error(context.Background(), logErr, "error", 99)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"error","%s":[%q,%d],"message":"test error"}`+"\n",
		"customArgs4", "error", 99,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

// TestNewZerologLogger_WithNilLogger verifies that NewZerologLogger returns an error when provided with a nil logger.
func TestNewZerologLogger_WithNilLogger(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewZerologLogger(nil)
	if err == nil {
		t.Fatal("expected error when providing nil logger, got nil")
	}

	if !errors.Is(err, adapter.ErrNilZerologLogger) {
		t.Errorf("expected error to be ErrNilZerologLogger, got %v", err)
	}
}

// TestNewZerologLogger_WithNilOption verifies that NewZerologLogger returns an error when provided with a nil option.
func TestNewZerologLogger_WithNilOption(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewZerologLogger(&zerolog.Logger{}, nil)
	if err == nil {
		t.Fatal("expected error when providing nil option, got nil")
	}

	if !errors.Is(err, adapter.ErrNilZerologOption) {
		t.Errorf("expected error to be ErrNilZerologOption, got %v", err)
	}
}

// TestNewZerologLogger_WithInvalidArgsKey verifies that NewZerologLogger returns
// an error when provided with an empty args key.
func TestNewZerologLogger_WithInvalidArgsKey(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewZerologLogger(&zerolog.Logger{}, adapter.WithArgsKey(""))
	if err == nil {
		t.Fatal("expected error when providing empty args key, got nil")
	}

	if !errors.Is(err, adapter.ErrEmptyArgsKey) {
		t.Errorf("expected error to be ErrEmptyArgsKey, got %v", err)
	}

	if !strings.Contains(err.Error(), "failed to apply option:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"failed to apply option:",
		)
	}
}
