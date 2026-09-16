package integration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
)

// TestAdapterStack_InfoMasksContextAndCallSite checks Level → Context → Mask →
// Callback → adapter: the backend emits flattened masked context and call-site
// fields.
func TestAdapterStack_InfoMasksContextAndCallSite(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)

		var gotMsg string

		var gotArgs []any

		logger := wrapRecommendedLevelDebug(
			t,
			sink.logger,
			masker,
			decorator.WithInfoCallback(func(_ context.Context, msg string, args ...any) {
				gotMsg = msg
				gotArgs = copyArgs(args)
			}),
		)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(ctx, "login secret", "password", stackSecret, "user", "alice")

		if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
			t.Errorf("expected masker message 'login secret', got %v", masker.messages)
		}

		assertContextToken(t, masker.args[0], stackSecret)

		if gotMsg != "login ***" {
			t.Errorf("expected callback message 'login ***', got %q", gotMsg)
		}

		assertContextToken(t, gotArgs, stackRedacted)

		entries := sink.snapshot()
		if len(entries) != 1 {
			t.Fatalf("expected 1 %s entry, got %d", sink.name, len(entries))
		}

		if entries[0].level != levelInfo {
			t.Errorf("expected %s level info, got %s", sink.name, entries[0].level)
		}

		assertMaskedLogin(t, sink, entries[0])
	})
}

// TestAdapterStack_ErrorMasksBeforeCallback checks that Error through the full
// adapter stack emits the masked error text and still unwraps for hooks.
func TestAdapterStack_ErrorMasksBeforeCallback(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		expectedErr := errors.New("secret leaked")

		var gotErr error

		logger := wrapRecommendedLevelDebug(
			t,
			sink.logger,
			masker,
			decorator.WithErrorCallback(func(_ context.Context, err error, _ ...any) {
				gotErr = err
			}),
		)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Error(ctx, expectedErr, "password", stackSecret)

		if gotErr == nil || gotErr.Error() != "*** leaked" {
			t.Errorf("expected callback error '*** leaked', got %v", gotErr)
		}

		if !errors.Is(gotErr, expectedErr) {
			t.Errorf("expected callback error to unwrap to %v", expectedErr)
		}

		entry := entryByLevel(t, sink, levelError)
		if entry.message != "*** leaked" {
			t.Errorf("expected %s error message '*** leaked', got %q", sink.name, entry.message)
		}

		assertField(t, sink, entry, "token", stackRedacted)
		assertField(t, sink, entry, "password", stackRedacted)
	})
}

// TestAdapterStack_NestedContextMapIsMasked checks that a nested context map
// is emitted as a nested object with the inner token redacted.
func TestAdapterStack_NestedContextMapIsMasked(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "auth", map[string]any{
			"token": stackSecret,
		})
		logger.Info(ctx, "ok")

		entry := entryByLevel(t, sink, levelInfo)
		if entry.message != "ok" {
			t.Errorf("expected %s message 'ok', got %q", sink.name, entry.message)
		}

		assertNestedField(t, sink, entry, "auth", "token", stackRedacted)
		assertNoField(t, sink, entry, "token")
	})
}

// TestAdapterStack_CallSiteOverridesContext checks that a call-site key after
// the context map wins; slog and zap may also keep the context value.
func TestAdapterStack_CallSiteOverridesContext(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "user", "from-ctx")
		logger.Info(ctx, "hello", "user", "from-call")

		entry := entryByLevel(t, sink, levelInfo)

		got, ok := lastField(entry, "user")
		if !ok {
			t.Fatalf("expected %s field user, got %v", sink.name, entry.fieldMap)
		}

		if got != "from-call" {
			t.Errorf("expected %s last user value 'from-call', got %#v", sink.name, got)
		}

		if sink.lastWins {
			return
		}

		values := fieldValues(entry, "user")
		if len(values) < 2 {
			t.Errorf(
				"expected %s to keep both user values, got %v",
				sink.name,
				values,
			)
		}
	})
}

// TestAdapterStack_EmptyContextMasksCallSiteOnly checks that an empty context
// does not add fields, while call-site secrets are still masked.
func TestAdapterStack_EmptyContextMasksCallSiteOnly(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		logger.Info(context.Background(), "login secret", "password", stackSecret)

		entry := entryByLevel(t, sink, levelInfo)
		if entry.message != "login ***" {
			t.Errorf("expected %s message 'login ***', got %q", sink.name, entry.message)
		}

		assertNoField(t, sink, entry, "token")
		assertField(t, sink, entry, "password", stackRedacted)
	})
}

// TestAdapterStack_AllLevelsAtLevelDebug checks that Debug, Info, Warn, and
// Error all emit through the full stack when the Level decorator is Debug.
func TestAdapterStack_AllLevelsAtLevelDebug(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Debug(ctx, "debug secret")
		logger.Info(ctx, "info secret")
		logger.Warn(ctx, "warn secret")
		logger.Error(ctx, errors.New("secret boom"))

		debugEntry := entryByLevel(t, sink, levelDebug)
		infoEntry := entryByLevel(t, sink, levelInfo)
		warnEntry := entryByLevel(t, sink, levelWarn)
		errorEntry := entryByLevel(t, sink, levelError)

		if debugEntry.message != "debug ***" {
			t.Errorf("expected debug 'debug ***', got %q", debugEntry.message)
		}

		if infoEntry.message != "info ***" {
			t.Errorf("expected info 'info ***', got %q", infoEntry.message)
		}

		if warnEntry.message != "warn ***" {
			t.Errorf("expected warn 'warn ***', got %q", warnEntry.message)
		}

		if errorEntry.message != "*** boom" {
			t.Errorf("expected error '*** boom', got %q", errorEntry.message)
		}

		for _, entry := range []capturedEntry{debugEntry, infoEntry, warnEntry, errorEntry} {
			assertField(t, sink, entry, "token", stackRedacted)
		}
	})
}

// TestAdapterStack_LevelInfoDropsDebug checks that LevelInfo outside the
// recommended adapter stack drops Debug before Mask and the backend.
func TestAdapterStack_LevelInfoDropsDebug(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelInfo(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Debug(ctx, "debug secret")
		logger.Info(ctx, "info secret")

		if len(masker.messages) != 1 || masker.messages[0] != "info secret" {
			t.Errorf("expected masker to see only Info, got %v", masker.messages)
		}

		assertNoLevel(t, sink, levelDebug)

		entry := entryByLevel(t, sink, levelInfo)
		if entry.message != "info ***" {
			t.Errorf("expected info 'info ***', got %q", entry.message)
		}

		assertField(t, sink, entry, "token", stackRedacted)
	})
}

// TestAdapterStack_ErrorNilAdapterMessages checks adapter-specific nil Error
// messages while context fields are still present.
func TestAdapterStack_ErrorNilAdapterMessages(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "request_id", "req-1")
		logger.Error(ctx, nil)

		entry := entryByLevel(t, sink, levelError)
		assertField(t, sink, entry, "request_id", "req-1")

		switch sink.name {
		case "logrus":
			if entry.message != "<nil>" {
				t.Errorf("expected logrus nil error message '<nil>', got %q", entry.message)
			}
		case "zerolog":
			if entry.hasMessage {
				t.Errorf("expected zerolog to omit message for nil error, got %q", entry.message)
			}
		default:
			if entry.message != "" {
				t.Errorf(
					"expected %s nil error message to be empty, got %q",
					sink.name,
					entry.message,
				)
			}
		}
	})
}

// TestAdapterStack_MaskOutsideContextLeaksContext checks that Mask outside
// Context leaves the context token unmasked in adapter fields.
func TestAdapterStack_MaskOutsideContextLeaksContext(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)

		callbackLogger, err := decorator.NewCallback(sink.logger)
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

		entry := entryByLevel(t, sink, levelInfo)
		if entry.message != "login ***" {
			t.Errorf("expected message 'login ***', got %q", entry.message)
		}

		assertField(t, sink, entry, "token", stackSecret)
		assertField(t, sink, entry, "password", stackRedacted)
	})
}

// TestAdapterStack_CallbackOutsideMaskSeesUnmasked checks that Callback outside
// Mask sees secrets while the adapter backend already has them redacted.
func TestAdapterStack_CallbackOutsideMaskSeesUnmasked(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)

		var gotMsg string

		var gotArgs []any

		maskLogger, err := decorator.NewMask(sink.logger, masker)
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

		entry := entryByLevel(t, sink, levelInfo)
		if entry.message != "login ***" {
			t.Errorf("expected backend message 'login ***', got %q", entry.message)
		}

		assertField(t, sink, entry, "token", stackRedacted)
		assertField(t, sink, entry, "password", stackRedacted)
	})
}

// TestAdapterStack_TrailingOrphanBecomesArgField checks that a trailing
// unpaired argument after the context map is emitted as the "arg" field.
func TestAdapterStack_TrailingOrphanBecomesArgField(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "token", "req")
		logger.Info(ctx, "hello", "orphan-value")

		entry := entryByLevel(t, sink, levelInfo)
		assertField(t, sink, entry, "token", "req")
		assertField(t, sink, entry, "arg", "orphan-value")
	})
}

// TestAdapterStack_MixedMapAndKeyValueFlattens checks that a call-site map is
// flattened together with context fields and key-value pairs.
func TestAdapterStack_MixedMapAndKeyValueFlattens(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(
			ctx,
			"login secret",
			map[string]any{"action": "login"},
			"password",
			stackSecret,
		)

		entry := entryByLevel(t, sink, levelInfo)
		if entry.message != "login ***" {
			t.Errorf("expected message 'login ***', got %q", entry.message)
		}

		assertField(t, sink, entry, "token", stackRedacted)
		assertField(t, sink, entry, "action", "login")
		assertField(t, sink, entry, "password", stackRedacted)
	})
}

// TestAdapterStack_LogrusEntryKeepsPreSetFields checks that wrapping a
// *logrus.Entry keeps fields already on that Entry through the recommended
// stack, together with masked context and call-site args.
func TestAdapterStack_LogrusEntryKeepsPreSetFields(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	sink := newLogrusEntrySink(t, logrus.DebugLevel)
	logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret, "user", "alice")

	if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
		t.Errorf("expected masker message 'login secret', got %v", masker.messages)
	}

	entry := entryByLevel(t, sink, levelInfo)
	assertMaskedLogin(t, sink, entry)
	assertField(t, sink, entry, "service", "api")
}

// TestAdapterStack_NopLoggerInnerStillRunsMaskAndCallback checks that NopLogger
// as the inner sink still lets Mask and Callback run.
func TestAdapterStack_NopLoggerInnerStillRunsMaskAndCallback(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	called := false

	logger := wrapRecommendedLevelDebug(
		t,
		cakelog.NopLogger(),
		masker,
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			called = true
		}),
	)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret)

	if !called {
		t.Error("expected info callback to run with NopLogger inner")
	}

	if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
		t.Errorf("expected masker to see the unmasked message, got %v", masker.messages)
	}
}
