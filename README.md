# 🍰 Cakelog

[![Go Version](https://img.shields.io/github/go-mod/go-version/yuppyweb/cakelog)](https://github.com/yuppyweb/cakelog)
[![Go Report Card](https://goreportcard.com/badge/github.com/yuppyweb/cakelog)](https://goreportcard.com/report/github.com/yuppyweb/cakelog)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**Cakelog** is a flexible Go logging library that provides a unified logger interface with support for multiple popular logging frameworks (Logrus, Slog, Zap, Zerolog) through adapters and extensible functionality through decorators.

## ✨ Features

✅ **Unified Logger Interface** — Write code once, switch loggers anytime  
✅ **Multiple Adapters** — Logrus, Slog, Zap, Zerolog support  
✅ **Extensible Decorators** — Context enrichment, Callback hooks, Data masking  
✅ **Context-Based Enrichment** — Pass metadata through context values  
✅ **Callback Hooks** — Execute custom logic on log events  
✅ **Data Masking** — Automatically mask sensitive information  
✅ **Composable Pattern** — Chain decorators for advanced functionality  
✅ **NopLogger for Testing** — Built-in no-op implementation 

## 🎯 Core Idea

Cakelog solves the problem of binding code to a specific logging library. Instead of depending on one logger directly, you work with a unified `Logger` interface that can be adapted to any popular logger or combine multiple loggers simultaneously through decorators.

```
┌─ Your Application
│
├─ cakelog.Logger ◄──── Unified Interface
│
├─ Adapters (choose one):
│  • Logrus
│  • Slog (Go 1.21+)
│  • Zap
│  • Zerolog
│
└─ Decorators (stack as needed):
   • Context (enrichment)
   • Callback (hooks)
   • Mask (sanitization)
```

---

## 📦 Installation

```bash
go get github.com/yuppyweb/cakelog
```

Install adapters as needed:

```bash
go get github.com/sirupsen/logrus     # For Logrus adapter
go get github.com/rs/zerolog          # For Zerolog adapter
go get go.uber.org/zap                # For Zap adapter
# Slog is built-in to Go 1.21+
```

## 🚀 Quick Start

### Basic Usage with Slog 📚

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/yuppyweb/cakelog/adapter"
)

func main() {
    // Create base logger
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger, _ := adapter.NewSlogLogger(slogLogger)
    
    ctx := context.Background()
    
    // Use logger
    logger.Info(ctx, "Application started")
    logger.Debug(ctx, "Debug message")
    logger.Warn(ctx, "Warning message")
    logger.Error(ctx, errors.New("connection failed"))
}
```

## 📖 Core Components

### Logger Interface 📋

The main interface contains four logging methods:

```go
type Logger interface {
    // Debug logs a debug-level message
    Debug(ctx context.Context, msg string, args ...any)
    // Info logs an info-level message
    Info(ctx context.Context, msg string, args ...any)
    // Warn logs a warning-level message
    Warn(ctx context.Context, msg string, args ...any)
    // Error logs an error-level message with an error value
    Error(ctx context.Context, err error, args ...any)
}
```

### NopLogger 🚫

Built-in no-operation logger for testing:

```go
logger := cakelog.NewNopLogger() // discards all log messages
logger.Info(ctx, "This won't be logged")
```

## 🔌 Adapters

Cakelog provides adapters for the most popular Go logging frameworks:

### Slog Adapter (Go 1.21+) 🔮

```go
import (
    "log/slog"
    "github.com/yuppyweb/cakelog/adapter"
)

slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger, err := adapter.NewSlogLogger(slogLogger)
if err != nil {
    log.Fatal(err)
}
```

### Zap Adapter ⚡

```go
import (
    "go.uber.org/zap"
    "github.com/yuppyweb/cakelog/adapter"
)

zapLogger, _ := zap.NewProduction()
logger, err := adapter.NewZapLogger(zapLogger)
if err != nil {
    log.Fatal(err)
}
```

### Logrus Adapter 📊

```go
import (
    "github.com/sirupsen/logrus"
    "github.com/yuppyweb/cakelog/adapter"
)

logrusLogger := logrus.New()
logger, err := adapter.NewLogrusLogger(logrusLogger)
if err != nil {
    log.Fatal(err)
}
```

### Zerolog Adapter 📬

```go
import (
    "github.com/rs/zerolog"
    "github.com/yuppyweb/cakelog/adapter"
)

zerologLogger := zerolog.New(os.Stdout)
logger, err := adapter.NewZerologLogger(zerologLogger)
if err != nil {
    log.Fatal(err)
}
```

### Custom Args Key ⚙️

All adapters support customizing the key under which arguments are stored:

```go
logger, err := adapter.NewSlogLogger(slogLogger, adapter.WithArgsKey("fields"))
```

## 🎨 Decorators

Decorators allow you to enhance logger functionality:

### Context Decorator (Enrichment) 🎁

Add metadata to all logs automatically using context values:

```go
import (
    "github.com/yuppyweb/cakelog/decorator"
)

// Add values to context
ctx := decorator.WithContextValue(context.Background(), "user_id", "user123")
ctx = decorator.WithContextValue(ctx, "request_id", "req456")

// Get all stored values
values := decorator.ContextValues(ctx)
// values: {"user_id": "user123", "request_id": "req456"}

// Get single value
userID := decorator.ContextValue(ctx, "user_id")
```

### Callback Decorator (Hooks) ↩️

Execute custom logic on each log event:

```go
import (
    "github.com/yuppyweb/cakelog/decorator"
)

// Define callbacks
debugCallback := func(ctx context.Context) { /* track debug events */ }
errorCallback := func(ctx context.Context) { /* alert on errors */ }

// Create callback logger
callbackLogger, err := decorator.NewCallbackLogger(
    baseLogger,
    decorator.WithDebugCallback(debugCallback),
    decorator.WithErrorCallback(errorCallback),
)
```

### Mask Decorator (Sanitization) 🎭

Mask sensitive data in all logs:

```go
import (
    "github.com/yuppyweb/cakelog/decorator"
)

// Define masker rules
type MyMasker struct{}

func (m MyMasker) MaskMessage(msg string) string {
    return strings.ReplaceAll(msg, "password", "***")
}

func (m MyMasker) MaskError(err error) error {
    return errors.New(strings.ReplaceAll(err.Error(), "secret", "***"))
}

func (m MyMasker) MaskArgument(arg any) any {
    if str, ok := arg.(string); ok {
        return strings.ReplaceAll(str, "token", "xxx")
    }
    return arg
}

// Create mask logger
maskLogger, err := decorator.NewMaskLogger(baseLogger, MyMasker{})
```

### Composing Decorators 🧩

Stack multiple decorators together:

```go
baseLogger, _ := adapter.NewSlogLogger(slogLogger)
callbackLogger, _ := decorator.NewCallbackLogger(baseLogger, opts...)
maskLogger, _ := decorator.NewMaskLogger(callbackLogger, masker)

// Now maskLogger has both callback and masking capabilities
```

## 📝 Examples

### Real-World Example: Service with Context Enrichment 🏗️

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/yuppyweb/cakelog/adapter"
    "github.com/yuppyweb/cakelog/decorator"
)

func processRequest(logger cakelog.Logger, requestID string, userID string) {
    // Enrich context with request metadata
    ctx := decorator.WithContextValue(context.Background(), "request_id", requestID)
    ctx = decorator.WithContextValue(ctx, "user_id", userID)
    
    logger.Info(ctx, "Request processing started")
    
    // ... do work ...
    
    logger.Info(ctx, "Request processing completed")
}

func main() {
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger, _ := adapter.NewSlogLogger(slogLogger)
    
    processRequest(logger, "req-123", "user-456")
}
```

### Combining Adapters and Decorators 🧩

The main advantage of Cakelog is the ability to combine components:

```go
import (
    "context"
    "log/slog"
    "os"
    "github.com/yuppyweb/cakelog"
    "github.com/yuppyweb/cakelog/adapter"
    "github.com/yuppyweb/cakelog/decorator"
)

func setupProductionLogger() (cakelog.Logger, error) {
    // 1. Create base logger through adapter
    slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    baseLogger, err := adapter.NewSlogLogger(slogLogger)
    if err != nil {
        return nil, err
    }
    
    // 2. Add callback hooks
    callbackLogger, err := decorator.NewCallbackLogger(
        baseLogger,
        decorator.WithErrorCallback(func(ctx context.Context) {
            // Trigger alert on error
            println("ERROR LOGGED!")
        }),
    )
    if err != nil {
        return nil, err
    }
    
    return callbackLogger, nil
}

func main() {
    logger, _ := setupProductionLogger()
    ctx := context.Background()
    
    // Add context values
    ctx = decorator.WithContextValue(ctx, "request_id", "req_12345")
    ctx = decorator.WithContextValue(ctx, "user_id", "user_789")
    
    // Log with enriched context
    logger.Info(ctx, "Server ready")
}

### Microservice with Context Tracking 🔗

```go
import (
    "context"
    "github.com/yuppyweb/cakelog/decorator"
)

func handleUserRequest(ctx context.Context, userID string, logger cakelog.Logger) {
    // Add correlation values to context
    ctx = decorator.WithContextValue(ctx, "user_id", userID)
    ctx = decorator.WithContextValue(ctx, "request_id", generateRequestID())
    
    logger.Info(ctx, "User request received")
    
    user, err := fetchUser(ctx, logger)
    if err != nil {
        logger.Error(ctx, err, map[string]any{"action": "fetch_user"})
        return
    }
    
    logger.Info(ctx, "User loaded", map[string]any{"name": user.Name})
}
```

### Masking Sensitive Data 🔐

```go
import (
    "regexp"
    "github.com/yuppyweb/cakelog/decorator"
    "github.com/yuppyweb/cakelog/adapter"
)

type PatternMasker struct {
    pattern *regexp.Regexp
    replace string
}

func (pm *PatternMasker) MaskMessage(msg string) string {
    return pm.pattern.ReplaceAllString(msg, pm.replace)
}

func (pm *PatternMasker) MaskError(err error) error {
    return err
}

func (pm *PatternMasker) MaskArgument(arg any) any {
    return arg
}

func setupSecureLogger(baseLogger cakelog.Logger) (cakelog.Logger, error) {
    maskers := []decorator.Masker{
        &PatternMasker{
            pattern: regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`),
            replace: "[REDACTED-CARD]",
        },
    }
    
    return decorator.NewMaskLogger(baseLogger, maskers...)
}
```

###  Using Callbacks for Error Handling 🚨

```go
import (
    "context"
    "github.com/yuppyweb/cakelog/decorator"
)

func setupLoggerWithAlerts(baseLogger cakelog.Logger) (cakelog.Logger, error) {
    callbackLogger, err := decorator.NewCallbackLogger(
        baseLogger,
        decorator.WithErrorCallback(func(ctx context.Context) {
            // Send alert when error is logged
            println("ERROR ALERT!")
        }),
        decorator.WithWarnCallback(func(ctx context.Context) {
            // Track warning
            println("WARNING ALERT!")
        }),
    )
    if err != nil {
        return nil, err
    }
    
    return callbackLogger, nil
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

For more information about the MIT License, visit [opensource.org/licenses/MIT](https://opensource.org/licenses/MIT).