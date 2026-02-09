package adapter_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/yuppyweb/cakelog/adapter"
)

type mockSlogHandler struct {
	contexts []context.Context
	records  []slog.Record
}

func (h *mockSlogHandler) Handle(ctx context.Context, record slog.Record) error {
	h.contexts = append(h.contexts, ctx)
	h.records = append(h.records, record)

	return nil
}

func (h *mockSlogHandler) WithAttrs([]slog.Attr) slog.Handler {
	return slog.DiscardHandler
}

func (h *mockSlogHandler) WithGroup(string) slog.Handler {
	return slog.DiscardHandler
}

func (h *mockSlogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

var _ slog.Handler = (*mockSlogHandler)(nil)

func TestSlogLogger_DebugWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "debugTestKey", "debug test value")
	values := []any{"key1", "debug", "key2", 42}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != adapter.DefaultSlogArgsKey {
			t.Errorf("expected attribute key '%s', got '%s'", adapter.DefaultSlogArgsKey, a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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

func TestSlogLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "infoTestKey", "info test value")
	values := []any{"key1", "info", "key2", 67}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != adapter.DefaultSlogArgsKey {
			t.Errorf("expected attribute key '%s', got '%s'", adapter.DefaultSlogArgsKey, a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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

func TestSlogLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "warnTestKey", "warn test value")
	values := []any{"key1", "warn", "key2", 89}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != adapter.DefaultSlogArgsKey {
			t.Errorf("expected attribute key '%s', got '%s'", adapter.DefaultSlogArgsKey, a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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

func TestSlogLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")
	values := []any{"key1", "error", "key2", 123}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != adapter.DefaultSlogArgsKey {
			t.Errorf("expected attribute key '%s', got '%s'", adapter.DefaultSlogArgsKey, a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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

func TestSlogLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))
	log.ArgsKey = "debugArgs"

	ctx := context.Background()
	ctx = context.WithValue(ctx, "debugTestKey", "debug test value")
	values := []any{"key1", "debug", "key2", 42}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "debugArgs" {
			t.Errorf("expected attribute key 'debugArgs', got '%s'", a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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

func TestSlogLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))
	log.ArgsKey = "infoArgs"

	ctx := context.Background()
	ctx = context.WithValue(ctx, "infoTestKey", "info test value")
	values := []any{"key1", "info", "key2", 67}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "infoArgs" {
			t.Errorf("expected attribute key 'infoArgs', got '%s'", a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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

func TestSlogLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))
	log.ArgsKey = "warnArgs"

	ctx := context.Background()
	ctx = context.WithValue(ctx, "warnTestKey", "warn test value")
	values := []any{"key1", "warn", "key2", 89}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "warnArgs" {
			t.Errorf("expected attribute key 'warnArgs', got '%s'", a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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

func TestSlogLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	handler := &mockSlogHandler{}
	log := adapter.NewSlogLogger(slog.New(handler))
	log.ArgsKey = "errorArgs"

	ctx := context.Background()
	ctx = context.WithValue(ctx, "errorTestKey", "error test value")
	values := []any{"key1", "error", "key2", 123}

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

	record.Attrs(func(a slog.Attr) bool {
		if a.Key != "errorArgs" {
			t.Errorf("expected attribute key 'errorArgs', got '%s'", a.Key)
		}

		if a.Value.String() != slog.AnyValue(values).String() {
			t.Errorf("expected attribute value '%s', got '%s'", slog.AnyValue(values).String(), a.Value.String())
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
