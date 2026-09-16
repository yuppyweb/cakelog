package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/logrusadapter"
	"github.com/yuppyweb/cakelog/adapter/slogadapter"
	"github.com/yuppyweb/cakelog/adapter/zapadapter"
	"github.com/yuppyweb/cakelog/adapter/zerologadapter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
)

// capturedField is a single structured field emitted by an adapter backend.
type capturedField struct {
	key   string
	value any
}

// capturedEntry is one log record captured from an adapter backend.
type capturedEntry struct {
	level      string
	message    string
	hasMessage bool
	fields     []capturedField
	fieldMap   map[string]any
	ctx        context.Context
}

// adapterSink is a cakelog.Logger backed by a real adapter with captured output.
type adapterSink struct {
	name      string
	logger    cakelog.Logger
	entries   *[]capturedEntry
	lastWins  bool
	tracksCtx bool
	refresh   func()
}

func (s *adapterSink) snapshot() []capturedEntry {
	if s.refresh != nil {
		s.refresh()
	}

	if s.entries == nil {
		return nil
	}

	return *s.entries
}

func forEachAdapter(t *testing.T, fn func(t *testing.T, sink *adapterSink)) {
	t.Helper()

	cases := []struct {
		name string
		new  func(*testing.T) *adapterSink
	}{
		{
			name: "slog",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newSlogSink(t, slog.LevelDebug)
			},
		},
		{
			name: "zap",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newZapSink(t, zapcore.DebugLevel)
			},
		},
		{
			name: "logrus",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newLogrusSink(t, logrus.DebugLevel)
			},
		},
		{
			name: "zerolog",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newZerologSink(t, zerolog.DebugLevel)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fn(t, tc.new(t))
		})
	}
}

func forEachWarnAdapter(t *testing.T, fn func(t *testing.T, sink *adapterSink)) {
	t.Helper()

	cases := []struct {
		name string
		new  func(*testing.T) *adapterSink
	}{
		{
			name: "slog",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newSlogSink(t, slog.LevelWarn)
			},
		},
		{
			name: "zap",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newZapSink(t, zapcore.WarnLevel)
			},
		},
		{
			name: "logrus",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newLogrusSink(t, logrus.WarnLevel)
			},
		},
		{
			name: "zerolog",
			new: func(t *testing.T) *adapterSink {
				t.Helper()

				return newZerologSink(t, zerolog.WarnLevel)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fn(t, tc.new(t))
		})
	}
}

type slogCaptureHandler struct {
	minLevel slog.Level
	entries  *[]capturedEntry
}

func (h *slogCaptureHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.minLevel
}

func (h *slogCaptureHandler) Handle(ctx context.Context, record slog.Record) error {
	cloned := record.Clone()
	fields := make([]capturedField, 0, cloned.NumAttrs())
	fieldMap := make(map[string]any, cloned.NumAttrs())

	cloned.Attrs(func(attr slog.Attr) bool {
		value := attr.Value.Any()
		fields = append(fields, capturedField{key: attr.Key, value: value})
		fieldMap[attr.Key] = value

		return true
	})

	*h.entries = append(*h.entries, capturedEntry{
		level:      strings.ToLower(cloned.Level.String()),
		message:    cloned.Message,
		hasMessage: true,
		fields:     fields,
		fieldMap:   fieldMap,
		ctx:        ctx,
	})

	return nil
}

func (h *slogCaptureHandler) WithAttrs([]slog.Attr) slog.Handler {
	return slog.DiscardHandler
}

func (h *slogCaptureHandler) WithGroup(string) slog.Handler {
	return slog.DiscardHandler
}

var _ slog.Handler = (*slogCaptureHandler)(nil)

func newSlogSink(t *testing.T, minLevel slog.Level) *adapterSink {
	t.Helper()

	entries := make([]capturedEntry, 0)
	handler := &slogCaptureHandler{minLevel: minLevel, entries: &entries}

	logger, err := slogadapter.New(slog.New(handler))
	if err != nil {
		t.Fatalf("failed to create slog adapter: %v", err)
	}

	return &adapterSink{
		name:      "slog",
		logger:    logger,
		entries:   &entries,
		lastWins:  false,
		tracksCtx: true,
	}
}

type zapCaptureCore struct {
	minLevel zapcore.Level
	entries  *[]capturedEntry
}

func (c *zapCaptureCore) Enabled(level zapcore.Level) bool {
	return level >= c.minLevel
}

func (c *zapCaptureCore) With([]zapcore.Field) zapcore.Core {
	return c
}

func (c *zapCaptureCore) Check(
	ent zapcore.Entry,
	ce *zapcore.CheckedEntry,
) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func (c *zapCaptureCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	captured := make([]capturedField, 0, len(fields))
	fieldMap := make(map[string]any, len(fields))

	for _, field := range fields {
		value := zapFieldValue(field)
		captured = append(captured, capturedField{key: field.Key, value: value})
		fieldMap[field.Key] = value
	}

	*c.entries = append(*c.entries, capturedEntry{
		level:      ent.Level.String(),
		message:    ent.Message,
		hasMessage: true,
		fields:     captured,
		fieldMap:   fieldMap,
	})

	return nil
}

func (c *zapCaptureCore) Sync() error {
	return nil
}

var _ zapcore.Core = (*zapCaptureCore)(nil)

func zapFieldValue(field zapcore.Field) any {
	encoder := zapcore.NewMapObjectEncoder()
	field.AddTo(encoder)

	return encoder.Fields[field.Key]
}

func newZapSink(t *testing.T, minLevel zapcore.Level) *adapterSink {
	t.Helper()

	entries := make([]capturedEntry, 0)
	core := &zapCaptureCore{minLevel: minLevel, entries: &entries}

	logger, err := zapadapter.New(zap.New(core))
	if err != nil {
		t.Fatalf("failed to create zap adapter: %v", err)
	}

	return &adapterSink{
		name:      "zap",
		logger:    logger,
		entries:   &entries,
		lastWins:  false,
		tracksCtx: false,
	}
}

type logrusCaptureHook struct {
	entries *[]capturedEntry
}

func (h *logrusCaptureHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logrusCaptureHook) Fire(entry *logrus.Entry) error {
	fields := make([]capturedField, 0, len(entry.Data))
	fieldMap := make(map[string]any, len(entry.Data))

	for key, value := range entry.Data {
		fields = append(fields, capturedField{key: key, value: value})
		fieldMap[key] = value
	}

	*h.entries = append(*h.entries, capturedEntry{
		level:      normalizeLevel(entry.Level.String()),
		message:    entry.Message,
		hasMessage: true,
		fields:     fields,
		fieldMap:   fieldMap,
		ctx:        entry.Context,
	})

	return nil
}

var _ logrus.Hook = (*logrusCaptureHook)(nil)

func newLogrusSink(t *testing.T, minLevel logrus.Level) *adapterSink {
	t.Helper()

	entries := make([]capturedEntry, 0)

	log := logrus.New()
	log.SetOutput(io.Discard)
	log.Level = minLevel
	log.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})
	log.AddHook(&logrusCaptureHook{entries: &entries})

	logger, err := logrusadapter.New(log)
	if err != nil {
		t.Fatalf("failed to create logrus adapter: %v", err)
	}

	return &adapterSink{
		name:      "logrus",
		logger:    logger,
		entries:   &entries,
		lastWins:  true,
		tracksCtx: true,
	}
}

type zerologCaptureHook struct {
	ctxs *[]context.Context
}

func (h *zerologCaptureHook) Run(event *zerolog.Event, _ zerolog.Level, _ string) {
	*h.ctxs = append(*h.ctxs, event.GetCtx())
}

func newZerologSink(t *testing.T, minLevel zerolog.Level) *adapterSink {
	t.Helper()

	buf := &bytes.Buffer{}
	ctxs := make([]context.Context, 0)
	log := zerolog.New(buf).Level(minLevel).Hook(&zerologCaptureHook{ctxs: &ctxs})

	logger, err := zerologadapter.New(&log)
	if err != nil {
		t.Fatalf("failed to create zerolog adapter: %v", err)
	}

	entries := make([]capturedEntry, 0)
	sink := &adapterSink{
		name:      "zerolog",
		logger:    logger,
		entries:   &entries,
		lastWins:  true,
		tracksCtx: true,
	}

	sink.refresh = func() {
		parsed := parseZerologBuffer(t, buf, ctxs)
		*sink.entries = parsed
	}

	return sink
}

func parseZerologBuffer(
	t *testing.T,
	buf *bytes.Buffer,
	ctxs []context.Context,
) []capturedEntry {
	t.Helper()

	output := bytes.TrimSpace(buf.Bytes())
	if len(output) == 0 {
		return nil
	}

	lines := bytes.Split(output, []byte("\n"))
	entries := make([]capturedEntry, 0, len(lines))

	for idx, line := range lines {
		raw := make(map[string]any)
		if unmarshalErr := json.Unmarshal(line, &raw); unmarshalErr != nil {
			t.Fatalf("failed to unmarshal zerolog output %q: %v", string(line), unmarshalErr)
		}

		level, _ := raw["level"].(string)
		message, hasMessage := raw["message"].(string)

		delete(raw, "level")
		delete(raw, "message")

		fields := make([]capturedField, 0, len(raw))
		fieldMap := make(map[string]any, len(raw))

		for key, value := range raw {
			fields = append(fields, capturedField{key: key, value: value})
			fieldMap[key] = value
		}

		entry := capturedEntry{
			level:      normalizeLevel(level),
			message:    message,
			hasMessage: hasMessage,
			fields:     fields,
			fieldMap:   fieldMap,
		}

		if idx < len(ctxs) {
			entry.ctx = ctxs[idx]
		}

		entries = append(entries, entry)
	}

	return entries
}

func normalizeLevel(level string) string {
	if level == "warning" {
		return levelWarn
	}

	return strings.ToLower(level)
}

func fieldValues(entry capturedEntry, key string) []any {
	if len(entry.fields) == 0 {
		value, ok := entry.fieldMap[key]
		if !ok {
			return nil
		}

		return []any{value}
	}

	values := make([]any, 0)

	for _, field := range entry.fields {
		if field.key == key {
			values = append(values, field.value)
		}
	}

	return values
}

func lastField(entry capturedEntry, key string) (any, bool) {
	if value, ok := entry.fieldMap[key]; ok {
		return value, true
	}

	values := fieldValues(entry, key)
	if len(values) == 0 {
		return nil, false
	}

	return values[len(values)-1], true
}

func entryByLevel(t *testing.T, sink *adapterSink, level string) capturedEntry {
	t.Helper()

	for _, entry := range sink.snapshot() {
		if entry.level == level {
			return entry
		}
	}

	t.Fatalf("expected %s sink to emit level %s, got %+v", sink.name, level, sink.snapshot())

	return capturedEntry{}
}

func assertNoLevel(t *testing.T, sink *adapterSink, level string) {
	t.Helper()

	for _, entry := range sink.snapshot() {
		if entry.level == level {
			t.Errorf("expected %s sink not to emit level %s, got %+v", sink.name, level, entry)
		}
	}
}

func assertField(t *testing.T, sink *adapterSink, entry capturedEntry, key string, want any) {
	t.Helper()

	got, ok := lastField(entry, key)
	if !ok {
		t.Errorf("expected %s field %q, got %v", sink.name, key, entry.fieldMap)

		return
	}

	if !fieldEqual(got, want) {
		t.Errorf("expected %s field %q to be %#v, got %#v", sink.name, key, want, got)
	}
}

func assertNoField(t *testing.T, sink *adapterSink, entry capturedEntry, key string) {
	t.Helper()

	if _, ok := lastField(entry, key); ok {
		t.Errorf("expected %s field %q to be absent, got %v", sink.name, key, entry.fieldMap)
	}
}

func assertNestedField(
	t *testing.T,
	sink *adapterSink,
	entry capturedEntry,
	key string,
	inner string,
	want any,
) {
	t.Helper()

	got, ok := lastField(entry, key)
	if !ok {
		t.Errorf("expected %s field %q, got %v", sink.name, key, entry.fieldMap)

		return
	}

	nested, ok := got.(map[string]any)
	if !ok {
		t.Errorf("expected %s field %q to be map[string]any, got %T", sink.name, key, got)

		return
	}

	innerGot, ok := nested[inner]
	if !ok {
		t.Errorf("expected %s nested field %s.%s, got %v", sink.name, key, inner, nested)

		return
	}

	if !fieldEqual(innerGot, want) {
		t.Errorf(
			"expected %s nested field %s.%s to be %#v, got %#v",
			sink.name,
			key,
			inner,
			want,
			innerGot,
		)
	}
}

func fieldEqual(got, want any) bool {
	gotMap, gotIsMap := got.(map[string]any)
	wantMap, wantIsMap := want.(map[string]any)

	if gotIsMap && wantIsMap {
		if len(gotMap) != len(wantMap) {
			return false
		}

		for key, wantValue := range wantMap {
			gotValue, ok := gotMap[key]
			if !ok || !fieldEqual(gotValue, wantValue) {
				return false
			}
		}

		return true
	}

	return got == want
}

func assertMaskedLogin(t *testing.T, sink *adapterSink, entry capturedEntry) {
	t.Helper()

	if entry.message != "login ***" {
		t.Errorf("expected %s message 'login ***', got %q", sink.name, entry.message)
	}

	assertField(t, sink, entry, "token", stackRedacted)
	assertField(t, sink, entry, "password", stackRedacted)
	assertField(t, sink, entry, "user", "alice")
}
