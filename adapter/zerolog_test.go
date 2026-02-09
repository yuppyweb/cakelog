package adapter_test

import (
	"bytes"
	"context"
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
	ctx := context.Background()
	ctx = context.WithValue(ctx, "debugTestKey", "debug test value")

	logger := adapter.NewZerologLogger(&log)

	logger.Debug(ctx, "debug message", "key1", "debug", "key2", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"debug","%s":[%q,%q,%q,%d],"message":"debug message"}`+"\n",
		adapter.DefaultZerologArgsKey,
		"key1", "debug", "key2", 42,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "infoTestKey", "info test value")

	logger := adapter.NewZerologLogger(&log)

	logger.Info(ctx, "info message", "key1", "info", "key2", 65)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"info","%s":[%q,%q,%q,%d],"message":"info message"}`+"\n",
		adapter.DefaultZerologArgsKey,
		"key1", "info", "key2", 65,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "warnTestKey", "warn test value")

	logger := adapter.NewZerologLogger(&log)

	logger.Warn(ctx, "warn message", "key1", "warn", "key2", 80)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"warn","%s":[%q,%q,%q,%d],"message":"warn message"}`+"\n",
		adapter.DefaultZerologArgsKey,
		"key1", "warn", "key2", 80,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")

	logger := adapter.NewZerologLogger(&log)

	err := fmt.Errorf("test error")
	logger.Error(ctx, err, "key1", "error", "key2", 99)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"error","%s":[%q,%q,%q,%d],"message":"test error"}`+"\n",
		adapter.DefaultZerologArgsKey,
		"key1", "error", "key2", 99,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "debugTestKey", "debug test value")

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs"

	logger.Debug(ctx, "debug message", "key1", "debug", "key2", 42)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"debug","%s":[%q,%q,%q,%d],"message":"debug message"}`+"\n",
		"customArgs",
		"key1", "debug", "key2", 42,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "infoTestKey", "info test value")

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs"

	logger.Info(ctx, "info message", "key1", "info", "key2", 65)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"info","%s":[%q,%q,%q,%d],"message":"info message"}`+"\n",
		"customArgs",
		"key1", "info", "key2", 65,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "warnTestKey", "warn test value")

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs"

	logger.Warn(ctx, "warn message", "key1", "warn", "key2", 80)

	output, err := io.ReadAll(buf)

	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"warn","%s":[%q,%q,%q,%d],"message":"warn message"}`+"\n",
		"customArgs",
		"key1", "warn", "key2", 80,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}

func TestZerologLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")

	logger := adapter.NewZerologLogger(&log)
	logger.ArgsKey = "customArgs"

	err := fmt.Errorf("test error")
	logger.Error(ctx, err, "key1", "error", "key2", 99)

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	expected := fmt.Sprintf(`{"level":"error","%s":[%q,%q,%q,%d],"message":"test error"}`+"\n",
		"customArgs",
		"key1", "error", "key2", 99,
	)

	if string(output) != expected {
		t.Errorf("unexpected log output:\nGot:  %s\nWant: %s", string(output), expected)
	}
}
