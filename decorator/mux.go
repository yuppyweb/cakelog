package decorator

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/yuppyweb/cakelog"
)

// ErrNoLoggerProvided is returned by NewMux when logs is empty.
var ErrNoLoggerProvided = errors.New("no logger provided")

// muxLogger is a decorator that forwards each log call to every
// underlying logger in order.
//
// The same ctx and message are passed to each logger. The args slice is
// copied per logger, not the values it contains. A panic in one logger
// is not recovered and stops remaining loggers.
type muxLogger struct {
	// logs is the ordered list of loggers that receive each call.
	logs []cakelog.Logger
}

// NewMux wraps logs so each log call is forwarded to every logger in
// order. The loggers slice is copied; later mutations of the caller's
// slice do not affect the mux. Passing the same logger more than once
// forwards the call to it on every occurrence.
//
// An empty logs list returns a wrapped ErrNoLoggerProvided. A nil logger,
// including a typed nil such as a nil pointer stored in Logger,
// returns a wrapped ErrNilLogger for that 1-based position.
//
// The same ctx and message are passed to each logger. The args slice is
// copied per logger, not the values it contains. A panic in one logger
// is not recovered and stops remaining loggers.
func NewMux(logs ...cakelog.Logger) (cakelog.Logger, error) {
	if len(logs) == 0 {
		return nil, fmt.Errorf("mux logger: %w", ErrNoLoggerProvided)
	}

	for i, log := range logs {
		if err := requireLogger(log); err != nil {
			return nil, fmt.Errorf("mux logger: %d: %w", i+1, err)
		}
	}

	copiedLogs := make([]cakelog.Logger, len(logs))
	copy(copiedLogs, logs)

	return &muxLogger{logs: copiedLogs}, nil
}

// Debug forwards the message to every underlying logger in order.
func (ml *muxLogger) Debug(ctx context.Context, msg string, args ...any) {
	for _, log := range ml.logs {
		log.Debug(ctx, msg, slices.Clone(args)...)
	}
}

// Info forwards the message to every underlying logger in order.
func (ml *muxLogger) Info(ctx context.Context, msg string, args ...any) {
	for _, log := range ml.logs {
		log.Info(ctx, msg, slices.Clone(args)...)
	}
}

// Warn forwards the message to every underlying logger in order.
func (ml *muxLogger) Warn(ctx context.Context, msg string, args ...any) {
	for _, log := range ml.logs {
		log.Warn(ctx, msg, slices.Clone(args)...)
	}
}

// Error forwards the error to every underlying logger in order.
func (ml *muxLogger) Error(ctx context.Context, err error, args ...any) {
	for _, log := range ml.logs {
		log.Error(ctx, err, slices.Clone(args)...)
	}
}

// Ensure muxLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*muxLogger)(nil)
