package adapter_test

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog/adapter"
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

// TestSlogLogger_DebugWithDefaultArgsKey verifies that SlogLogger correctly logs debug messages with default args key.
func TestSlogLogger_DebugWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "debugTestKey", "debug test value")
	values := []any{"debug", 42}

	log.Debug(ctx, "debug message", values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "context" {
			t.Errorf("expected attribute key 'context', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_InfoWithDefaultArgsKey verifies that SlogLogger correctly logs info messages with default args key.
func TestSlogLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "infoTestKey", "info test value")
	values := []any{"info", 67}

	log.Info(ctx, "info message", values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "context" {
			t.Errorf("expected attribute key 'context', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_WarnWithDefaultArgsKey verifies that SlogLogger correctly logs warn messages with default args key.
func TestSlogLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "warnTestKey", "warn test value")
	values := []any{"warn", 89}

	log.Warn(ctx, "warn message", values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "context" {
			t.Errorf("expected attribute key 'context', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_ErrorWithDefaultArgsKey verifies that SlogLogger correctly logs error messages with default args key.
func TestSlogLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")
	values := []any{"error", 123}

	log.Error(ctx, errors.New("error message"), values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "context" {
			t.Errorf("expected attribute key 'context', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_DebugWithCustomArgsKey verifies that SlogLogger correctly logs debug messages with custom args key.
func TestSlogLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler), adapter.WithArgsKey("debugArgs"))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "debugTestKey", "debug test value")
	values := []any{"debug", 42}

	log.Debug(ctx, "debug message", values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "debugArgs" {
			t.Errorf("expected attribute key 'debugArgs', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_InfoWithCustomArgsKey verifies that SlogLogger correctly logs info messages with custom args key.
func TestSlogLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler), adapter.WithArgsKey("infoArgs"))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "infoTestKey", "info test value")
	values := []any{"info", 67}

	log.Info(ctx, "info message", values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "infoArgs" {
			t.Errorf("expected attribute key 'infoArgs', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_WarnWithCustomArgsKey verifies that SlogLogger correctly logs warn messages with custom args key.
func TestSlogLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler), adapter.WithArgsKey("warnArgs"))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "warnTestKey", "warn test value")
	values := []any{"warn", 89}

	log.Warn(ctx, "warn message", values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "warnArgs" {
			t.Errorf("expected attribute key 'warnArgs', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestSlogLogger_ErrorWithCustomArgsKey verifies that SlogLogger correctly logs error messages with custom args key.
func TestSlogLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := new(mockSlogHandler)

	log, err := adapter.NewSlogLogger(slog.New(handler), adapter.WithArgsKey("errorArgs"))
	if err != nil {
		t.Fatalf("failed to create SlogLogger: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")
	values := []any{"error", 123}

	log.Error(ctx, errors.New("error message"), values...)

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

	if record.NumAttrs() != 1 {
		t.Fatalf("expected 1 attribute, got %d", record.NumAttrs())
	}

	checkedAttr := false

	record.Attrs(func(attr slog.Attr) bool {
		if attr.Key != "errorArgs" {
			t.Errorf("expected attribute key 'errorArgs', got '%s'", attr.Key)
		}

		if attr.Value.String() != slog.AnyValue(values).String() {
			t.Errorf(
				"expected attribute value '%s', got '%s'",
				slog.AnyValue(values).String(),
				attr.Value.String(),
			)
		}

		checkedAttr = true

		return true
	})

	if !checkedAttr {
		t.Fatal("expected to check the log record attributes, but did not")
	}

	if len(handler.contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(handler.contexts))
	}

	if handler.contexts[0] != ctx {
		t.Errorf("expected context %v, got %v", ctx, handler.contexts[0])
	}
}

// TestNewSlogLogger_WithNilLogger verifies that NewSlogLogger returns an error when provided with a nil logger.
func TestNewSlogLogger_WithNilLogger(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewSlogLogger(nil)
	if err == nil {
		t.Fatal("expected error when creating SlogLogger with nil logger, but got nil")
	}

	if !errors.Is(err, adapter.ErrNilSlogLogger) {
		t.Fatalf(
			"unexpected error when creating SlogLogger with nil logger:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrNilSlogLogger,
		)
	}
}

// TestNewSlogLogger_WithNilOption verifies that NewSlogLogger returns an error when provided with a nil option.
func TestNewSlogLogger_WithNilOption(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewSlogLogger(slog.New(new(mockSlogHandler)), nil)
	if err == nil {
		t.Fatal("expected error when creating SlogLogger with nil option, but got nil")
	}

	if !errors.Is(err, adapter.ErrNilSlogOption) {
		t.Fatalf(
			"unexpected error when creating SlogLogger with nil option:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrNilSlogOption,
		)
	}
}

// TestNewSlogLogger_WithInvalidArgsKey verifies that NewSlogLogger returns
// an error when provided with an empty args key.
func TestNewSlogLogger_WithInvalidArgsKey(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewSlogLogger(slog.New(&mockSlogHandler{}), adapter.WithArgsKey(""))
	if err == nil {
		t.Fatal("expected error when creating SlogLogger with empty args key, but got nil")
	}

	if !errors.Is(err, adapter.ErrEmptyArgsKey) {
		t.Fatalf(
			"unexpected error when creating SlogLogger with empty args key:\nGot:  %v\nWant: %v",
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
