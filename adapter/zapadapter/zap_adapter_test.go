package zapadapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/zapadapter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// mockZapCore is a custom zapcore.Core implementation used for testing the zap adapter.
// It captures log entries and fields for verification in tests.
type mockZapCore struct {
	entry  zapcore.Entry
	fields []zapcore.Field
}

// Enabled always returns true, indicating that all log levels are enabled for this mock core.
func (c *mockZapCore) Enabled(zapcore.Level) bool {
	return true
}

// With appends the provided fields to the mock core's fields and returns the core itself for chaining.
func (c *mockZapCore) With(fields []zapcore.Field) zapcore.Core {
	c.fields = append(c.fields, fields...)

	return c
}

// Check adds the log entry to the checked entry if the log level is enabled,
// allowing the zap adapter to write log messages through this mock core.
func (c *mockZapCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

// Write captures the log entry and fields when a log message is written through the zap adapter,
// allowing tests to verify that the correct log level, message, and fields are being used.
func (c *mockZapCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	c.entry = ent
	c.fields = fields

	return nil
}

// Sync is a no-op for the mock core, as it does not perform any actual I/O operations.
func (c *mockZapCore) Sync() error {
	return nil
}

// TestAdapter_Debug verifies that the zap adapter correctly logs debug messages.
func TestAdapter_Debug(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	if mockCore.entry.Level != zap.DebugLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.DebugLevel)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "debug", 42)
}

// TestAdapter_Info verifies that the zap adapter correctly logs info messages.
func TestAdapter_Info(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 65)

	if mockCore.entry.Level != zap.InfoLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.InfoLevel)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "info", 65)
}

// TestAdapter_Warn verifies that the zap adapter correctly logs warn messages.
func TestAdapter_Warn(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 99)

	if mockCore.entry.Level != zap.WarnLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.WarnLevel)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "warn", 99)
}

// TestAdapter_Error verifies that the zap adapter correctly logs error messages.
func TestAdapter_Error(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	expectedErr := errors.New("error message")

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	logger.Error(context.Background(), expectedErr, "code", 123)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	if len(mockCore.fields) != 2 {
		t.Fatalf("unexpected number of fields: got %d, want 2", len(mockCore.fields))
	}

	assertZapError(t, mockCore.fields, expectedErr)
	assertZapField(t, mockCore.fields, "code", 123)
}

// TestAdapter_ErrorNil verifies that the zap adapter logs a nil error with an empty message.
func TestAdapter_ErrorNil(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	logger.Error(context.Background(), nil, "error", 123)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	if mockCore.entry.Message != "" {
		t.Errorf("unexpected log message: got %q, want empty", mockCore.entry.Message)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "error", 123)
}

// TestAdapter_FieldsMap verifies that the zap adapter expands a map into individual fields.
func TestAdapter_FieldsMap(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	logger.Info(context.Background(), "mapped", map[string]any{"method": "GET", "status": 200})

	if len(mockCore.fields) != 2 {
		t.Fatalf("unexpected number of fields: got %d, want 2", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "method", "GET")
	assertZapField(t, mockCore.fields, "status", 200)
}

// TestAdapter_NoArgs verifies that the zap adapter logs messages without extra arguments.
func TestAdapter_NoArgs(t *testing.T) {
	t.Parallel()

	loggedErr := errors.New("error message")
	testCases := []struct {
		name       string
		call       func(cakelog.Logger)
		level      zapcore.Level
		msg        string
		wantFields int
		errField   error
	}{
		{
			name: "debug",
			call: func(logger cakelog.Logger) {
				logger.Debug(context.Background(), "debug message")
			},
			level: zap.DebugLevel,
			msg:   "debug message",
		},
		{
			name: "info",
			call: func(logger cakelog.Logger) {
				logger.Info(context.Background(), "info message")
			},
			level: zap.InfoLevel,
			msg:   "info message",
		},
		{
			name: "warn",
			call: func(logger cakelog.Logger) {
				logger.Warn(context.Background(), "warn message")
			},
			level: zap.WarnLevel,
			msg:   "warn message",
		},
		{
			name: "error",
			call: func(logger cakelog.Logger) {
				logger.Error(context.Background(), loggedErr)
			},
			level:      zap.ErrorLevel,
			msg:        "error message",
			wantFields: 1,
			errField:   loggedErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockCore := new(mockZapCore)

			logger, err := zapadapter.New(zap.New(mockCore))
			if err != nil {
				t.Fatalf("failed to create zap adapter: %v", err)
			}

			tc.call(logger)

			if mockCore.entry.Level != tc.level {
				t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, tc.level)
			}

			if mockCore.entry.Message != tc.msg {
				t.Errorf("unexpected log message: got %q, want %q", mockCore.entry.Message, tc.msg)
			}

			if len(mockCore.fields) != tc.wantFields {
				t.Errorf("expected %d fields, got %v", tc.wantFields, mockCore.fields)
			}

			if tc.errField != nil {
				assertZapError(t, mockCore.fields, tc.errField)
			}
		})
	}
}

// TestNew_NilLogger verifies that New returns an error when provided with a nil logger.
func TestNew_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := zapadapter.New(nil)
	if err == nil {
		t.Fatal("expected error when creating zap adapter with nil logger, but got nil")
	}

	if !errors.Is(err, zapadapter.ErrNilZapLogger) {
		t.Errorf(
			"unexpected error when creating zap adapter with nil logger:\nGot:  %v\nWant: %v",
			err,
			zapadapter.ErrNilZapLogger,
		)
	}
}

// TestAdapter_DuplicateKeys verifies that zap keeps both values for a repeated key.
func TestAdapter_DuplicateKeys(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	logger.Info(context.Background(), "dup", "status", "ok", "status", "fail")

	if len(mockCore.fields) != 2 {
		t.Fatalf("unexpected number of fields: got %d, want 2", len(mockCore.fields))
	}

	if mockCore.fields[0].Key != "status" || mockCore.fields[1].Key != "status" {
		t.Fatalf("expected two status fields, got %v", mockCore.fields)
	}

	assertZapField(t, mockCore.fields[:1], "status", "ok")
	assertZapField(t, mockCore.fields[1:], "status", "fail")
}

func assertZapError(t *testing.T, fields []zapcore.Field, want error) {
	t.Helper()

	wantField := zap.Error(want)

	for _, field := range fields {
		if field.Key != "error" {
			continue
		}

		if !field.Equals(wantField) {
			t.Errorf("unexpected error field: got %#v, want %#v", field, wantField)
		}

		return
	}

	t.Errorf("expected zap.Error field, got %v", fields)
}

func assertZapField(t *testing.T, fields []zapcore.Field, key string, want any) {
	t.Helper()

	wantField := zap.Any(key, want)

	for _, field := range fields {
		if field.Key != key {
			continue
		}

		if !field.Equals(wantField) {
			t.Errorf("unexpected field %q: got %#v, want %#v", key, field, wantField)
		}

		return
	}

	t.Errorf("expected field %q, got %v", key, fields)
}
