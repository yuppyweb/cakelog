package decorator

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yuppyweb/cakelog"
)

type PrometheusLoggerCounter struct {
	Debug prometheus.Counter
	Info  prometheus.Counter
	Warn  prometheus.Counter
	Error prometheus.Counter
}

type PrometheusLogger struct {
	logger  cakelog.Logger
	counter PrometheusLoggerCounter
}

func NewPrometheusLogger(logger cakelog.Logger, counter PrometheusLoggerCounter) *PrometheusLogger {
	return &PrometheusLogger{
		logger:  logger,
		counter: counter,
	}
}

func (pl *PrometheusLogger) Debug(ctx context.Context, msg string, args ...any) {
	pl.logger.Debug(ctx, msg, args...)

	if pl.counter.Debug != nil {
		pl.counter.Debug.Inc()
	}
}

func (pl *PrometheusLogger) Info(ctx context.Context, msg string, args ...any) {
	pl.logger.Info(ctx, msg, args...)

	if pl.counter.Info != nil {
		pl.counter.Info.Inc()
	}
}

func (pl *PrometheusLogger) Warn(ctx context.Context, msg string, args ...any) {
	pl.logger.Warn(ctx, msg, args...)

	if pl.counter.Warn != nil {
		pl.counter.Warn.Inc()
	}
}

func (pl *PrometheusLogger) Error(ctx context.Context, err error, args ...any) {
	pl.logger.Error(ctx, err, args...)

	if pl.counter.Error != nil {
		pl.counter.Error.Inc()
	}
}

var _ cakelog.Logger = (*PrometheusLogger)(nil)
