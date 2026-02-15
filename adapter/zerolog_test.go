package adapter_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/yuppyweb/cakelog/adapter"
)

func TestZerologLogger_DebugWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)

	logger.Debug(context.Background(), "debug message", "debug", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"debug","%s":[%q,%d],"message":"debug message"}`+"\n",
		adapter.DefaultZerologArgsKey, "debug", 42,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)

	logger.Info(context.Background(), "info message", "info", 65)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"info","%s":[%q,%d],"message":"info message"}`+"\n",
		adapter.DefaultZerologArgsKey, "info", 65,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)

	logger.Warn(context.Background(), "warn message", "warn", 80)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"warn","%s":[%q,%d],"message":"warn message"}`+"\n",
		adapter.DefaultZerologArgsKey, "warn", 80,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)

	err := errors.New("test error")
	logger.Error(context.Background(), err, "error", 99)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(
		`{"level":"error","%s":[%q,%d],"message":"test error"}`+"\n",
		adapter.DefaultZerologArgsKey, "error", 99,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs1"

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

func TestZerologLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs2"

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

func TestZerologLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs3"

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

func TestZerologLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs4"

	err := errors.New("test error")
	logger.Error(context.Background(), err, "error", 99)

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
