package slogadapter_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/slogadapter"
)

// mockSlogHandler is a custom slog.Handler implementation used for testing the slog adapter.
// It captures log records and contexts for verification in tests.
type mockSlogHandler struct {
	contexts []context.Context
	records  []slog.Record
}

// Handle captures the log record and context when a log message is handled through the slog adapter,
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
// This allows us to use it as a handler for testing the slog adapter.
var _ slog.Handler = (*mockSlogHandler)(nil)

// TestAdapter_Debug verifies that the slog adapter correctly logs debug messages.
func TestAdapter_Debug(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
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

// TestAdapter_Info verifies that the slog adapter correctly logs info messages.
func TestAdapter_Info(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
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

// TestAdapter_Warn verifies that the slog adapter correctly logs warn messages.
func TestAdapter_Warn(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
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

// TestAdapter_Error verifies that the slog adapter correctly logs error messages.
func TestAdapter_Error(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")

	loggedErr := errors.New("error message")
	log.Error(ctx, loggedErr, "code", 123)

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
	if len(attrs) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(attrs))
	}

	assertSlogAttr(t, attrs, "error", loggedErr)
	assertSlogAttr(t, attrs, "code", 123)

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestAdapter_ErrorNil verifies that the slog adapter logs a nil error with an empty message.
func TestAdapter_ErrorNil(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
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

// TestAdapter_FlattenMap verifies that the slog adapter expands a map into individual attributes.
func TestAdapter_FlattenMap(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
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

// TestAdapter_NoArgs verifies that the slog adapter logs messages without extra arguments.
func TestAdapter_NoArgs(t *testing.T) {
	t.Parallel()

	loggedErr := errors.New("error message")
	testCases := []struct {
		name      string
		call      func(cakelog.Logger)
		level     slog.Level
		msg       string
		wantAttrs int
		errAttr   error
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
				logger.Error(context.Background(), loggedErr)
			},
			level:     slog.LevelError,
			msg:       "error message",
			wantAttrs: 1,
			errAttr:   loggedErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			handler := new(mockSlogHandler)

			logger, err := slogadapter.New(slog.New(handler))
			if err != nil {
				t.Fatalf("failed to create slog adapter: %v", err)
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

			attrs := slogAttrs(record)
			if record.NumAttrs() != tc.wantAttrs {
				t.Errorf("expected %d attributes, got %v", tc.wantAttrs, attrs)
			}

			if tc.errAttr != nil {
				assertSlogAttr(t, attrs, "error", tc.errAttr)
			}
		})
	}
}

// TestNew_NilLogger verifies that New returns an error when provided with a nil logger.
func TestNew_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := slogadapter.New(nil)
	if err == nil {
		t.Fatal("expected error when creating slog adapter with nil logger, but got nil")
	}

	if !errors.Is(err, slogadapter.ErrNilSlogLogger) {
		t.Fatalf(
			"unexpected error when creating slog adapter with nil logger:\nGot:  %v\nWant: %v",
			err,
			slogadapter.ErrNilSlogLogger,
		)
	}
}

// TestAdapter_DuplicateKeys verifies that slog keeps both values for a repeated key.
func TestAdapter_DuplicateKeys(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
	}

	log.Info(context.Background(), "dup", "status", "ok", "status", "fail")

	if len(handler.records) != 1 {
		t.Fatalf("expected 1 log record, got %d", len(handler.records))
	}

	got := slogAttrList(handler.records[0])
	want := []slog.Attr{
		slog.String("status", "ok"),
		slog.String("status", "fail"),
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d attributes, got %v", len(want), got)
	}

	for i := range want {
		if got[i].Key != want[i].Key || got[i].Value.String() != want[i].Value.String() {
			t.Errorf("unexpected attribute %d: got %#v, want %#v", i, got[i], want[i])
		}
	}
}

func slogAttrList(record slog.Record) []slog.Attr {
	attrs := make([]slog.Attr, 0, record.NumAttrs())

	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)

		return true
	})

	return attrs
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
