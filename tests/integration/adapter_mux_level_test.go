package integration_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
	"go.uber.org/zap/zapcore"
)

func assertMaskedInfoEntry(t *testing.T, sink *adapterSink) {
	t.Helper()

	entry := entryByLevel(t, sink, levelInfo)
	if entry.message != "login ***" {
		t.Errorf("expected %s message 'login ***', got %q", sink.name, entry.message)
	}

	assertField(t, sink, entry, "token", stackRedacted)
	assertField(t, sink, entry, "password", stackRedacted)
}

// TestAdapterStack_MuxTwoAdapters checks that Mux under the recommended stack
// delivers the same masked payload to two different adapters.
func TestAdapterStack_MuxTwoAdapters(t *testing.T) {
	t.Parallel()

	t.Run("slog_zap", func(t *testing.T) {
		t.Parallel()

		masker := new(stackMasker)
		slogSink := newSlogSink(t, slog.LevelDebug)
		zapSink := newZapSink(t, zapcore.DebugLevel)

		logger := wrapRecommendedStack(
			t,
			newMuxLogger(t, slogSink.logger, zapSink.logger),
			masker,
		)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(ctx, "login secret", "password", stackSecret)

		if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
			t.Errorf("expected masker to run once, got %v", masker.messages)
		}

		assertMaskedInfoEntry(t, slogSink)
		assertMaskedInfoEntry(t, zapSink)
	})

	t.Run("logrus_zerolog", func(t *testing.T) {
		t.Parallel()

		masker := new(stackMasker)
		logrusSink := newLogrusSink(t, logrus.DebugLevel)
		zerologSink := newZerologSink(t, zerolog.DebugLevel)

		logger := wrapRecommendedStack(
			t,
			newMuxLogger(t, logrusSink.logger, zerologSink.logger),
			masker,
		)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(ctx, "login secret", "password", stackSecret)

		if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
			t.Errorf("expected masker to run once, got %v", masker.messages)
		}

		assertMaskedInfoEntry(t, logrusSink)
		assertMaskedInfoEntry(t, zerologSink)
	})

	t.Run("slog_logrus_entry", func(t *testing.T) {
		t.Parallel()

		masker := new(stackMasker)
		slogSink := newSlogSink(t, slog.LevelDebug)
		logrusSink := newLogrusEntrySink(t, logrus.DebugLevel)

		logger := wrapRecommendedStack(
			t,
			newMuxLogger(t, slogSink.logger, logrusSink.logger),
			masker,
		)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(ctx, "login secret", "password", stackSecret)

		if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
			t.Errorf("expected masker to run once, got %v", masker.messages)
		}

		assertMaskedInfoEntry(t, slogSink)
		assertMaskedInfoEntry(t, logrusSink)
		assertField(t, logrusSink, entryByLevel(t, logrusSink, levelInfo), "service", "api")
		assertNoField(t, slogSink, entryByLevel(t, slogSink, levelInfo), "service")
	})
}

// TestAdapterStack_MuxAllAdapters checks that one Info call fans out to all
// four adapters with the same masked fields.
func TestAdapterStack_MuxAllAdapters(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	slogSink := newSlogSink(t, slog.LevelDebug)
	zapSink := newZapSink(t, zapcore.DebugLevel)
	logrusSink := newLogrusSink(t, logrus.DebugLevel)
	zerologSink := newZerologSink(t, zerolog.DebugLevel)

	logger := wrapRecommendedStack(
		t,
		newMuxLogger(t, slogSink.logger, zapSink.logger, logrusSink.logger, zerologSink.logger),
		masker,
	)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "login secret", "password", stackSecret)

	if len(masker.messages) != 1 {
		t.Errorf("expected masker to run once, got %v", masker.messages)
	}

	for _, sink := range []*adapterSink{slogSink, zapSink, logrusSink, zerologSink} {
		assertMaskedInfoEntry(t, sink)
	}
}

// TestAdapterStack_MuxErrorMasksAllAdapters checks Error fan-out to every
// adapter under the recommended stack.
func TestAdapterStack_MuxErrorMasksAllAdapters(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	slogSink := newSlogSink(t, slog.LevelDebug)
	zapSink := newZapSink(t, zapcore.DebugLevel)
	logrusSink := newLogrusSink(t, logrus.DebugLevel)
	zerologSink := newZerologSink(t, zerolog.DebugLevel)
	expectedErr := errors.New("secret leaked")

	logger := wrapRecommendedStack(
		t,
		newMuxLogger(t, slogSink.logger, zapSink.logger, logrusSink.logger, zerologSink.logger),
		masker,
	)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Error(ctx, expectedErr, "password", stackSecret)

	if len(masker.errs) != 1 || !errors.Is(masker.errs[0], expectedErr) {
		t.Errorf("expected masker to see 1 unmasked error, got %v", masker.errs)
	}

	for _, sink := range []*adapterSink{slogSink, zapSink, logrusSink, zerologSink} {
		entry := entryByLevel(t, sink, levelError)
		if entry.message != "*** leaked" {
			t.Errorf("expected %s error '*** leaked', got %q", sink.name, entry.message)
		}

		assertField(t, sink, entry, "token", stackRedacted)
		assertField(t, sink, entry, "password", stackRedacted)
	}
}

// TestAdapterStack_MuxPerBranchLevel checks different Level thresholds on
// adapter branches under a shared Context→Mask stack.
func TestAdapterStack_MuxPerBranchLevel(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	slogSink := newSlogSink(t, slog.LevelDebug)
	zapSink := newZapSink(t, zapcore.DebugLevel)

	mux := newMuxLogger(
		t,
		newLevelLogger(t, slogSink.logger, int(decorator.LevelWarn)),
		newLevelLogger(t, zapSink.logger, int(decorator.LevelError)),
	)
	logger := wrapContextMask(t, mux, masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Warn(ctx, "warn secret")
	logger.Error(ctx, errors.New("secret boom"))

	assertNoLevel(t, zapSink, levelWarn)

	warnEntry := entryByLevel(t, slogSink, levelWarn)
	if warnEntry.message != "warn ***" {
		t.Errorf("expected slog warn 'warn ***', got %q", warnEntry.message)
	}

	assertField(t, slogSink, warnEntry, "token", stackRedacted)

	for _, sink := range []*adapterSink{slogSink, zapSink} {
		entry := entryByLevel(t, sink, levelError)
		if entry.message != "*** boom" {
			t.Errorf("expected %s error '*** boom', got %q", sink.name, entry.message)
		}

		assertField(t, sink, entry, "token", stackRedacted)
	}
}

// TestAdapterStack_BackendWarnDropsInfoAfterCallback checks that an Info
// callback still runs when the adapter backend is filtered to Warn.
func TestAdapterStack_BackendWarnDropsInfoAfterCallback(t *testing.T) {
	t.Parallel()

	forEachWarnAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		called := false

		logger := wrapRecommendedLevelDebug(
			t,
			sink.logger,
			masker,
			decorator.WithInfoCallback(func(context.Context, string, ...any) {
				called = true
			}),
		)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(ctx, "login secret")

		if !called {
			t.Error("expected info callback to run even when the backend drops Info")
		}

		if len(masker.messages) != 1 || masker.messages[0] != "login secret" {
			t.Errorf("expected masker to see Info, got %v", masker.messages)
		}

		if len(sink.snapshot()) != 0 {
			t.Errorf("expected %s backend to drop Info, got %+v", sink.name, sink.snapshot())
		}
	})
}

// TestAdapterStack_ZapStillEmitsContextFields checks that zap ignores the
// context object but still logs fields extracted by the Context decorator.
func TestAdapterStack_ZapStillEmitsContextFields(t *testing.T) {
	t.Parallel()

	masker := new(stackMasker)
	sink := newZapSink(t, zapcore.DebugLevel)
	logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

	ctx := withTestValue(context.Background(), "token", stackSecret)
	logger.Info(ctx, "hello")

	if sink.tracksCtx {
		t.Fatal("expected zap sink not to track context identity")
	}

	entry := entryByLevel(t, sink, levelInfo)
	assertField(t, sink, entry, "token", stackRedacted)

	if entry.ctx != nil {
		t.Errorf("expected zap captured context to be nil, got %v", entry.ctx)
	}
}

// TestAdapterStack_ContextReachesBackend checks that slog, logrus, and zerolog
// receive the same context value that was passed to the stacked logger.
func TestAdapterStack_ContextReachesBackend(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		if !sink.tracksCtx {
			return
		}

		masker := new(stackMasker)
		logger := wrapRecommendedLevelDebug(t, sink.logger, masker)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(ctx, "hello")

		entry := entryByLevel(t, sink, levelInfo)
		if entry.ctx != ctx {
			t.Errorf("expected %s backend context to match the log call", sink.name)
		}
	})
}

// TestAdapterStack_MuxAdapterWithNopLogger checks that Mux(adapter, NopLogger)
// still emits on the adapter branch.
func TestAdapterStack_MuxAdapterWithNopLogger(t *testing.T) {
	t.Parallel()

	forEachAdapter(t, func(t *testing.T, sink *adapterSink) {
		t.Helper()

		masker := new(stackMasker)
		logger := wrapRecommendedStack(
			t,
			newMuxLogger(t, sink.logger, cakelog.NopLogger()),
			masker,
		)

		ctx := withTestValue(context.Background(), "token", stackSecret)
		logger.Info(ctx, "login secret", "password", stackSecret)

		assertMaskedInfoEntry(t, sink)
	})
}
