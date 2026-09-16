package integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
)

func newLevelLogger(t *testing.T, log cakelog.Logger, minLevel int) cakelog.Logger {
	t.Helper()

	var logger cakelog.Logger

	var err error

	switch minLevel {
	case int(decorator.LevelDebug):
		logger, err = decorator.NewLevel(log, decorator.LevelDebug)
	case int(decorator.LevelInfo):
		logger, err = decorator.NewLevel(log, decorator.LevelInfo)
	case int(decorator.LevelWarn):
		logger, err = decorator.NewLevel(log, decorator.LevelWarn)
	case int(decorator.LevelError):
		logger, err = decorator.NewLevel(log, decorator.LevelError)
	default:
		t.Fatalf("unsupported min level %d", minLevel)
	}

	if err != nil {
		t.Fatalf("failed to create level logger: %v", err)
	}

	return logger
}

// TestDecoratorStack_LevelInfoOutsideDropsDebug checks that LevelInfo outside
// the recommended stack drops Debug before Mask and Callback, while Info,
// Warn, and Error still reach the base logger masked.
func TestDecoratorStack_LevelInfoOutsideDropsDebug(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)

	debugCalled := false
	infoCalled := false

	base, stack := newRecommendedStack(
		t,
		masker,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			debugCalled = true
		}),
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			infoCalled = true
		}),
	)

	logger := newLevelLogger(t, stack, int(decorator.LevelInfo))

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Info(ctx, "info secret")
	logger.Warn(ctx, "warn secret")
	logger.Error(ctx, errors.New("secret boom"))

	if debugCalled {
		t.Error("expected debug callback not to run for a dropped Debug call")
	}

	if !infoCalled {
		t.Error("expected info callback to run")
	}

	if len(masker.messages) != 2 {
		t.Fatalf("expected masker to see Info and Warn, got %v", masker.messages)
	}

	if masker.messages[0] != "info secret" || masker.messages[1] != "warn secret" {
		t.Errorf("expected unmasked Info and Warn, got %v", masker.messages)
	}

	if len(masker.errs) != 1 || masker.errs[0] == nil || masker.errs[0].Error() != "secret boom" {
		t.Errorf("expected masker to see unmasked Error, got %v", masker.errs)
	}

	if len(base.debugIn) != 0 {
		t.Errorf("expected Debug to be dropped, got %+v", base.debugIn)
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

	assertContextToken(t, base.infoIn[0].args, stackRedacted)
	assertContextToken(t, base.warnIn[0].args, stackRedacted)
	assertContextToken(t, base.errorIn[0].args, stackRedacted)
}

// TestDecoratorStack_LevelErrorOutsideForwardsOnlyError checks that LevelError
// outside the recommended stack drops Debug, Info, and Warn before Mask and
// Callback, while Error is still masked and forwarded.
func TestDecoratorStack_LevelErrorOutsideForwardsOnlyError(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	expectedErr := errors.New("secret leaked")

	warnCalled := false
	errorCalled := false

	base, stack := newRecommendedStack(
		t,
		masker,
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			warnCalled = true
		}),
		decorator.WithErrorCallback(func(context.Context, error, ...any) {
			errorCalled = true
		}),
	)

	logger := newLevelLogger(t, stack, int(decorator.LevelError))

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Info(ctx, "info secret")
	logger.Warn(ctx, "warn secret")
	logger.Error(ctx, expectedErr)

	if warnCalled {
		t.Error("expected warn callback not to run for a dropped Warn call")
	}

	if !errorCalled {
		t.Error("expected error callback to run")
	}

	if len(masker.messages) != 0 {
		t.Errorf("expected masker not to see dropped messages, got %v", masker.messages)
	}

	if len(masker.errs) != 1 || !errors.Is(masker.errs[0], expectedErr) {
		t.Errorf("expected masker to see the original error, got %v", masker.errs)
	}

	if len(base.debugIn) != 0 || len(base.infoIn) != 0 || len(base.warnIn) != 0 {
		t.Errorf(
			"expected only Error to reach the base logger, got debug=%+v info=%+v warn=%+v",
			base.debugIn,
			base.infoIn,
			base.warnIn,
		)
	}

	if len(base.errorIn) != 1 ||
		base.errorIn[0].err == nil ||
		base.errorIn[0].err.Error() != "*** leaked" {
		t.Errorf("expected Error '*** leaked', got %+v", base.errorIn)
	}

	if !errors.Is(base.errorIn[0].err, expectedErr) {
		t.Errorf("expected logger error to unwrap to %v", expectedErr)
	}

	assertContextToken(t, base.errorIn[0].args, stackRedacted)
}

// TestDecoratorStack_LevelBetweenContextAndMask checks that Level between
// Context and Mask still extracts context fields for dropped calls, but does
// not run Mask or Callback.
func TestDecoratorStack_LevelBetweenContextAndMask(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	base := new(mockLogger)

	debugCalled := false
	warnCalled := false
	fieldsCalls := 0

	callbackLogger, err := decorator.NewCallback(
		base,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			debugCalled = true
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			warnCalled = true
		}),
	)
	if err != nil {
		t.Fatalf("failed to create CallbackLogger: %v", err)
	}

	maskLogger, err := decorator.NewMask(callbackLogger, masker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
	}

	levelLogger := newLevelLogger(t, maskLogger, int(decorator.LevelWarn))

	logger, err := decorator.NewContext(levelLogger, func(ctx context.Context) map[string]any {
		fieldsCalls++

		return testFields(ctx)
	})
	if err != nil {
		t.Fatalf("failed to create context logger: %v", err)
	}

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Warn(ctx, "warn secret")

	if fieldsCalls != 2 {
		t.Errorf("expected context fields to run for Debug and Warn, got %d", fieldsCalls)
	}

	if debugCalled {
		t.Error("expected debug callback not to run")
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

	if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn ***" {
		t.Errorf("expected Warn message 'warn ***', got %+v", base.warnIn)
	}

	assertContextToken(t, base.warnIn[0].args, stackRedacted)
}

// TestDecoratorStack_LevelBetweenMaskAndCallback checks that Level between
// Mask and Callback still lets the masker see dropped Debug calls, while the
// callback and base logger do not.
func TestDecoratorStack_LevelBetweenMaskAndCallback(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	base := new(mockLogger)

	debugCalled := false
	warnCalled := false

	callbackLogger, err := decorator.NewCallback(
		base,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			debugCalled = true
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			warnCalled = true
		}),
	)
	if err != nil {
		t.Fatalf("failed to create CallbackLogger: %v", err)
	}

	levelLogger := newLevelLogger(t, callbackLogger, int(decorator.LevelWarn))
	logger := wrapContextMask(t, levelLogger, masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Warn(ctx, "warn secret")

	if debugCalled {
		t.Error("expected debug callback not to run")
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
		t.Errorf("expected Debug to be dropped, got %+v", base.debugIn)
	}

	if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn ***" {
		t.Errorf("expected Warn message 'warn ***', got %+v", base.warnIn)
	}

	assertContextToken(t, base.warnIn[0].args, stackRedacted)
}

// TestDecoratorStack_MuxOutsideIndependentStacks checks that Mux wrapping two
// recommended stacks runs each masker and callback once and delivers the same
// masked payload to both bases.
func TestDecoratorStack_MuxOutsideIndependentStacks(t *testing.T) {
	t.Parallel()

	maskerA := new(stackMasker)
	maskerB := new(stackMasker)

	callsA := 0
	callsB := 0

	baseA, stackA := newRecommendedStack(t, maskerA, decorator.WithInfoCallback(
		func(context.Context, string, ...any) {
			callsA++
		},
	))
	baseB, stackB := newRecommendedStack(t, maskerB, decorator.WithInfoCallback(
		func(context.Context, string, ...any) {
			callsB++
		},
	))

	logger := newMuxLogger(t, stackA, stackB)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret)

	if callsA != 1 || callsB != 1 {
		t.Errorf("expected one callback per stack, got A=%d B=%d", callsA, callsB)
	}

	if len(maskerA.messages) != 1 || maskerA.messages[0] != "login secret" {
		t.Errorf("expected masker A to see 1 unmasked message, got %v", maskerA.messages)
	}

	if len(maskerB.messages) != 1 || maskerB.messages[0] != "login secret" {
		t.Errorf("expected masker B to see 1 unmasked message, got %v", maskerB.messages)
	}

	for idx, base := range []*mockLogger{baseA, baseB} {
		if len(base.infoIn) != 1 {
			t.Fatalf("expected logger %d Info to be called once, got %d", idx, len(base.infoIn))
		}

		if base.infoIn[0].msg != "login ***" {
			t.Errorf("expected logger %d message 'login ***', got %q", idx, base.infoIn[0].msg)
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
}

// TestDecoratorStack_MuxAsymmetricBranches checks that Mux does not mask on
// its own: a Context→Mask branch is redacted while a bare logger still sees
// the raw secret.
func TestDecoratorStack_MuxAsymmetricBranches(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	masked := new(mockLogger)
	raw := new(mockLogger)

	logger := newMuxLogger(t, wrapContextMask(t, masked, masker), raw)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret)

	if len(masked.infoIn) != 1 || masked.infoIn[0].msg != "login ***" {
		t.Errorf("expected masked branch message 'login ***', got %+v", masked.infoIn)
	}

	assertContextToken(t, masked.infoIn[0].args, stackRedacted)

	if len(masked.infoIn[0].args) != 3 || masked.infoIn[0].args[2] != stackRedacted {
		t.Errorf("expected masked branch password %q, got %v", stackRedacted, masked.infoIn[0].args)
	}

	if len(raw.infoIn) != 1 || raw.infoIn[0].msg != "login secret" {
		t.Errorf("expected raw branch message 'login secret', got %+v", raw.infoIn)
	}

	assertNoContextMap(t, raw.infoIn[0].args)

	if len(raw.infoIn[0].args) != 2 || raw.infoIn[0].args[1] != stackSecret {
		t.Errorf("expected raw branch password %q, got %v", stackSecret, raw.infoIn[0].args)
	}
}

// TestDecoratorStack_CallbackOnceOutsideMux checks that a Callback wrapping
// Mux fires once while both underlying loggers still receive the call.
func TestDecoratorStack_CallbackOnceOutsideMux(t *testing.T) {
	t.Parallel()

	first := new(mockLogger)
	second := new(mockLogger)
	calls := 0

	callbackLogger, err := decorator.NewCallback(
		newMuxLogger(t, first, second),
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			calls++
		}),
	)
	if err != nil {
		t.Fatalf("failed to create CallbackLogger: %v", err)
	}

	callbackLogger.Info(context.Background(), "hello", "user", "alice")

	if calls != 1 {
		t.Errorf("expected callback to run once, got %d", calls)
	}

	for idx, base := range []*mockLogger{first, second} {
		if len(base.infoIn) != 1 || base.infoIn[0].msg != "hello" {
			t.Errorf("expected logger %d Info 'hello', got %+v", idx, base.infoIn)
		}
	}
}

// TestDecoratorStack_CallbackPerMuxBranch checks that Mux of two Callback
// loggers invokes each hook once.
func TestDecoratorStack_CallbackPerMuxBranch(t *testing.T) {
	t.Parallel()

	first := new(mockLogger)
	second := new(mockLogger)
	callsA := 0
	callsB := 0

	callbackA, err := decorator.NewCallback(
		first,
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			callsA++
		}),
	)
	if err != nil {
		t.Fatalf("failed to create first CallbackLogger: %v", err)
	}

	callbackB, err := decorator.NewCallback(
		second,
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			callsB++
		}),
	)
	if err != nil {
		t.Fatalf("failed to create second CallbackLogger: %v", err)
	}

	logger := newMuxLogger(t, callbackA, callbackB)
	logger.Info(context.Background(), "hello")

	if callsA != 1 || callsB != 1 {
		t.Errorf("expected one callback per branch, got A=%d B=%d", callsA, callsB)
	}

	if len(first.infoIn) != 1 || len(second.infoIn) != 1 {
		t.Errorf(
			"expected both loggers to receive Info, got first=%d second=%d",
			len(first.infoIn),
			len(second.infoIn),
		)
	}
}

// TestDecoratorStack_MuxErrorMasksBoth checks that Error through a recommended
// stack with Mux inside delivers the masked error to every underlying logger.
func TestDecoratorStack_MuxErrorMasksBoth(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	first := new(mockLogger)
	second := new(mockLogger)
	expectedErr := errors.New("secret leaked")

	logger := wrapRecommendedStack(t, newMuxLogger(t, first, second), masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Error(ctx, expectedErr, "password", stackSecret)

	if len(masker.errs) != 1 || !errors.Is(masker.errs[0], expectedErr) {
		t.Errorf("expected masker to see 1 unmasked error, got %v", masker.errs)
	}

	for idx, base := range []*mockLogger{first, second} {
		if len(base.errorIn) != 1 {
			t.Fatalf("expected logger %d Error to be called once, got %d", idx, len(base.errorIn))
		}

		if base.errorIn[0].err == nil || base.errorIn[0].err.Error() != "*** leaked" {
			t.Errorf("expected logger %d error '*** leaked', got %v", idx, base.errorIn[0].err)
		}

		if !errors.Is(base.errorIn[0].err, expectedErr) {
			t.Errorf("expected logger %d error to unwrap to %v", idx, expectedErr)
		}

		assertContextToken(t, base.errorIn[0].args, stackRedacted)

		if len(base.errorIn[0].args) != 3 || base.errorIn[0].args[2] != stackRedacted {
			t.Errorf(
				"expected logger %d password to be %q, got %v",
				idx,
				stackRedacted,
				base.errorIn[0].args,
			)
		}
	}
}

// TestDecoratorStack_MuxNilErrorForwardsNil checks that a nil Error through Mux
// stays nil on every underlying logger.
func TestDecoratorStack_MuxNilErrorForwardsNil(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	first := new(mockLogger)
	second := new(mockLogger)

	logger := wrapRecommendedStack(t, newMuxLogger(t, first, second), masker)

	ctx := withTestValue(context.Background(), "request_id", "req-1")
	logger.Error(ctx, nil)

	if len(masker.errs) != 1 || masker.errs[0] != nil {
		t.Errorf("expected masker to receive nil error, got %v", masker.errs)
	}

	for idx, base := range []*mockLogger{first, second} {
		if len(base.errorIn) != 1 {
			t.Fatalf("expected logger %d Error to be called once, got %d", idx, len(base.errorIn))
		}

		if base.errorIn[0].err != nil {
			t.Errorf("expected logger %d error to be nil, got %v", idx, base.errorIn[0].err)
		}

		values := contextMapFromArgs(t, base.errorIn[0].args)
		if values["request_id"] != "req-1" {
			t.Errorf("expected logger %d request_id 'req-1', got %v", idx, values["request_id"])
		}
	}
}

// TestDecoratorStack_MuxPerBranchLevel checks that Mux can wrap Level
// loggers with different thresholds under a shared Context→Mask stack.
func TestDecoratorStack_MuxPerBranchLevel(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	first := new(mockLogger)
	second := new(mockLogger)

	mux := newMuxLogger(
		t,
		newLevelLogger(t, first, int(decorator.LevelWarn)),
		newLevelLogger(t, second, int(decorator.LevelError)),
	)
	logger := wrapContextMask(t, mux, masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Warn(ctx, "warn secret")
	logger.Error(ctx, errors.New("secret boom"))

	if len(masker.messages) != 1 || masker.messages[0] != "warn secret" {
		t.Errorf("expected masker to see Warn, got %v", masker.messages)
	}

	if len(masker.errs) != 1 {
		t.Errorf("expected masker to see Error, got %v", masker.errs)
	}

	if len(first.warnIn) != 1 || first.warnIn[0].msg != "warn ***" {
		t.Errorf("expected first logger Warn 'warn ***', got %+v", first.warnIn)
	}

	assertContextToken(t, first.warnIn[0].args, stackRedacted)

	if len(second.warnIn) != 0 {
		t.Errorf("expected second logger to drop Warn, got %+v", second.warnIn)
	}

	for idx, base := range []*mockLogger{first, second} {
		if len(base.errorIn) != 1 ||
			base.errorIn[0].err == nil ||
			base.errorIn[0].err.Error() != "*** boom" {
			t.Errorf("expected logger %d Error '*** boom', got %+v", idx, base.errorIn)
		}

		assertContextToken(t, base.errorIn[0].args, stackRedacted)
	}
}

// TestDecoratorStack_LevelOutsideMuxDropsBoth checks that Level wrapping Mux
// drops Debug before either branch sees it.
func TestDecoratorStack_LevelOutsideMuxDropsBoth(t *testing.T) {
	t.Parallel()

	first := new(mockLogger)
	second := new(mockLogger)

	logger := newLevelLogger(t, newMuxLogger(t, first, second), int(decorator.LevelWarn))

	logger.Debug(context.Background(), "debug secret")
	logger.Warn(context.Background(), "warn secret")

	for idx, base := range []*mockLogger{first, second} {
		if len(base.debugIn) != 0 {
			t.Errorf("expected logger %d Debug to be dropped, got %+v", idx, base.debugIn)
		}

		if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn secret" {
			t.Errorf("expected logger %d Warn 'warn secret', got %+v", idx, base.warnIn)
		}
	}
}

// TestDecoratorStack_MuxWithInnerLevelDropsOneBranch checks that Mux wrapping
// Level on only one branch still forwards Debug to the unfiltered logger.
func TestDecoratorStack_MuxWithInnerLevelDropsOneBranch(t *testing.T) {
	t.Parallel()

	first := new(mockLogger)
	second := new(mockLogger)

	logger := newMuxLogger(t, newLevelLogger(t, first, int(decorator.LevelWarn)), second)

	logger.Debug(context.Background(), "debug secret")
	logger.Warn(context.Background(), "warn secret")

	if len(first.debugIn) != 0 {
		t.Errorf("expected first logger Debug to be dropped, got %+v", first.debugIn)
	}

	if len(second.debugIn) != 1 || second.debugIn[0].msg != "debug secret" {
		t.Errorf("expected second logger Debug 'debug secret', got %+v", second.debugIn)
	}

	for idx, base := range []*mockLogger{first, second} {
		if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn secret" {
			t.Errorf("expected logger %d Warn 'warn secret', got %+v", idx, base.warnIn)
		}
	}
}

// TestDecoratorStack_MuxWithNopLogger checks that Mux with NopLogger still
// writes to the real branch and does not panic.
func TestDecoratorStack_MuxWithNopLogger(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	logger := newMuxLogger(t, base, cakelog.NopLogger())

	logger.Info(context.Background(), "hello", "user", "alice")
	logger.Error(context.Background(), errors.New("boom"))

	if len(base.infoIn) != 1 || base.infoIn[0].msg != "hello" {
		t.Errorf("expected Info 'hello', got %+v", base.infoIn)
	}

	if len(base.errorIn) != 1 ||
		base.errorIn[0].err == nil ||
		base.errorIn[0].err.Error() != "boom" {
		t.Errorf("expected Error 'boom', got %+v", base.errorIn)
	}
}

// TestDecoratorStack_FullStackLevelMux checks the production stack
// Level → Context → Mask → Callback → Mux: Debug is dropped before Mask,
// while Warn and Error reach both bases masked.
func TestDecoratorStack_FullStackLevelMux(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	first := new(mockLogger)
	second := new(mockLogger)

	debugCalled := false
	warnCalled := false

	stack := wrapRecommendedStack(
		t,
		newMuxLogger(t, first, second),
		masker,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			debugCalled = true
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			warnCalled = true
		}),
	)
	logger := newLevelLogger(t, stack, int(decorator.LevelWarn))

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Debug(ctx, "debug secret")
	logger.Warn(ctx, "warn secret")
	logger.Error(ctx, errors.New("secret boom"))

	if debugCalled {
		t.Error("expected debug callback not to run")
	}

	if !warnCalled {
		t.Error("expected warn callback to run")
	}

	if len(masker.messages) != 1 || masker.messages[0] != "warn secret" {
		t.Errorf("expected masker to see only Warn, got %v", masker.messages)
	}

	if len(masker.errs) != 1 {
		t.Errorf("expected masker to see Error, got %v", masker.errs)
	}

	for idx, base := range []*mockLogger{first, second} {
		if len(base.debugIn) != 0 {
			t.Errorf("expected logger %d Debug to be dropped, got %+v", idx, base.debugIn)
		}

		if len(base.warnIn) != 1 || base.warnIn[0].msg != "warn ***" {
			t.Errorf("expected logger %d Warn 'warn ***', got %+v", idx, base.warnIn)
		}

		assertContextToken(t, base.warnIn[0].args, stackRedacted)

		if len(base.errorIn) != 1 ||
			base.errorIn[0].err == nil ||
			base.errorIn[0].err.Error() != "*** boom" {
			t.Errorf("expected logger %d Error '*** boom', got %+v", idx, base.errorIn)
		}

		assertContextToken(t, base.errorIn[0].args, stackRedacted)
	}
}
