package decorator

import (
	"errors"
	"reflect"

	"github.com/yuppyweb/cakelog"
)

// ErrNilLogger is returned by decorator constructors when log is nil,
// including a typed nil such as a nil pointer stored in Logger.
var ErrNilLogger = errors.New("logger is nil")

func requireLogger(log cakelog.Logger) error {
	if isNil(log) {
		return ErrNilLogger
	}

	return nil
}

func isNil(value any) bool {
	if value == nil {
		return true
	}

	val := reflect.ValueOf(value)

	//nolint:exhaustive // IsNil is valid only for nillable kinds.
	switch val.Kind() {
	case reflect.Pointer,
		reflect.Interface,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer:
		if val.IsNil() {
			return true
		}
	}

	return false
}
