package decorator

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yuppyweb/cakelog"
)

// Holds Prometheus counters for each log level.
type PrometheusLoggerCounter struct {
	// Counter for debug log messages.
	Debug prometheus.Counter

	// Counter for info log messages.
	Info prometheus.Counter

	// Counter for warning log messages.
	Warn prometheus.Counter

	// Counter for error log messages.
	Error prometheus.Counter
}

// Is a cakelog.Logger decorator that increments Prometheus counters for each log level.
type PrometheusLogger struct {
	// The underlying logger to which log messages will be forwarded.
	logger cakelog.Logger

	// The Prometheus counters to be incremented for each log level.
	counter PrometheusLoggerCounter
}

// Creates a new PrometheusLogger that wraps the provided cakelog.Logger
// and uses the given PrometheusLoggerCounter to track log levels.
func NewPrometheusLogger(logger cakelog.Logger, counter PrometheusLoggerCounter) *PrometheusLogger {
	return &PrometheusLogger{
		logger:  logger,
		counter: counter,
	}
}

// Sends a debug message to the underlying logger and increments the debug counter if it is set.
func (pl *PrometheusLogger) Debug(ctx context.Context, msg string, args ...any) {
	pl.logger.Debug(ctx, msg, args...)

	if pl.counter.Debug != nil {
		pl.counter.Debug.Inc()
	}
}

// Sends an info message to the underlying logger and increments the info counter if it is set.
func (pl *PrometheusLogger) Info(ctx context.Context, msg string, args ...any) {
	pl.logger.Info(ctx, msg, args...)

	if pl.counter.Info != nil {
		pl.counter.Info.Inc()
	}
}

// Sends a warning message to the underlying logger and increments the warn counter if it is set.
func (pl *PrometheusLogger) Warn(ctx context.Context, msg string, args ...any) {
	pl.logger.Warn(ctx, msg, args...)

	if pl.counter.Warn != nil {
		pl.counter.Warn.Inc()
	}
}

// Sends an error message to the underlying logger and increments the error counter if it is set.
func (pl *PrometheusLogger) Error(ctx context.Context, err error, args ...any) {
	pl.logger.Error(ctx, err, args...)

	if pl.counter.Error != nil {
		pl.counter.Error.Inc()
	}
}

// Ensures that PrometheusLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*PrometheusLogger)(nil)
