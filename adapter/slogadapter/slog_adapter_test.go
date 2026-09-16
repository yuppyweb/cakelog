package slogadapter_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/slogadapter"
)

// mockSlogHandler is a custom slog.Handler implementation used for testing the SlogLogger adapter.
// It captures log records and contexts for verification in tests.
type mockSlogHandler struct {
	contexts []context.Context
	records  []slog.Record
}

// Handle captures the log record and context when a log message is handled through the SlogLogger,
// allowing tests to verify that the correct log level, message, and attributes are being used.
func (h *mockSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	h.contexts = append(h.contexts, ctx)
	h.records = append(h.records, record)

	return nil
}

// WithAttrs returns a new handler that discards attributes, as the mock handler
// does not need to handle attributes for testing purposes.
func (h *mockSlogHandler) WithAttrs([]slog.Attr) slog.Handler {
	return slog.DiscardHandler
}

// WithGroup returns a new handler that discards groups, as the mock handler
// does not need to handle groups for testing purposes.
func (h *mockSlogHandler) WithGroup(string) slog.Handler {
	return slog.DiscardHandler
}

// Enabled always returns true, indicating that all log levels are enabled for this mock handler.
func (h *mockSlogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

// Assert that mockSlogHandler implements the slog.Handler interface.
// This allows us to use it as a handler for testing the SlogLogger adapter.
var _ slog.Handler = (*mockSlogHandler)(nil)

// TestSlogLogger_Debug verifies that SlogLogger correctly logs debug messages.
func TestSlogLogger_Debug(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "debugTestKey", "debug test value")

	log.Debug(ctx, "debug message", "debug", 42)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	record := handler.records[0]

	if record.Message != "debug message" {
		t.Errorf("expected message 'debug message', got '%s'", record.Message)
	}

	if record.Level != slog.LevelDebug {
		t.Errorf("expected level Debug, got %s", record.Level)
	}

	attrs := slogAttrs(record)
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(attrs))
	}

	assertSlogAttr(t, attrs, "debug", 42)

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_Info verifies that SlogLogger correctly logs info messages.
func TestSlogLogger_Info(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "infoTestKey", "info test value")

	log.Info(ctx, "info message", "info", 67)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	record := handler.records[0]

	if record.Message != "info message" {
		t.Errorf("expected message 'info message', got '%s'", record.Message)
	}

	if record.Level != slog.LevelInfo {
		t.Errorf("expected level Info, got %s", record.Level)
	}

	attrs := slogAttrs(record)
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(attrs))
	}

	assertSlogAttr(t, attrs, "info", 67)

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_Warn verifies that SlogLogger correctly logs warn messages.
func TestSlogLogger_Warn(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "warnTestKey", "warn test value")

	log.Warn(ctx, "warn message", "warn", 89)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	record := handler.records[0]

	if record.Message != "warn message" {
		t.Errorf("expected message 'warn message', got '%s'", record.Message)
	}

	if record.Level != slog.LevelWarn {
		t.Errorf("expected level Warn, got %s", record.Level)
	}

	attrs := slogAttrs(record)
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(attrs))
	}

	assertSlogAttr(t, attrs, "warn", 89)

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_Error verifies that SlogLogger correctly logs error messages.
func TestSlogLogger_Error(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")

	log.Error(ctx, errors.New("error message"), "error", 123)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	record := handler.records[0]

	if record.Message != "error message" {
		t.Errorf("expected message 'error message', got '%s'", record.Message)
	}

	if record.Level != slog.LevelError {
		t.Errorf("expected level Error, got %s", record.Level)
	}

	attrs := slogAttrs(record)
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(attrs))
	}

	assertSlogAttr(t, attrs, "error", 123)

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_ErrorNil verifies that SlogLogger logs a nil error with an empty message.
func TestSlogLogger_ErrorNil(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorNilTestKey", "error nil test value")

	log.Error(ctx, nil, "error", 123)

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	record := handler.records[0]

	if record.Message != "" {
		t.Errorf("expected empty message, got '%s'", record.Message)
	}

	if record.Level != slog.LevelError {
		t.Errorf("expected level Error, got %s", record.Level)
	}

	attrs := slogAttrs(record)
	if len(attrs) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(attrs))
	}

	assertSlogAttr(t, attrs, "error", 123)

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_FlattenMap verifies that SlogLogger expands a map into individual attributes.
func TestSlogLogger_FlattenMap(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	log.Info(context.Background(), "mapped", map[string]any{"method": "GET", "status": 200})

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	attrs := slogAttrs(handler.records[0])
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(attrs))
	}

	assertSlogAttr(t, attrs, "method", "GET")
	assertSlogAttr(t, attrs, "status", 200)
}

// TestSlogLogger_NoArgs verifies that SlogLogger logs messages without extra arguments.
func TestSlogLogger_NoArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		call  func(cakelog.Logger)
		level slog.Level
		msg   string
	}{
		{
			name: "debug",
			call: func(logger cakelog.Logger) {
				logger.Debug(context.Background(), "debug message")
			},
			level: slog.LevelDebug,
			msg:   "debug message",
		},
		{
			name: "info",
			call: func(logger cakelog.Logger) {
				logger.Info(context.Background(), "info message")
			},
			level: slog.LevelInfo,
			msg:   "info message",
		},
		{
			name: "warn",
			call: func(logger cakelog.Logger) {
				logger.Warn(context.Background(), "warn message")
			},
			level: slog.LevelWarn,
			msg:   "warn message",
		},
		{
			name: "error",
			call: func(logger cakelog.Logger) {
				logger.Error(context.Background(), errors.New("error message"))
			},
			level: slog.LevelError,
			msg:   "error message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := new(mockSlogHandler)

			logger, err := slogadapter.New(slog.New(handler))
			if err != nil {
				t.Fatalf("failed to create SlogLogger: %v", err)
			}

			tc.call(logger)

			if len(handler.records) != 1 {
				t.Fatalf("expected 1 log record, got %d", len(handler.records))
			}

			record := handler.records[0]

			if record.Message != tc.msg {
				t.Errorf("expected message %q, got %q", tc.msg, record.Message)
			}

			if record.Level != tc.level {
				t.Errorf("expected level %s, got %s", tc.level, record.Level)
			}

			if record.NumAttrs() != 0 {
				t.Errorf("expected no attributes, got %v", slogAttrs(record))
			}
		})
	}
}

// TestNewSlogLogger_WithNilLogger verifies that NewSlogLogger returns an error when provided with a nil logger.
func TestNewSlogLogger_WithNilLogger(t *testing.T) {
	t.Parallel()

	_, err := slogadapter.New(nil)
	if err == nil {
		t.Fatal("expected error when creating SlogLogger with nil logger, but got nil")
	}

	if !errors.Is(err, slogadapter.ErrNilSlogLogger) {
		t.Fatalf(
			"unexpected error when creating SlogLogger with nil logger:\nGot:  %v\nWant: %v",
			err,
			slogadapter.ErrNilSlogLogger,
		)
	}
}

func slogAttrs(record slog.Record) map[string]slog.Value {
	attrs := make(map[string]slog.Value, record.NumAttrs())

	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value

		return true
	})

	return attrs
}

func assertSlogAttr(t *testing.T, attrs map[string]slog.Value, key string, want any) {
	t.Helper()

	got, ok := attrs[key]
	if !ok {
		t.Fatalf("expected attribute %q, got %v", key, attrs)
	}

	if got.String() != slog.AnyValue(want).String() {
		t.Errorf("expected attribute %q to be %v, got %s", key, want, got)
	}
}
