package zerologadapter_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/zerologadapter"
)

// TestAdapter_Debug verifies that the zerolog adapter correctly logs debug messages.
func TestAdapter_Debug(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	assertZerologEntry(t, readZerologEntry(t, buf), "debug", "debug message", map[string]any{
		"debug": float64(42),
	})
}

// TestAdapter_Info verifies that the zerolog adapter correctly logs info messages.
func TestAdapter_Info(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 65)

	assertZerologEntry(t, readZerologEntry(t, buf), "info", "info message", map[string]any{
		"info": float64(65),
	})
}

// TestAdapter_Warn verifies that the zerolog adapter correctly logs warn messages.
func TestAdapter_Warn(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 80)

	assertZerologEntry(t, readZerologEntry(t, buf), "warn", "warn message", map[string]any{
		"warn": float64(80),
	})
}

// TestAdapter_Error verifies that the zerolog adapter correctly logs error messages.
func TestAdapter_Error(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	loggedErr := errors.New("test error")
	logger.Error(context.Background(), loggedErr, "code", 99)

	assertZerologEntry(t, readZerologEntry(t, buf), "error", "test error", map[string]any{
		"code":  float64(99),
		"error": loggedErr.Error(),
	})
}

// TestAdapter_ErrorNil verifies that the zerolog adapter logs a nil error with an empty message.
func TestAdapter_ErrorNil(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	logger.Error(context.Background(), nil, "error", 99)

	entry := readZerologEntry(t, buf)

	if entry["level"] != "error" {
		t.Errorf("unexpected level: got %v, want error", entry["level"])
	}

	if _, ok := entry["message"]; ok {
		t.Errorf("expected no message field for nil error, got %v", entry["message"])
	}

	got, ok := entry["error"]
	if !ok {
		t.Errorf("expected field %q, got %v", "error", entry)
	} else if got != float64(99) {
		t.Errorf(
			"unexpected field %q: got %v (%T), want %v (%T)",
			"error",
			got,
			got,
			float64(99),
			float64(99),
		)
	}
}

// TestAdapter_FieldsMap verifies that the zerolog adapter expands a map into individual fields.
func TestAdapter_FieldsMap(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	logger.Info(context.Background(), "mapped", map[string]any{"method": "GET", "status": 200})

	assertZerologEntry(t, readZerologEntry(t, buf), "info", "mapped", map[string]any{
		"method": "GET",
		"status": float64(200),
	})
}

// TestAdapter_NoArgs verifies that the zerolog adapter logs messages without extra arguments.
func TestAdapter_NoArgs(t *testing.T) {
	t.Parallel()

	loggedErr := errors.New("test error")
	testCases := []struct {
		name      string
		call      func(cakelog.Logger)
		level     string
		msg       string
		wantKeys  int
		errString string
	}{
		{
			name: "debug",
			call: func(logger cakelog.Logger) {
				logger.Debug(context.Background(), "debug message")
			},
			level: "debug",
			msg:   "debug message",
		},
		{
			name: "info",
			call: func(logger cakelog.Logger) {
				logger.Info(context.Background(), "info message")
			},
			level: "info",
			msg:   "info message",
		},
		{
			name: "warn",
			call: func(logger cakelog.Logger) {
				logger.Warn(context.Background(), "warn message")
			},
			level: "warn",
			msg:   "warn message",
		},
		{
			name: "error",
			call: func(logger cakelog.Logger) {
				logger.Error(context.Background(), loggedErr)
			},
			level:     "error",
			msg:       "test error",
			wantKeys:  3,
			errString: "test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			log := zerolog.New(buf)

			logger, err := zerologadapter.New(&log)
			if err != nil {
				t.Fatalf("failed to create zerolog adapter: %v", err)
			}

			tc.call(logger)

			entry := readZerologEntry(t, buf)
			assertZerologEntry(t, entry, tc.level, tc.msg, nil)

			if tc.wantKeys == 0 {
				if len(entry) != 2 {
					t.Errorf("expected only level and message, got %v", entry)
				}
			} else if len(entry) != tc.wantKeys {
				t.Errorf("expected %d keys, got %v", tc.wantKeys, entry)
			}

			if tc.errString != "" && entry["error"] != tc.errString {
				t.Errorf("expected error field %q, got %v", tc.errString, entry["error"])
			}
		})
	}
}

// TestNew_NilLogger verifies that New returns an error when provided with a nil logger.
func TestNew_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := zerologadapter.New(nil)
	if err == nil {
		t.Fatal("expected error when providing nil logger, got nil")
	}

	if !errors.Is(err, zerologadapter.ErrNilZerologLogger) {
		t.Errorf(
			"unexpected error when creating zerolog adapter with nil logger:\nGot:  %v\nWant: %v",
			err,
			zerologadapter.ErrNilZerologLogger,
		)
	}
}

// TestAdapter_DuplicateKeys verifies that zerolog keeps the last value for a repeated key.
func TestAdapter_DuplicateKeys(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	log := zerolog.New(buf)

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	logger.Info(context.Background(), "dup", "status", "ok", "status", "fail")

	assertZerologEntry(t, readZerologEntry(t, buf), "info", "dup", map[string]any{
		"status": "fail",
	})
}

func readZerologEntry(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	output, err := io.ReadAll(buf)
	if err != nil {
		t.Fatalf("failed to read log output: %v", err)
	}

	entry := make(map[string]any)
	if err := json.Unmarshal(output, &entry); err != nil {
		t.Fatalf("failed to unmarshal log output %q: %v", string(output), err)
	}

	return entry
}

func assertZerologEntry(
	t *testing.T,
	entry map[string]any,
	level,
	message string,
	fields map[string]any,
) {
	t.Helper()

	if entry["level"] != level {
		t.Errorf("unexpected level: got %v, want %s", entry["level"], level)
	}

	if entry["message"] != message {
		t.Errorf("unexpected message: got %v, want %s", entry["message"], message)
	}

	for key, want := range fields {
		got, ok := entry[key]
		if !ok {
			t.Errorf("expected field %q, got %v", key, entry)

			continue
		}

		if got != want {
			t.Errorf("unexpected field %q: got %v (%T), want %v (%T)", key, got, got, want, want)
		}
	}
}
