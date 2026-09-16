package integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog/decorator"
)

// TestDecoratorStack_InfoMasksContextAndCallSite checks the recommended stack
// Context → Mask → Callback. The masker sees unmasked context fields; the
// callback and underlying logger receive redacted data.
func TestDecoratorStack_InfoMasksContextAndCallSite(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)

	var gotMsg string

	var gotArgs []any

	base, logger := newRecommendedStack(t, masker, decorator.WithInfoCallback(
		func(_ context.Context, msg string, args ...any) {
			gotMsg = msg
			gotArgs = copyArgs(args)
		},
	))

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret, "user", "alice")

	if len(masker.messages) != 1 {
		t.Fatalf("expected masker to see 1 message, got %d", len(masker.messages))
	}

	if masker.messages[0] != "login secret" {
		t.Errorf(
			"expected masker message to be 'login secret', got %q",
			masker.messages[0],
		)
	}

	if len(masker.args) != 1 {
		t.Fatalf("expected masker to see 1 args slice, got %d", len(masker.args))
	}

	assertContextToken(t, masker.args[0], stackSecret)

	if len(masker.args[0]) != 5 {
		t.Fatalf("expected masker to receive 5 arguments, got %d", len(masker.args[0]))
	}

	if masker.args[0][1] != "password" || masker.args[0][2] != stackSecret {
		t.Errorf(
			"expected masker call-site password %q, got %v %v",
			stackSecret,
			masker.args[0][1],
			masker.args[0][2],
		)
	}

	if gotMsg != "login ***" {
		t.Errorf("expected callback message to be 'login ***', got %q", gotMsg)
	}

	assertContextToken(t, gotArgs, stackRedacted)

	if len(gotArgs) != 5 {
		t.Fatalf("expected callback to receive 5 arguments, got %d", len(gotArgs))
	}

	if gotArgs[1] != "password" || gotArgs[2] != stackRedacted {
		t.Errorf(
			"expected callback password to be %q, got %v %v",
			stackRedacted,
			gotArgs[1],
			gotArgs[2],
		)
	}

	if gotArgs[3] != "user" || gotArgs[4] != "alice" {
		t.Errorf("expected callback user to stay 'alice', got %v %v", gotArgs[3], gotArgs[4])
	}

	if len(base.infoIn) != 1 {
		t.Fatalf("expected Info to be called once, got %d", len(base.infoIn))
	}

	if base.infoIn[0].msg != "login ***" {
		t.Errorf("expected logger message to be 'login ***', got %q", base.infoIn[0].msg)
	}

	assertContextToken(t, base.infoIn[0].args, stackRedacted)

	if len(base.infoIn[0].args) != 5 {
		t.Fatalf("expected logger to receive 5 arguments, got %d", len(base.infoIn[0].args))
	}

	if base.infoIn[0].args[2] != stackRedacted {
		t.Errorf(
			"expected logger password to be %q, got %v",
			stackRedacted,
			base.infoIn[0].args[2],
		)
	}
}

// TestDecoratorStack_ErrorMasksBeforeCallback checks that Error redacts the
// error and context map before the callback and underlying logger see them.
func TestDecoratorStack_ErrorMasksBeforeCallback(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	expectedErr := errors.New("secret leaked")

	var gotErr error

	var gotArgs []any

	base, logger := newRecommendedStack(t, masker, decorator.WithErrorCallback(
		func(_ context.Context, err error, args ...any) {
			gotErr = err
			gotArgs = copyArgs(args)
		},
	))

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Error(ctx, expectedErr, "password", stackSecret)

	if len(masker.errs) != 1 {
		t.Fatalf("expected masker to see 1 error, got %d", len(masker.errs))
	}

	if !errors.Is(masker.errs[0], expectedErr) {
		t.Errorf(
			"expected masker error to be %v, got %v",
			expectedErr,
			masker.errs[0],
		)
	}

	assertContextToken(t, masker.args[0], stackSecret)

	if gotErr == nil {
		t.Fatal("expected callback to receive a masked error, got nil")
	}

	if gotErr.Error() != "*** leaked" {
		t.Errorf("expected callback error to be '*** leaked', got %q", gotErr.Error())
	}

	if !errors.Is(gotErr, expectedErr) {
		t.Errorf("expected callback error to unwrap to %v", expectedErr)
	}

	assertContextToken(t, gotArgs, stackRedacted)

	if len(gotArgs) != 3 {
		t.Fatalf("expected callback to receive 3 arguments, got %d", len(gotArgs))
	}

	if gotArgs[2] != stackRedacted {
		t.Errorf("expected callback password to be %q, got %v", stackRedacted, gotArgs[2])
	}

	if len(base.errorIn) != 1 {
		t.Fatalf("expected Error to be called once, got %d", len(base.errorIn))
	}

	if base.errorIn[0].err == nil || base.errorIn[0].err.Error() != "*** leaked" {
		t.Errorf(
			"expected logger error to be '*** leaked', got %v",
			base.errorIn[0].err,
		)
	}

	if !errors.Is(base.errorIn[0].err, expectedErr) {
		t.Errorf("expected logger error to unwrap to %v", expectedErr)
	}

	assertContextToken(t, base.errorIn[0].args, stackRedacted)
}

// TestDecoratorStack_ErrorNilForwardsNil checks that a nil error stays nil
// through MaskLogger and the error callback.
func TestDecoratorStack_ErrorNilForwardsNil(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)

	var gotErr error

	called := false

	base, logger := newRecommendedStack(t, masker, decorator.WithErrorCallback(
		func(_ context.Context, err error, _ ...any) {
			gotErr = err
			called = true
		},
	))

	ctx := withTestValue(context.Background(), "request_id", "req-1")
	logger.Error(ctx, nil)

	if !called {
		t.Fatal("expected error callback to run")
	}

	if gotErr != nil {
		t.Errorf("expected callback error to be nil, got %v", gotErr)
	}

	if len(masker.errs) != 1 || masker.errs[0] != nil {
		t.Errorf("expected masker to receive nil error, got %v", masker.errs)
	}

	if len(base.errorIn) != 1 {
		t.Fatalf("expected Error to be called once, got %d", len(base.errorIn))
	}

	if base.errorIn[0].err != nil {
		t.Errorf("expected logger error to be nil, got %v", base.errorIn[0].err)
	}

	loggerCtx := contextMapFromArgs(t, base.errorIn[0].args)
	if loggerCtx["request_id"] != "req-1" {
		t.Errorf(
			"expected logger context request_id to be 'req-1', got %v",
			loggerCtx["request_id"],
		)
	}
}

// TestDecoratorStack_AllLevelsForward checks that Debug, Info, Warn, and Error
// all go through the recommended stack with context values prepended.
func TestDecoratorStack_AllLevelsForward(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	expectedErr := errors.New("secret boom")

	base, logger := newRecommendedStack(t, masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Info(ctx, "info secret")
	logger.Warn(ctx, "warn secret")
	logger.Error(ctx, expectedErr)

	if len(base.debugIn) != 1 || base.debugIn[0].msg != "debug ***" {
		t.Errorf("expected Debug message 'debug ***', got %+v", base.debugIn)
	}

	if len(base.infoIn) != 1 || base.infoIn[0].msg != "info ***" {
		t.Errorf("expected Info message 'info ***', got %+v", base.infoIn)
	}

	if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn ***" {
		t.Errorf("expected Warn message 'warn ***', got %+v", base.warnIn)
	}

	if len(base.errorIn) != 1 ||
		base.errorIn[0].err == nil ||
		base.errorIn[0].err.Error() != "*** boom" {
		t.Errorf("expected Error '*** boom', got %+v", base.errorIn)
	}

	loggedArgs := [][]any{
		base.debugIn[0].args,
		base.infoIn[0].args,
		base.warnIn[0].args,
		base.errorIn[0].args,
	}

	for _, args := range loggedArgs {
		assertContextToken(t, args, stackRedacted)
	}

	if len(masker.messages) != 3 {
		t.Fatalf("expected 3 masked messages, got %d", len(masker.messages))
	}

	if len(masker.errs) != 1 {
		t.Fatalf("expected 1 masked error, got %d", len(masker.errs))
	}
}

// TestDecoratorStack_CallbackOutsideMaskSeesUnmasked checks that wrapping
// CallbackLogger outside MaskLogger leaks unmasked data to hooks.
func TestDecoratorStack_CallbackOutsideMaskSeesUnmasked(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	base := new(mockLogger)

	var gotMsg string

	var gotArgs []any

	maskLogger, err := decorator.NewMask(base, masker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
	}

	callbackLogger, err := decorator.NewCallback(
		maskLogger,
		decorator.WithInfoCallback(func(_ context.Context, msg string, args ...any) {
			gotMsg = msg
			gotArgs = copyArgs(args)
		}),
	)
	if err != nil {
		t.Fatalf("failed to create CallbackLogger: %v", err)
	}

	logger, err := decorator.NewContext(callbackLogger, testFields)
	if err != nil {
		t.Fatalf("failed to create context logger: %v", err)
	}

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret)

	if gotMsg != "login secret" {
		t.Errorf("expected callback message to stay unmasked, got %q", gotMsg)
	}

	assertContextToken(t, gotArgs, stackSecret)

	if len(gotArgs) != 3 || gotArgs[2] != stackSecret {
		t.Errorf(
			"expected callback password to stay %q, got %v",
			stackSecret,
			gotArgs,
		)
	}

	if len(base.infoIn) != 1 {
		t.Fatalf("expected Info to be called once, got %d", len(base.infoIn))
	}

	if base.infoIn[0].msg != "login ***" {
		t.Errorf("expected logger message to be masked, got %q", base.infoIn[0].msg)
	}

	assertContextToken(t, base.infoIn[0].args, stackRedacted)
}

// TestDecoratorStack_DebugAndWarnMaskBeforeCallback checks that Debug and
// Warn follow the same masking order as Info: hooks and the underlying
// logger receive redacted messages and context fields.
func TestDecoratorStack_DebugAndWarnMaskBeforeCallback(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)

	var gotDebugMsg, gotWarnMsg string

	var gotDebugArgs, gotWarnArgs []any

	base, logger := newRecommendedStack(
		t,
		masker,
		decorator.WithDebugCallback(func(_ context.Context, msg string, args ...any) {
			gotDebugMsg = msg
			gotDebugArgs = copyArgs(args)
		}),
		decorator.WithWarnCallback(func(_ context.Context, msg string, args ...any) {
			gotWarnMsg = msg
			gotWarnArgs = copyArgs(args)
		}),
	)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret", "password", stackSecret)
	logger.Warn(ctx, "warn secret", "password", stackSecret)

	if gotDebugMsg != "debug ***" {
		t.Errorf("expected debug callback message to be 'debug ***', got %q", gotDebugMsg)
	}

	if gotWarnMsg != "warn ***" {
		t.Errorf("expected warn callback message to be 'warn ***', got %q", gotWarnMsg)
	}

	assertContextToken(t, gotDebugArgs, stackRedacted)
	assertContextToken(t, gotWarnArgs, stackRedacted)

	if len(gotDebugArgs) != 3 || gotDebugArgs[2] != stackRedacted {
		t.Errorf("expected debug callback password to be %q, got %v", stackRedacted, gotDebugArgs)
	}

	if len(gotWarnArgs) != 3 || gotWarnArgs[2] != stackRedacted {
		t.Errorf("expected warn callback password to be %q, got %v", stackRedacted, gotWarnArgs)
	}

	if len(base.debugIn) != 1 || base.debugIn[0].msg != "debug ***" {
		t.Errorf("expected Debug message 'debug ***', got %+v", base.debugIn)
	}

	if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn ***" {
		t.Errorf("expected Warn message 'warn ***', got %+v", base.warnIn)
	}

	assertContextToken(t, base.debugIn[0].args, stackRedacted)
	assertContextToken(t, base.warnIn[0].args, stackRedacted)

	if len(masker.messages) != 2 {
		t.Fatalf("expected masker to see 2 messages, got %d", len(masker.messages))
	}

	if masker.messages[0] != "debug secret" || masker.messages[1] != "warn secret" {
		t.Errorf("expected masker to see unmasked messages, got %v", masker.messages)
	}
}

// TestDecoratorStack_EmptyContextMasksCallSiteOnly checks that a context
// with no extracted fields leaves args unchanged: no map is prepended, and
// call-site secrets are still masked before the callback and logger.
func TestDecoratorStack_EmptyContextMasksCallSiteOnly(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)

	var gotMsg string

	var gotArgs []any

	base, logger := newRecommendedStack(t, masker, decorator.WithInfoCallback(
		func(_ context.Context, msg string, args ...any) {
			gotMsg = msg
			gotArgs = copyArgs(args)
		},
	))

	logger.Info(context.Background(), "login secret", "password", stackSecret)

	if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
		t.Errorf("expected masker message 'login secret', got %v", masker.messages)
	}

	if len(masker.args) != 1 {
		t.Fatalf("expected masker to see 1 args slice, got %d", len(masker.args))
	}

	assertNoContextMap(t, masker.args[0])

	if len(masker.args[0]) != 2 || masker.args[0][1] != stackSecret {
		t.Errorf(
			"expected masker call-site password %q, got %v",
			stackSecret,
			masker.args[0],
		)
	}

	if gotMsg != "login ***" {
		t.Errorf("expected callback message to be 'login ***', got %q", gotMsg)
	}

	assertNoContextMap(t, gotArgs)

	if len(gotArgs) != 2 || gotArgs[1] != stackRedacted {
		t.Errorf("expected callback password to be %q, got %v", stackRedacted, gotArgs)
	}

	if len(base.infoIn) != 1 {
		t.Fatalf("expected Info to be called once, got %d", len(base.infoIn))
	}

	if base.infoIn[0].msg != "login ***" {
		t.Errorf("expected logger message to be 'login ***', got %q", base.infoIn[0].msg)
	}

	assertNoContextMap(t, base.infoIn[0].args)

	if len(base.infoIn[0].args) != 2 || base.infoIn[0].args[1] != stackRedacted {
		t.Errorf(
			"expected logger password to be %q, got %v",
			stackRedacted,
			base.infoIn[0].args,
		)
	}
}

// TestDecoratorStack_NestedContextMapIsMasked checks that a nested map in
// extracted context fields is walked by the masker, so hooks and the
// underlying logger receive the redacted inner value.
func TestDecoratorStack_NestedContextMapIsMasked(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)

	var gotArgs []any

	base, logger := newRecommendedStack(t, masker, decorator.WithInfoCallback(
		func(_ context.Context, _ string, args ...any) {
			gotArgs = copyArgs(args)
		},
	))

	ctx := withTestValue(context.Background(), "auth", map[string]any{
		"token": stackSecret,
	})
	logger.Info(ctx, "ok")

	if len(masker.args) != 1 {
		t.Fatalf("expected masker to see 1 args slice, got %d", len(masker.args))
	}

	assertNestedAuthToken(t, masker.args[0], stackSecret)
	assertNestedAuthToken(t, gotArgs, stackRedacted)

	if len(base.infoIn) != 1 {
		t.Fatalf("expected Info to be called once, got %d", len(base.infoIn))
	}

	assertNestedAuthToken(t, base.infoIn[0].args, stackRedacted)
}

// TestDecoratorStack_MaskOutsideContextLeaksContext checks that wrapping
// MaskLogger outside ContextLogger leaves context fields unmasked: the
// masker never sees the prepended map, so hooks and the logger receive
// the raw token while call-site args are still redacted.
func TestDecoratorStack_MaskOutsideContextLeaksContext(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	base := new(mockLogger)

	var gotMsg string

	var gotArgs []any

	callbackLogger, err := decorator.NewCallback(
		base,
		decorator.WithInfoCallback(func(_ context.Context, msg string, args ...any) {
			gotMsg = msg
			gotArgs = copyArgs(args)
		}),
	)
	if err != nil {
		t.Fatalf("failed to create CallbackLogger: %v", err)
	}

	contextLogger, err := decorator.NewContext(callbackLogger, testFields)
	if err != nil {
		t.Fatalf("failed to create context logger: %v", err)
	}

	logger, err := decorator.NewMask(contextLogger, masker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
	}

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret)

	if len(masker.args) != 1 {
		t.Fatalf("expected masker to see 1 args slice, got %d", len(masker.args))
	}

	assertNoContextMap(t, masker.args[0])

	if len(masker.args[0]) != 2 || masker.args[0][1] != stackSecret {
		t.Errorf(
			"expected masker to see call-site args only, got %v",
			masker.args[0],
		)
	}

	if gotMsg != "login ***" {
		t.Errorf("expected callback message to be masked, got %q", gotMsg)
	}

	assertContextToken(t, gotArgs, stackSecret)

	if len(gotArgs) != 3 || gotArgs[2] != stackRedacted {
		t.Errorf(
			"expected callback password to be %q, got %v",
			stackRedacted,
			gotArgs,
		)
	}

	if len(base.infoIn) != 1 {
		t.Fatalf("expected Info to be called once, got %d", len(base.infoIn))
	}

	if base.infoIn[0].msg != "login ***" {
		t.Errorf("expected logger message to be masked, got %q", base.infoIn[0].msg)
	}

	assertContextToken(t, base.infoIn[0].args, stackSecret)

	if len(base.infoIn[0].args) != 3 || base.infoIn[0].args[2] != stackRedacted {
		t.Errorf(
			"expected logger password to be %q, got %v",
			stackRedacted,
			base.infoIn[0].args,
		)
	}
}

// TestDecoratorStack_LevelOutsideDropsBeforeMaskAndCallback checks that a
// LevelLogger outside the recommended stack drops calls before Context,
// Mask, or Callback run.
func TestDecoratorStack_LevelOutsideDropsBeforeMaskAndCallback(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)

	debugCalled := false
	warnCalled := false

	base, stack := newRecommendedStack(
		t,
		masker,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			debugCalled = true
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			warnCalled = true
		}),
	)

	logger, err := decorator.NewLevel(stack, decorator.LevelWarn)
	if err != nil {
		t.Fatalf("failed to create level logger: %v", err)
	}

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Info(ctx, "info secret")
	logger.Warn(ctx, "warn secret")

	if debugCalled {
		t.Error("expected debug callback not to run for a dropped Debug call")
	}

	if !warnCalled {
		t.Error("expected warn callback to run")
	}

	if len(masker.messages) != 1 || masker.messages[0] != "warn secret" {
		t.Errorf("expected masker to see only Warn, got %v", masker.messages)
	}

	if len(base.debugIn) != 0 {
		t.Errorf("expected Debug to be dropped, got %+v", base.debugIn)
	}

	if len(base.infoIn) != 0 {
		t.Errorf("expected Info to be dropped, got %+v", base.infoIn)
	}

	if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn ***" {
		t.Errorf("expected Warn message 'warn ***', got %+v", base.warnIn)
	}

	assertContextToken(t, base.warnIn[0].args, stackRedacted)
}

// TestDecoratorStack_LevelInsideRunsMaskAndCallback checks that a
// LevelLogger inside the recommended stack still lets Mask and Callback
// run for dropped levels; only the underlying logger misses the call.
func TestDecoratorStack_LevelInsideRunsMaskAndCallback(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	base := new(mockLogger)

	debugCalled := false
	warnCalled := false

	inner, err := decorator.NewLevel(base, decorator.LevelWarn)
	if err != nil {
		t.Fatalf("failed to create level logger: %v", err)
	}

	logger := wrapRecommendedStack(
		t,
		inner,
		masker,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			debugCalled = true
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			warnCalled = true
		}),
	)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Warn(ctx, "warn secret")

	if !debugCalled {
		t.Error("expected debug callback to run even when Level drops Debug")
	}

	if !warnCalled {
		t.Error("expected warn callback to run")
	}

	if len(masker.messages) != 2 {
		t.Fatalf("expected masker to see Debug and Warn, got %d", len(masker.messages))
	}

	if masker.messages[0] != "debug secret" || masker.messages[1] != "warn secret" {
		t.Errorf("expected masker to see unmasked messages, got %v", masker.messages)
	}

	if len(base.debugIn) != 0 {
		t.Errorf("expected Debug to be dropped by Level, got %+v", base.debugIn)
	}

	if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn ***" {
		t.Errorf("expected Warn message 'warn ***', got %+v", base.warnIn)
	}

	assertContextToken(t, base.warnIn[0].args, stackRedacted)
}

// TestDecoratorStack_MuxReceivesMaskedContext checks that MuxLogger inside
// the recommended stack forwards the same masked message and context map
// to every underlying logger.
func TestDecoratorStack_MuxReceivesMaskedContext(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	first := new(mockLogger)
	second := new(mockLogger)

	mux, err := decorator.NewMux(first, second)
	if err != nil {
		t.Fatalf("failed to create mux logger: %v", err)
	}

	logger := wrapRecommendedStack(t, mux, masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret)

	for idx, base := range []*mockLogger{first, second} {
		if len(base.infoIn) != 1 {
			t.Fatalf("expected logger %d Info to be called once, got %d", idx, len(base.infoIn))
		}

		if base.infoIn[0].msg != "login ***" {
			t.Errorf(
				"expected logger %d message to be 'login ***', got %q",
				idx,
				base.infoIn[0].msg,
			)
		}

		assertContextToken(t, base.infoIn[0].args, stackRedacted)

		if len(base.infoIn[0].args) != 3 || base.infoIn[0].args[2] != stackRedacted {
			t.Errorf(
				"expected logger %d password to be %q, got %v",
				idx,
				stackRedacted,
				base.infoIn[0].args,
			)
		}
	}

	if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
		t.Errorf("expected masker to see 1 unmasked message, got %v", masker.messages)
	}
}
