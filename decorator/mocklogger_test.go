package decorator_test

import (
	"context"

	"github.com/yuppyweb/cakelog"
)

type mockMsgArgs struct {
	ctx  context.Context
	msg  string
	args []any
}

type mockErrArgs struct {
	ctx  context.Context
	err  error
	args []any
}

type mockLogger struct {
	debugIn []mockMsgArgs
	infoIn  []mockMsgArgs
	warnIn  []mockMsgArgs
	errorIn []mockErrArgs
}

func (ml *mockLogger) Debug(ctx context.Context, msg string, args ...any) {
	ml.debugIn = append(ml.debugIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

func (ml *mockLogger) Info(ctx context.Context, msg string, args ...any) {
	ml.infoIn = append(ml.infoIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

func (ml *mockLogger) Warn(ctx context.Context, msg string, args ...any) {
	ml.warnIn = append(ml.warnIn, mockMsgArgs{
		ctx:  ctx,
		msg:  msg,
		args: args,
	})
}

func (ml *mockLogger) Error(ctx context.Context, err error, args ...any) {
	ml.errorIn = append(ml.errorIn, mockErrArgs{
		ctx:  ctx,
		err:  err,
		args: args,
	})
}

var _ cakelog.Logger = (*mockLogger)(nil)
