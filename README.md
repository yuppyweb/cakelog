# 🍰 Cakelog

[![Go Version](https://img.shields.io/github/go-mod/go-version/yuppyweb/cakelog)](https://github.com/yuppyweb/cakelog) 
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) 

**Cakelog** is a Go logging library with a small unified `Logger` interface. Swap the backend through adapters (Slog, Zap, Logrus, Zerolog) and add behavior with composable decorators.

## ✨ Features

✅ **Unified Logger Interface** — Write application code once, switch backends anytime  
✅ **Multiple Adapters** — Slog, Zap, Logrus, and Zerolog, each in its own package  
✅ **Structured Fields** — Mix slog-style key-value pairs and maps; adapters parse them the same way  
✅ **Context Enrichment** — Prepend fields extracted from `context.Context`  
✅ **Callback Hooks** — Run custom logic after each log call  
✅ **Data Masking** — Sanitize messages, errors, and arguments before they are emitted  
✅ **Level Filter** — Drop calls below a minimum severity  
✅ **Mux** — Fan the same call out to several loggers  
✅ **Composable Decorators** — Stack wrappers in a predictable order  
✅ **NopLogger for Testing** — Built-in no-op implementation

## 🎯 Core Idea

Application code depends on `cakelog.Logger`, not on a specific logging library. An adapter turns a backend logger into that interface. Decorators wrap any `Logger` and can be stacked.

```
┌─ Your Application
│
├─ cakelog.Logger ◄──── Unified Interface
│
├─ Decorators (outermost first):
│  • Level   — drop calls below a threshold
│  • Context — prepend fields from ctx
│  • Mask    — sanitize message, error, args
│  • Callback — hooks after the call is forwarded
│  • Mux     — fan out to several loggers
│
└─ Adapters (choose one or mux several):
   • Slog
   • Zap
   • Logrus
   • Zerolog
```

## 📦 Installation

```bash
go get github.com/yuppyweb/cakelog
```

Import only the adapter you need. Construct the backend logger yourself:

```bash
go get github.com/sirupsen/logrus     # Logrus adapter
go get github.com/rs/zerolog          # Zerolog adapter
go get go.uber.org/zap                # Zap adapter
# Slog is in the Go standard library
```

## 🚀 Quick Start

### Basic Usage with Slog 📚

```go
package main

import (
    "context"
    "errors"
    "log/slog"
    "os"

    "github.com/yuppyweb/cakelog/adapter/slogadapter"
)

func main() {
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger, err := slogadapter.New(slogLogger)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    logger.Info(ctx, "application started", "env", "prod")
    logger.Debug(ctx, "debug message")
    logger.Warn(ctx, "warning message")
    logger.Error(ctx, errors.New("connection failed"))
}
```

## 📖 Core Components

### Logger Interface 📋

```go
type Logger interface {
    Debug(ctx context.Context, msg string, args ...any)
    Info(ctx context.Context, msg string, args ...any)
    Warn(ctx context.Context, msg string, args ...any)
    Error(ctx context.Context, err error, args ...any)
}
```

`Error` has no separate message argument: the message is `err.Error()`. A nil `err` is allowed; each backend chooses the message shape (slog and zap use an empty string, zerolog omits the message field, logrus prints `<nil>`).

Built-in adapters also pass a non-nil `err` through the backend's error API under the key `error`: slog attribute `error`, `zap.Error`, zerolog `Event.Err`, logrus `WithError`. A call-site field with the same key follows that backend's duplicate-key rules. Args may mix alternating key-value pairs and maps.

### NopLogger 🚫

```go
logger := cakelog.NopLogger()
logger.Info(ctx, "this is discarded")
```

### Structured Fields 🧩

`adapter.Fields` is the shared parser used by every built-in adapter. Arguments are read left to right:

- a map is expanded into fields (non-string keys are formatted with `fmt.Sprint`)
- otherwise a string key is paired with the next value
- a trailing value without a key is stored under `"arg"`
- a map used as a pair value is kept as a single field
- duplicate keys are kept in encounter order

`adapter.Fields` always preserves duplicates. What is then emitted is up to the backend, not Cakelog:

- slog and zap keep every occurrence
- logrus and zerolog keep the last value (map / JSON object)

This is intentional: swapping adapters keeps Cakelog's parse rules, but not a backend's own field model.

```go
logger.Info(ctx, "request finished",
    "method", "GET",
    map[string]any{"status": 200, "path": "/health"},
)
```

## 🔌 Adapters

Each adapter lives in its own package and exposes `New`. A nil backend returns a package-specific error.

Cakelog does not rewrite caller or source location. If the backend records file:line (slog `AddSource`, zap `AddCaller`, logrus `ReportCaller`, zerolog `Caller`), it is the adapter method — and decorator frames when the logger is wrapped. A fixed skip such as zap `AddCallerSkip` only matches one wrap depth. Configure skip or a caller hook on the backend if you need the application call site.

### Slog Adapter 🔮

```go
import (
    "log/slog"

    "github.com/yuppyweb/cakelog/adapter/slogadapter"
)

slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger, err := slogadapter.New(slogLogger)
```

### Zap Adapter ⚡

Zap's standard methods do not take `context.Context`. The adapter still accepts `ctx` on the `Logger` interface, but does not forward it.

`cakelog.Logger` has no `Sync`. A production zap logger buffers output; call `Sync` on the `*zap.Logger` you passed to `New` before shutdown, or the last entries may be lost. `Sync` on stdout or stderr can return a non-nil error on some platforms; ignoring it is common.

```go
import (
    "go.uber.org/zap"

    "github.com/yuppyweb/cakelog/adapter/zapadapter"
)

zapLogger, err := zap.NewProduction()
if err != nil {
    log.Fatal(err)
}
defer func() { _ = zapLogger.Sync() }()

logger, err := zapadapter.New(zapLogger)
```

### Logrus Adapter 📊

`New` takes a `logrus.FieldLogger`: `*logrus.Logger`, `*logrus.Entry`, or another implementation. A nil value, including a typed nil pointer stored in the interface, returns `ErrNilLogrusLogger`.

```go
import (
    "github.com/sirupsen/logrus"
    "github.com/yuppyweb/cakelog/adapter/logrusadapter"
)

logrusLogger := logrus.New()
logger, err := logrusadapter.New(logrusLogger)

entry := logrusLogger.WithField("service", "api")
logger, err = logrusadapter.New(entry)
```

### Zerolog Adapter 📬

`New` takes a pointer to `zerolog.Logger`.

```go
import (
    "os"

    "github.com/rs/zerolog"
    "github.com/yuppyweb/cakelog/adapter/zerologadapter"
)

zlog := zerolog.New(os.Stdout)
logger, err := zerologadapter.New(&zlog)
```

## 🎨 Decorators

Every decorator constructor returns `(cakelog.Logger, error)`. A nil logger, including a typed nil stored in the `Logger` interface, returns a wrapped `decorator.ErrNilLogger`.

### Context Decorator 🎁

`NewContext` prepends `fields(ctx)` as a `map[string]any` before call-site args. A nil or empty map leaves args unchanged. Call-site pairs after the map override extracted fields when an adapter keeps the last value.

You supply the extractor; Cakelog does not store values on the context for you.

Typical uses:

- Attach request, correlation, user, or tenant IDs once on the context instead of repeating them at every call site
- Carry tracing identifiers through a request without coupling handlers to the logger
- Inject static service metadata (name, version, environment) by returning a map that ignores `ctx`
- Keep application logs focused on the event; shared fields come from the request scope

To mask context fields, wrap `NewMask` inside `NewContext` so the masker sees the prepended map. A masker must walk that map (and other nested maps); masking only top-level strings leaves context values unchanged.

```go
import (
    "context"

    "github.com/yuppyweb/cakelog/decorator"
)

type ctxKey struct{}

func withFields(ctx context.Context, fields map[string]any) context.Context {
    return context.WithValue(ctx, ctxKey{}, fields)
}

func fieldsFromCtx(ctx context.Context) map[string]any {
    fields, _ := ctx.Value(ctxKey{}).(map[string]any)
    return fields
}

logger, err := decorator.NewContext(baseLogger, fieldsFromCtx)
if err != nil {
    log.Fatal(err)
}

ctx := withFields(context.Background(), map[string]any{
    "user_id":    "user123",
    "request_id": "req456",
})

logger.Info(ctx, "request started")
```

### Callback Decorator ↩️

Callbacks run after the call is forwarded, on the `Logger` method itself, even if the adapter later discards the entry. If the underlying logger panics, the callback does not run. Put `NewMask` outside `NewCallback` if hooks must not see secrets. An `Error` callback may receive a nil error.

Typical uses:

- Increment metrics or counters per log level without changing call sites
- Page or notify on errors (chat, incident tools) from a single hook
- Record audit side effects when a warning or error is logged
- Observe log calls in tests without parsing backend output
- React to a log method being invoked even when the adapter later drops the entry

```go
callbackLogger, err := decorator.NewCallback(
    baseLogger,
    decorator.WithDebugCallback(func(ctx context.Context, msg string, args ...any) {
        // track debug events
    }),
    decorator.WithErrorCallback(func(ctx context.Context, err error, args ...any) {
        // alert on errors
    }),
)
```

Omitted levels use no-op callbacks. If the same level is configured more than once, the last option wins.

### Mask Decorator 🎭

Maskers run in order. Do not mutate received values, including nested maps and slices; return new values instead.

Typical uses:

- Redact passwords, tokens, cookies, and other secrets before logs leave the process
- Strip or hash PII (emails, phone numbers, names) for GDPR-style constraints
- Replace payment card numbers and similar PCI data in messages, errors, and fields
- Chain several maskers when different patterns belong to different policies
- Sanitize nested maps, including fields prepended by `NewContext`

```go
import (
    "strings"

    "github.com/yuppyweb/cakelog/decorator"
)

type maskedError struct {
    err error
    msg string
}

func (e maskedError) Error() string { return e.msg }
func (e maskedError) Unwrap() error { return e.err }

type MyMasker struct{}

func (m MyMasker) MaskMessage(msg string) string {
    return strings.ReplaceAll(msg, "password", "***")
}

func (m MyMasker) MaskError(err error) error {
    if err == nil {
        return nil
    }

    return maskedError{
        err: err,
        msg: strings.ReplaceAll(err.Error(), "secret", "***"),
    }
}

func (m MyMasker) MaskArguments(args ...any) []any {
    masked := make([]any, len(args))
    for i, arg := range args {
        masked[i] = maskArg(arg)
    }
    return masked
}

func maskArg(arg any) any {
    switch val := arg.(type) {
    case string:
        return strings.ReplaceAll(val, "token", "xxx")
    case map[string]any:
        out := make(map[string]any, len(val))
        for key, value := range val {
            out[key] = maskArg(value)
        }
        return out
    default:
        return arg
    }
}

maskLogger, err := decorator.NewMask(baseLogger, MyMasker{})
```

### Level Decorator 📶

`NewLevel` drops calls below an inclusive minimum. Valid levels are `LevelDebug`, `LevelInfo`, `LevelWarn`, and `LevelError`. The underlying logger may still drop the entry according to its own configuration.

Typical uses:

- Raise verbosity in development and keep only warnings and errors in production without swapping adapters
- Filter independently of a backend's own level, including when several adapters disagree
- Place it outermost so cheaper drops skip context extraction, masking, and callbacks
- Give a subsystem a quieter logger while the rest of the process stays verbose

```go
logger, err := decorator.NewLevel(baseLogger, decorator.LevelWarn)
// Debug and Info are dropped; Warn and Error are forwarded.
```

### Mux Decorator 🔀

`NewMux` forwards each call to every logger in order. The args slice is copied per logger, not the values it contains. A panic in one logger is not recovered and stops the remaining loggers.

Typical uses:

- Write the same events to several backends at once (console and a log aggregator)
- Dual-write during a migration from one logging library to another
- Keep a human-readable sink next to structured JSON for the same process
- Fan out to a dedicated audit or file logger alongside the primary adapter

Mux does not route by level: every wrapped logger receives every forwarded call. Combine it with `NewLevel` on a branch if one sink should be quieter.

```go
logger, err := decorator.NewMux(stdoutLogger, auditLogger)
```

### Composing Decorators 🧩

Recommended production stack, outermost first: **Context → Mask → Callback → adapter**. Context is outermost so maskers see the prepended map. Mask wraps Callback so hooks receive already-masked data.

```go
baseLogger, _ := slogadapter.New(slogLogger)
callbackLogger, _ := decorator.NewCallback(baseLogger, opts...)
maskLogger, _ := decorator.NewMask(callbackLogger, masker)
logger, _ := decorator.NewContext(maskLogger, fieldsFromCtx)
```

Add `NewLevel` outside that stack to filter before enrichment and masking. Use `NewMux` at the adapter layer to write to several backends.

## 📝 Examples

### Service with Context Enrichment 🏗️

```go
package main

import (
    "context"
    "log/slog"
    "os"

    "github.com/yuppyweb/cakelog"
    "github.com/yuppyweb/cakelog/adapter/slogadapter"
    "github.com/yuppyweb/cakelog/decorator"
)

type ctxKey struct{}

func withFields(ctx context.Context, fields map[string]any) context.Context {
    return context.WithValue(ctx, ctxKey{}, fields)
}

func fieldsFromCtx(ctx context.Context) map[string]any {
    fields, _ := ctx.Value(ctxKey{}).(map[string]any)
    return fields
}

func processRequest(logger cakelog.Logger, requestID, userID string) {
    ctx := withFields(context.Background(), map[string]any{
        "request_id": requestID,
        "user_id":    userID,
    })

    logger.Info(ctx, "request processing started")
    logger.Info(ctx, "request processing completed")
}

func main() {
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    baseLogger, _ := slogadapter.New(slogLogger)
    logger, _ := decorator.NewContext(baseLogger, fieldsFromCtx)

    processRequest(logger, "req-123", "user-456")
}
```

### Combining Adapters and Decorators 🧩

```go
func setupProductionLogger() (cakelog.Logger, error) {
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    baseLogger, err := slogadapter.New(slogLogger)
    if err != nil {
        return nil, err
    }

    callbackLogger, err := decorator.NewCallback(
        baseLogger,
        decorator.WithErrorCallback(func(_ context.Context, err error, _ ...any) {
            if err == nil {
                return
            }

            println("ERROR LOGGED:", err.Error())
        }),
    )
    if err != nil {
        return nil, err
    }

    return decorator.NewContext(callbackLogger, fieldsFromCtx)
}
```

Add `NewMask` between `NewContext` and `NewCallback` if hooks must not see secrets.

### Fan-out with Mux 🔗

```go
stdout, err := slogadapter.New(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
if err != nil {
    return nil, err
}

audit, err := slogadapter.New(slog.New(slog.NewJSONHandler(auditFile, nil)))
if err != nil {
    return nil, err
}

return decorator.NewMux(stdout, audit)
```

### Masking Sensitive Data 🔐

```go
func setupSecureLogger(baseLogger cakelog.Logger) (cakelog.Logger, error) {
    return decorator.NewMask(baseLogger, &PatternMasker{
        pattern: regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`),
        replace: "[REDACTED-CARD]",
    })
}
```

Walk nested `map[string]any` values in `MaskArguments` so context fields are redacted too.

### Using Callbacks for Error Handling 🚨

```go
func setupLoggerWithAlerts(baseLogger cakelog.Logger) (cakelog.Logger, error) {
    return decorator.NewCallback(
        baseLogger,
        decorator.WithErrorCallback(func(_ context.Context, err error, _ ...any) {
            if err == nil {
                return
            }

            println("ERROR ALERT:", err.Error())
        }),
        decorator.WithWarnCallback(func(_ context.Context, msg string, _ ...any) {
            println("WARNING ALERT:", msg)
        }),
    )
}
```

## 🤝 Contributing

Contributions are welcome. Please open a pull request.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
