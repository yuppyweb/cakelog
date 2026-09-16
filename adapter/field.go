package adapter

import (
	"fmt"
	"reflect"
)

const (
	// extraArgKey is the field name for a trailing unpaired argument.
	extraArgKey = "arg"
	// keyValuePairSize is the number of args consumed as one key-value pair.
	keyValuePairSize = 2
)

// Field is a single structured log Field.
type Field struct {
	// Key is the field name.
	Key string
	// Value is the field value.
	Value any
}

// Fields parses mixed slog-style key-value pairs and maps into a stable field list.
//
// Arguments are read left to right:
//   - map values are expanded; non-string keys are formatted with fmt.Sprint;
//   - otherwise a string key is paired with the next value;
//   - a non-string key is formatted with fmt.Sprint;
//   - a trailing value without a key is stored under "arg";
//   - a map used as a pair value is kept as a single field;
//   - duplicate keys are kept in encounter order, including equal values.
//
// Adapters emit every returned field. What a backend then emits is its own
// field model: slog and zap keep every occurrence; logrus and zerolog keep
// the last value.
func Fields(args []any) []Field {
	if len(args) == 0 {
		return nil
	}

	parsed := make([]Field, 0, len(args))

	i := 0
	for i < len(args) {
		parsed, i = appendArg(parsed, args, i)
	}

	return parsed
}

func appendArg(parsed []Field, args []any, idx int) ([]Field, int) {
	if expanded, ok := appendStringMap(parsed, args[idx]); ok {
		return expanded, idx + 1
	}

	if idx+1 >= len(args) {
		return append(parsed, Field{Key: extraArgKey, Value: args[idx]}), idx + 1
	}

	return append(
		parsed,
		Field{Key: stringifyKey(args[idx]), Value: args[idx+1]},
	), idx + keyValuePairSize
}

func appendStringMap(parsed []Field, arg any) ([]Field, bool) {
	value := reflect.ValueOf(arg)
	if value.Kind() != reflect.Map {
		return parsed, false
	}

	for iter := value.MapRange(); iter.Next(); {
		parsed = append(parsed, Field{
			Key:   stringifyKey(iter.Key().Interface()),
			Value: iter.Value().Interface(),
		})
	}

	return parsed, true
}

func stringifyKey(key any) string {
	if s, ok := key.(string); ok {
		return s
	}

	return fmt.Sprint(key)
}
