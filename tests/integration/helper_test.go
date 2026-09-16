package integration_test

import (
	"context"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
)

// mockMsgArgs holds context, message, and arguments for a message log call.
type mockMsgArgs struct {
	ctx  context.Context
	msg  string
	args []any
}

// mockErrArgs holds context, error, and arguments for an error log call.
type mockErrArgs struct {
	ctx  context.Context
	err  error
	args []any
}

// mockLogger is a mock implementation of cakelog.Logger for testing.
type mockLogger struct {
	debugIn []mockMsgArgs
	infoIn  []mockMsgArgs
	warnIn  []mockMsgArgs
	errorIn []mockErrArgs
}

// Debug records a debug-level message with context and arguments.
func (ml *mockLogger) Debug(ctx context.Context, msg string, args ...any) {
	ml.debugIn = append(ml.debugIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

// Info records an info-level message with context and arguments.
func (ml *mockLogger) Info(ctx context.Context, msg string, args ...any) {
	ml.infoIn = append(ml.infoIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

// Warn records a warning-level message with context and arguments.
func (ml *mockLogger) Warn(ctx context.Context, msg string, args ...any) {
	ml.warnIn = append(ml.warnIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

// Error records an error-level message with context, error, and arguments.
func (ml *mockLogger) Error(ctx context.Context, err error, args ...any) {
	ml.errorIn = append(ml.errorIn, mockErrArgs{
		ctx:  ctx,
		err:  err,
		args: args,
	})
}

// Ensure mockLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*mockLogger)(nil)

const (
	stackSecret   = "secret"
	stackRedacted = "***"
)

// stackMaskedError redacts Error() while preserving errors.Is / errors.As.
type stackMaskedError struct {
	err error
	msg string
}

func (e stackMaskedError) Error() string { return e.msg }
func (e stackMaskedError) Unwrap() error { return e.err }

// stackMasker redacts stackSecret in messages, errors, strings, and maps.
// It records the unmasked values it received so stack tests can assert order.
type stackMasker struct {
	messages []string
	errs     []error
	args     [][]any
}

func (m *stackMasker) MaskMessage(msg string) string {
	m.messages = append(m.messages, msg)

	return strings.ReplaceAll(msg, stackSecret, stackRedacted)
}

func (m *stackMasker) MaskError(err error) error {
	m.errs = append(m.errs, err)

	if err == nil {
		return nil
	}

	return stackMaskedError{
		err: err,
		msg: strings.ReplaceAll(err.Error(), stackSecret, stackRedacted),
	}
}

func (m *stackMasker) MaskArguments(args ...any) []any {
	copied := make([]any, len(args))
	copy(copied, args)
	m.args = append(m.args, copied)

	return maskStackArgs(args)
}

var _ decorator.Masker = (*stackMasker)(nil)

func maskStackArgs(args []any) []any {
	masked := make([]any, len(args))
	for i, arg := range args {
		masked[i] = maskStackValue(arg)
	}

	return masked
}

func maskStackValue(arg any) any {
	switch val := arg.(type) {
	case string:
		return strings.ReplaceAll(val, stackSecret, stackRedacted)
	case map[string]any:
		out := make(map[string]any, len(val))
		for key, value := range val {
			out[key] = maskStackValue(value)
		}

		return out
	default:
		return arg
	}
}

type testCtxKey struct{}

func withTestValue(ctx context.Context, key string, value any) context.Context {
	parent, _ := ctx.Value(testCtxKey{}).(map[string]any)

	next := make(map[string]any, len(parent)+1)
	for existingKey, existingValue := range parent {
		next[existingKey] = existingValue
	}

	next[key] = value

	return context.WithValue(ctx, testCtxKey{}, next)
}

func testFields(ctx context.Context) map[string]any {
	fields, _ := ctx.Value(testCtxKey{}).(map[string]any)

	return fields
}

func wrapRecommendedStack(
	t *testing.T,
	inner cakelog.Logger,
	masker decorator.Masker,
	opts ...decorator.CallbackOption,
) cakelog.Logger {
	t.Helper()

	callbackLogger, err := decorator.NewCallback(inner, opts...)
	if err != nil {
		t.Fatalf("failed to create CallbackLogger: %v", err)
	}

	maskLogger, err := decorator.NewMask(callbackLogger, masker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
	}

	logger, err := decorator.NewContext(maskLogger, testFields)
	if err != nil {
		t.Fatalf("failed to create context logger: %v", err)
	}

	return logger
}

func newRecommendedStack(
	t *testing.T,
	masker decorator.Masker,
	opts ...decorator.CallbackOption,
) (*mockLogger, cakelog.Logger) {
	t.Helper()

	base := new(mockLogger)

	return base, wrapRecommendedStack(t, base, masker, opts...)
}

func wrapContextMask(
	t *testing.T,
	inner cakelog.Logger,
	masker decorator.Masker,
) cakelog.Logger {
	t.Helper()

	maskLogger, err := decorator.NewMask(inner, masker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
	}

	logger, err := decorator.NewContext(maskLogger, testFields)
	if err != nil {
		t.Fatalf("failed to create context logger: %v", err)
	}

	return logger
}

func wrapRecommendedLevelDebug(
	t *testing.T,
	inner cakelog.Logger,
	masker decorator.Masker,
	opts ...decorator.CallbackOption,
) cakelog.Logger {
	t.Helper()

	stack := wrapRecommendedStack(t, inner, masker, opts...)

	logger, err := decorator.NewLevel(stack, decorator.LevelDebug)
	if err != nil {
		t.Fatalf("failed to create level logger: %v", err)
	}

	return logger
}

func wrapRecommendedLevelInfo(
	t *testing.T,
	inner cakelog.Logger,
	masker decorator.Masker,
	opts ...decorator.CallbackOption,
) cakelog.Logger {
	t.Helper()

	stack := wrapRecommendedStack(t, inner, masker, opts...)

	logger, err := decorator.NewLevel(stack, decorator.LevelInfo)
	if err != nil {
		t.Fatalf("failed to create level logger: %v", err)
	}

	return logger
}

func newMuxLogger(t *testing.T, logs ...cakelog.Logger) cakelog.Logger {
	t.Helper()

	logger, err := decorator.NewMux(logs...)
	if err != nil {
		t.Fatalf("failed to create mux logger: %v", err)
	}

	return logger
}

func contextMapFromArgs(t *testing.T, args []any) map[string]any {
	t.Helper()

	if len(args) == 0 {
		t.Fatal("expected args to include a context map")
	}

	values, ok := args[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first argument to be map[string]any, got %T", args[0])
	}

	return values
}

func assertContextToken(t *testing.T, args []any, want any) {
	t.Helper()

	values := contextMapFromArgs(t, args)
	if values["token"] != want {
		t.Errorf("expected context token to be %q, got %v", want, values["token"])
	}
}

func assertNestedAuthToken(t *testing.T, args []any, want any) {
	t.Helper()

	values := contextMapFromArgs(t, args)

	auth, ok := values["auth"].(map[string]any)
	if !ok {
		t.Fatalf("expected context auth to be map[string]any, got %T", values["auth"])
	}

	if auth["token"] != want {
		t.Errorf("expected nested context token to be %q, got %v", want, auth["token"])
	}
}

func assertNoContextMap(t *testing.T, args []any) {
	t.Helper()

	if len(args) == 0 {
		return
	}

	if _, ok := args[0].(map[string]any); ok {
		t.Errorf("expected no context map, got %v", args[0])
	}
}

func copyArgs(args []any) []any {
	copied := make([]any, len(args))
	copy(copied, args)

	return copied
}
