package decorator

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/yuppyweb/cakelog"
)

// Is the default key used to store the Sentry event ID in the log arguments.
const DefaultSentryEventIDKey = "sentryEventId"

// Holds separate Sentry hubs for different log levels.
// This allows for more granular control over which logs are sent to Sentry and how they are processed.
type SentryLoggerHub struct {
	// Debug level hub for capturing debug messages in Sentry.
	Debug *sentry.Hub

	// Info level hub for capturing informational messages in Sentry.
	Info *sentry.Hub

	// Warn level hub for capturing warning messages in Sentry.
	Warn *sentry.Hub

	// Error level hub for capturing error messages in Sentry.
	Error *sentry.Hub
}

// SentryLogger is a decorator that sends logs to Sentry and then forwards them to the underlying logger.
type SentryLogger struct {
	// The underlying logger to which logs will be forwarded after being sent to Sentry.
	log cakelog.Logger

	// Containing separate hubs for each log level.
	hub SentryLoggerHub

	// The key used to store the Sentry event ID in the log arguments.
	SentryEventIDKey string
}

// Creates a new SentryLogger with the provided underlying logger and Sentry hubs.
func NewSentryLogger(log cakelog.Logger, hub SentryLoggerHub) *SentryLogger {
	return &SentryLogger{
		log:              log,
		hub:              hub,
		SentryEventIDKey: DefaultSentryEventIDKey,
	}
}

// Sends a debug message to Sentry and then forwards it to the underlying logger.
func (sl *SentryLogger) Debug(ctx context.Context, msg string, args ...any) {
	sl.log.Debug(ctx, msg, sl.captureMessage(sl.hub.Debug, msg, args)...)
}

// Sends an info message to Sentry and then forwards it to the underlying logger.
func (sl *SentryLogger) Info(ctx context.Context, msg string, args ...any) {
	sl.log.Info(ctx, msg, sl.captureMessage(sl.hub.Info, msg, args)...)
}

// Sends a warning message to Sentry and then forwards it to the underlying logger.
func (sl *SentryLogger) Warn(ctx context.Context, msg string, args ...any) {
	sl.log.Warn(ctx, msg, sl.captureMessage(sl.hub.Warn, msg, args)...)
}

// Sends an error message to Sentry and then forwards it to the underlying logger.
func (sl *SentryLogger) Error(ctx context.Context, err error, args ...any) {
	sl.log.Error(ctx, err, sl.captureException(sl.hub.Error, err, args)...)
}

// Helper method to capture a message in Sentry and return the updated log arguments with the Sentry event ID.
func (sl *SentryLogger) captureMessage(hub *sentry.Hub, msg string, args []any) []any {
	if hub == nil {
		return args
	}

	if eventID := hub.CaptureMessage(msg); eventID != nil {
		return append(args, map[string]*sentry.EventID{sl.SentryEventIDKey: eventID})
	}

	return args
}

// Helper method to capture an exception in Sentry and return the updated log arguments with the Sentry event ID.
func (sl *SentryLogger) captureException(hub *sentry.Hub, err error, args []any) []any {
	if hub == nil {
		return args
	}

	if eventID := hub.CaptureException(err); eventID != nil {
		return append(args, map[string]*sentry.EventID{sl.SentryEventIDKey: eventID})
	}

	return args
}

// Ensures that SentryLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*SentryLogger)(nil)
