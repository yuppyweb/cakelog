package adapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog/adapter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mockZapCore struct {
	entry  zapcore.Entry
	fields []zapcore.Field
}

func (c *mockZapCore) Enabled(zapcore.Level) bool {
	return true
}

func (c *mockZapCore) With(fields []zapcore.Field) zapcore.Core {
	c.fields = append(c.fields, fields...)

	return c
}

func (c *mockZapCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *mockZapCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	c.entry = ent
	c.fields = fields

	return nil
}

func (c *mockZapCore) Sync() error {
	return nil
}

func TestZapLogger_DebugWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"debug", 42}

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)

	logger.Debug(context.Background(), "debug message", values...)

	if mockCore.entry.Level != zap.DebugLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.DebugLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(adapter.DefaultZapArgsKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

func TestZapLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"info", 65}

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)

	logger.Info(context.Background(), "info message", values...)

	if mockCore.entry.Level != zap.InfoLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.InfoLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(adapter.DefaultZapArgsKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

func TestZapLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"warn", 99}

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)

	logger.Warn(context.Background(), "warn message", values...)

	if mockCore.entry.Level != zap.WarnLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.WarnLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(adapter.DefaultZapArgsKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

func TestZapLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"error", 123}
	expectedErr := errors.New("error message")

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)

	logger.Error(context.Background(), expectedErr, values...)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(adapter.DefaultZapArgsKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}
}

func TestZapLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"debug", 42}
	customKey := "custom_debug"

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)
	logger.ArgsKey = customKey

	logger.Debug(context.Background(), "debug message", values...)

	if mockCore.entry.Level != zap.DebugLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.DebugLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

func TestZapLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"info", 65}
	customKey := "custom_info"

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)
	logger.ArgsKey = customKey

	logger.Info(context.Background(), "info message", values...)

	if mockCore.entry.Level != zap.InfoLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.InfoLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

func TestZapLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"warn", 99}
	customKey := "custom_warn"

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)
	logger.ArgsKey = customKey

	logger.Warn(context.Background(), "warn message", values...)

	if mockCore.entry.Level != zap.WarnLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.WarnLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

func TestZapLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"error", 123}
	customKey := "custom_error"
	expectedErr := errors.New("error message")

	log := zap.New(mockCore)
	logger := adapter.NewZapLogger(log)
	logger.ArgsKey = customKey

	logger.Error(context.Background(), expectedErr, values...)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}
