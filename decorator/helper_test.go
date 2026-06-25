package decorator_test

import (
	"context"

	"github.com/yuppyweb/cakelog"
)

// mockMsgArgs holds context, message, and arguments for a message log call.
type mockMsgArgs struct {
	ctx  context.Context
	msg  string
	args []any
}

// mockErrArgs holds context, error, and arguments for an error log call.
type mockErrArgs struct {
	ctx  context.Context
	err  error
	args []any
}

// mockLogger is a mock implementation of cakelog.Logger for testing.
type mockLogger struct {
	debugIn []mockMsgArgs
	infoIn  []mockMsgArgs
	warnIn  []mockMsgArgs
	errorIn []mockErrArgs
}

// Debug records a debug-level message with context and arguments.
func (ml *mockLogger) Debug(ctx context.Context, msg string, args ...any) {
	ml.debugIn = append(ml.debugIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

// Info records an info-level message with context and arguments.
func (ml *mockLogger) Info(ctx context.Context, msg string, args ...any) {
	ml.infoIn = append(ml.infoIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

// Warn records a warning-level message with context and arguments.
func (ml *mockLogger) Warn(ctx context.Context, msg string, args ...any) {
	ml.warnIn = append(ml.warnIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

// Error records an error-level message with context, error, and arguments.
func (ml *mockLogger) Error(ctx context.Context, err error, args ...any) {
	ml.errorIn = append(ml.errorIn, mockErrArgs{
		ctx:  ctx,
		err:  err,
		args: args,
	})
}

// Ensure mockLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*mockLogger)(nil)
