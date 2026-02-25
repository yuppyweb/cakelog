# 🍰 Cakelog

**Cakelog** is a flexible Go logging library that provides a unified logger interface with support for many popular logging frameworks through adapters and extensible functionality through decorators.

## 🎯 Core Idea

Cakelog solves the problem of binding code to a specific logging library. Instead of using one logger directly, you work with a unified `Logger` interface that can be adapted to any popular logger or combine multiple loggers simultaneously through decorators.

```
┌─────────────────────────────┐
│     Your Application        │
└────────────┬────────────────┘
             │
             ▼
     ┌───────────────┐
     │ cakelog.Logger│  ◄───── Unified Interface
     └───────┬───────┘
             │
    ┌────────┴────────┬──────────┬────────────┐
    ▼                 ▼          ▼            ▼
 Adapters:      Logrus       Slog       Zap       Zerolog
    │
 Decorators: Context  Prometheus  Sentry
```

## 📦 Logger Interface

The main interface contains four logging methods:

```go
type Logger interface {
    Debug(ctx context.Context, msg string, args ...any)
    Info(ctx context.Context, msg string, args ...any)
    Warn(ctx context.Context, msg string, args ...any)
    Error(ctx context.Context, err error, args ...any)
}
```

All methods:
- ✅ Accept `context.Context` for passing values between components
- ✅ Support variable number of arguments for additional information
- ✅ Work with async operation context

---

## 🔌 Adapters

Adapters allow you to use popular logging libraries as `cakelog.Logger`.

### 📊 Logrus Adapter

An adapter for [sirupsen/logrus](https://github.com/sirupsen/logrus) — one of the most popular loggers in Go.

```go
import (
    "context"
    "github.com/sirupsen/logrus"
    "github.com/yuppyweb/cakelog/adapter"
)

func main() {
    logrusLogger := logrus.New()
    logger := adapter.NewLogrusLogger(logrusLogger)
    
    ctx := context.Background()
    logger.Info(ctx, "Application started", map[string]any{"version": "1.0", "environment": "prod"})
}
```

**Features:**
- Context support via `WithContext()`
- Arguments stored in `context` field (by default)
- Flexible `ArgsKey` parameter for field name customization

```go
logger := adapter.NewLogrusLogger(logrusLogger)
logger.ArgsKey = "metadata"  // Change the key for arguments
```

---

### 🔮 Slog Adapter

An adapter for the built-in [log/slog](https://pkg.go.dev/log/slog) (available from Go 1.21+).

```go
import (
    "context"
    "log/slog"
    "github.com/yuppyweb/cakelog/adapter"
)

func main() {
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger := adapter.NewSlogLogger(slogLogger)
    
    ctx := context.Background()
    logger.Error(ctx, errors.New("database error"), map[string]any{"query": "SELECT * FROM users"})
}
```

**Features:**
- Built into the standard library
- Structured logging in JSON/Text format
- Full context support via `*Context` methods

---

### ⚡ Zap Adapter

An adapter for [Uber's Zap](https://github.com/uber-go/zap) — a high-performance logger.

```go
import (
    "context"
    "github.com/yuppyweb/cakelog/adapter"
    "go.uber.org/zap"
)

func main() {
    zapLogger, _ := zap.NewProduction()
    defer zapLogger.Sync()
    
    logger := adapter.NewZapLogger(zapLogger)
    
    ctx := context.Background()
    logger.Warn(ctx, "High memory usage", map[string]any{"usage": "85%", "threshold": "80%"})
}
```

**Features:**
- Extremely fast logging
- Minimal data copying
- Optimal for high-load applications

---

### 📬 Zerolog Adapter

An adapter for [rs/zerolog](https://github.com/rs/zerolog) — a reflection-free logger.

```go
import (
    "context"
    "github.com/rs/zerolog"
    "github.com/yuppyweb/cakelog/adapter"
)

func main() {
    zerologLogger := zerolog.New(os.Stdout)
    logger := adapter.NewZerologLogger(&zerologLogger)
    
    ctx := context.Background()
    logger.Debug(ctx, "Sending request", map[string]any{"url": "https://api.example.com", "timeout": "30s"})
}
```

**Features:**
- Minimalist and fast
- Good balance between performance and functionality
- Without reflection usage

---

## 🎨 Decorators

Decorators extend logger functionality by wrapping an existing `cakelog.Logger`.

### 🎁 Context Decorator

Enriches logs with values stored in the context. Allows passing data between functions through context.

```go
import (
    "context"
    "github.com/yuppyweb/cakelog/decorator"
    "github.com/yuppyweb/cakelog/adapter"
)

func processOrder(ctx context.Context, logger cakelog.Logger) {
    // Create context logger
    ctxLogger := decorator.NewContextLogger(logger)
    
    // Add values to context
    ctx = ctxLogger.PutContext(ctx, "order_id", "12345")
    ctx = ctxLogger.PutContext(ctx, "user_id", "user_789")
    
    // These values will be automatically added to all logs
    ctxLogger.Info(ctx, "Processing order")
    // Output: Info: Processing order order_id=12345 user_id=user_789
    
    processPayment(ctx, ctxLogger)
}

func processPayment(ctx context.Context, logger cakelog.Logger) {
    // Context values are still available!
    logger.Info(ctx, "Payment processed")
    // Output: Info: Payment processed order_id=12345 user_id=user_789
}
```

**Use cases:**
- Tracking request ID (request ID) through the entire call chain
- Preserving user information for log correlation
- Passing context between goroutines

---

### 📈 Prometheus Decorator

Tracks the number of logs at each level using Prometheus metrics.

```go
import (
    "context"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/yuppyweb/cakelog/decorator"
    "github.com/yuppyweb/cakelog/adapter"
)

func setupLogger(baseLogger cakelog.Logger) cakelog.Logger {
    // Create counters for each logging level
    counters := decorator.PrometheusLoggerCounter{
        Debug: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "logs_debug_total",
            Help: "Total debug logs",
        }),
        Info: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "logs_info_total",
            Help: "Total info logs",
        }),
        Warn: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "logs_warn_total",
            Help: "Total warning logs",
        }),
        Error: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "logs_error_total",
            Help: "Total error logs",
        }),
    }
    
    // Register counters
    prometheus.MustRegister(counters.Debug, counters.Info, counters.Warn, counters.Error)
    
    // Wrap logger in Prometheus decorator
    return decorator.NewPrometheusLogger(baseLogger, counters)
}

func main() {
    logger := setupLogger(adapter.NewSlogLogger(slog.Default()))
    ctx := context.Background()
    
    logger.Info(ctx, "Server started")
    logger.Warn(ctx, "Low memory")
    logger.Error(ctx, errors.New("connection error"))
    
    // Metrics in Prometheus:
    // logs_info_total 1
    // logs_warn_total 1
    // logs_error_total 1
}
```

**Benefits:**
- 🔍 Monitoring log frequency
- 📊 Problem detection through anomalies
- 🚨 Setting up alerts based on error count

---

### 🔴 Sentry Decorator

Integrates [Sentry](https://sentry.io) for error tracking and events in production.

```go
import (
    "context"
    "github.com/getsentry/sentry-go"
    "github.com/yuppyweb/cakelog/decorator"
    "github.com/yuppyweb/cakelog/adapter"
)

func setupLogger(baseLogger cakelog.Logger) cakelog.Logger {
    // Initialize Sentry
    sentry.Init(sentry.ClientOptions{
        Dsn: "https://your-key@sentry.io/your-project-id",
    })
    
    // Create Sentry hubs for each level (optional)
    hubs := decorator.SentryLoggerHub{
        Debug: sentry.CurrentHub(),
        Info:  sentry.CurrentHub(),
        Warn:  sentry.CurrentHub(),
        Error: sentry.CurrentHub(),
    }
    
    return decorator.NewSentryLogger(baseLogger, hubs)
}

func main() {
    logger := setupLogger(adapter.NewSlogLogger(slog.Default()))
    ctx := context.Background()
    
    // Will be sent to Sentry and to the base logger
    logger.Error(ctx, errors.New("DB connection expired"), map[string]any{"database": "postgres"})
    
    // Output in logs:
    // Error: DB connection expired sentryEventId=<uuid>
}
```

**Capabilities:**
- 🎯 Automatic exception capturing
- 🏷️ Tags and context for each event
- 🔄 Deduplication of similar errors
- 📧 Alerts on critical errors

---

## 🧩 Combining Adapters and Decorators

The main advantage of Cakelog is the ability to combine components:

```go
import (
    "context"
    "log/slog"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/getsentry/sentry-go"
    "github.com/yuppyweb/cakelog"
    "github.com/yuppyweb/cakelog/adapter"
    "github.com/yuppyweb/cakelog/decorator"
)

func setupProductionLogger() cakelog.Logger {
    // 1 Create base logger through adapter
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    baseLogger := adapter.NewSlogLogger(slogLogger)
    
    // 2 Add context enrichment
    ctxLogger := decorator.NewContextLogger(baseLogger)
    
    // 3 Add Prometheus metrics
    promountCounters := decorator.PrometheusLoggerCounter{
        Debug: /* ... */,
        Info:  /* ... */,
        // ...
    }
    promLogger := decorator.NewPrometheusLogger(ctxLogger, promountCounters)
    
    // 4 Add Sentry integration
    sentryHubs := decorator.SentryLoggerHub{
        Debug: sentry.CurrentHub(),
        Info:  sentry.CurrentHub(),
        Warn:  sentry.CurrentHub(),
        Error: sentry.CurrentHub(),
    }
    finalLogger := decorator.NewSentryLogger(promLogger, sentryHubs)
    
    return finalLogger
}

func main() {
    logger := setupProductionLogger()
    ctx := context.Background()
    
    // Processing chain:
    // logger.Info() 
    //   → Sentry will capture (if configured)
    //   → Prometheus counter increases
    //   → Context will be added
    //   → Slog outputs JSON
    
    logger.Info(ctx, "Server ready", map[string]any{"port": 8080})
}
```

---

## 📋 NopLogger

Built-in logger that does nothing. Useful in tests or to disable logging:

```go
import "github.com/yuppyweb/cakelog"

func main() {
    nopLogger := cakelog.NewNopLogger()
    
    // All calls will be ignored
    nopLogger.Info(context.Background(), "This will not be logged")
}
```

---

## 🚀 Installation

```bash
go get github.com/yuppyweb/cakelog
```

Then install the needed adapters:

```bash
go get github.com/sirupsen/logrus          # For Logrus
go get github.com/rs/zerolog               # For Zerolog
go get go.uber.org/zap                     # For Zap
go get github.com/prometheus/client_golang # For Prometheus
go get github.com/getsentry/sentry-go      # For Sentry
```

---

## 💡 Usage Examples

### Example 1: Microservice with Context Tracking

```go
func handleUserRequest(ctx context.Context, userID string, logger cakelog.Logger) {
    ctxLogger := decorator.NewContextLogger(logger)
    ctx = ctxLogger.PutContext(ctx, "user_id", userID)
    ctx = ctxLogger.PutContext(ctx, "request_id", generateRequestID())
    
    ctxLogger.Info(ctx, "User request received")
    
    user, err := fetchUser(ctx, ctxLogger)
    if err != nil {
        ctxLogger.Error(ctx, err, map[string]any{"action": "fetch_user"})
        return
    }
    
    ctxLogger.Info(ctx, "User loaded", map[string]any{"name": user.Name})
}
```

### Example 2: Using Different Loggers for Different Scenarios

```go
// For development
func devLogger() cakelog.Logger {
    slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    return adapter.NewSlogLogger(slogLogger)
}

// For production
func prodLogger() cakelog.Logger {
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    baseLogger := adapter.NewSlogLogger(slogLogger)
    return decorator.NewSentryLogger(baseLogger, /* ... */)
}

func main() {
    var logger cakelog.Logger
    if os.Getenv("ENV") == "production" {
        logger = prodLogger()
    } else {
        logger = devLogger()
    }
    
    logger.Info(context.Background(), "Application started")
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

For more information about the MIT License, visit [opensource.org/licenses/MIT](https://opensource.org/licenses/MIT).