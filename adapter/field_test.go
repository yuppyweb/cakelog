package adapter_test

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/yuppyweb/cakelog/adapter"
)

func TestFields_Empty(t *testing.T) {
	t.Parallel()

	got := adapter.Fields(nil)
	if got != nil {
		t.Errorf("expected nil fields, got %v", got)
	}
}

func TestFields_KeyValuePairs(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{"method", "GET", "status", 200})

	want := []adapter.Field{
		{Key: "method", Value: "GET"},
		{Key: "status", Value: 200},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_MapAny(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{map[string]any{"method": "GET", "status": 200}})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "method", Value: "GET"},
		{Key: "status", Value: 200},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_MapInt(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{map[string]int{"status": 200, "user_id": 42}})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "status", Value: 200},
		{Key: "user_id", Value: 42},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_MapInt32(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{map[string]int32{"status": 200, "user_id": 42}})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "status", Value: int32(200)},
		{Key: "user_id", Value: int32(42)},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_MapFloat64(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{map[string]float64{"duration": 1.5, "ratio": 0.25}})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "duration", Value: 1.5},
		{Key: "ratio", Value: 0.25},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_KnownStringMaps(t *testing.T) {
	t.Parallel()

	errVal := errors.New("boom")

	testCases := []struct {
		name string
		arg  any
		want any
	}{
		{name: "bool", arg: map[string]bool{"v": true}, want: true},
		{name: "error", arg: map[string]error{"v": errVal}, want: errVal},
		{name: "int8", arg: map[string]int8{"v": 8}, want: int8(8)},
		{name: "int16", arg: map[string]int16{"v": 16}, want: int16(16)},
		{name: "int64", arg: map[string]int64{"v": 64}, want: int64(64)},
		{name: "uint", arg: map[string]uint{"v": 1}, want: uint(1)},
		{name: "uint8", arg: map[string]uint8{"v": 8}, want: uint8(8)},
		{name: "uint16", arg: map[string]uint16{"v": 16}, want: uint16(16)},
		{name: "uint32", arg: map[string]uint32{"v": 32}, want: uint32(32)},
		{name: "uint64", arg: map[string]uint64{"v": 64}, want: uint64(64)},
		{name: "uintptr", arg: map[string]uintptr{"v": 8}, want: uintptr(8)},
		{name: "float32", arg: map[string]float32{"v": 1.5}, want: float32(1.5)},
		{name: "complex64", arg: map[string]complex64{"v": 1 + 2i}, want: complex64(1 + 2i)},
		{name: "complex128", arg: map[string]complex128{"v": 3 + 4i}, want: complex128(3 + 4i)},
		{
			name: "string_slice",
			arg:  map[string][]string{"v": {"admin", "user"}},
			want: []string{"admin", "user"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := adapter.Fields([]any{tc.arg})
			if len(got) != 1 {
				t.Fatalf("expected 1 field, got %#v", got)
			}

			want := []adapter.Field{
				{Key: "v", Value: tc.want},
			}

			if !equalFields(got, want) {
				t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
			}
		})
	}
}

func TestFields_MapString(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{map[string]string{"method": "GET", "path": "/health"}})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "method", Value: "GET"},
		{Key: "path", Value: "/health"},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_Mixed(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{
		"method", "GET",
		map[string]any{"status": 200, "user_id": 42},
	})

	if len(got) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "method", Value: "GET"},
		{Key: "status", Value: 200},
		{Key: "user_id", Value: 42},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_MapAsPairValue(t *testing.T) {
	t.Parallel()

	meta := map[string]any{"ip": "127.0.0.1"}
	got := adapter.Fields([]any{"meta", meta})

	if len(got) != 1 {
		t.Fatalf("expected 1 field, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "meta", Value: meta},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_TrailingValue(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{"method", "GET", "orphan"})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "method", Value: "GET"},
		{Key: "arg", Value: "orphan"},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_NonStringKey(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{42, "value"})

	if len(got) != 1 {
		t.Fatalf("expected 1 field, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "42", Value: "value"},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_NonStringKeyMap(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{map[int]string{1: "a"}, "value"})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "1", Value: "a"},
		{Key: "arg", Value: "value"},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_DefinedStringKeyMap(t *testing.T) {
	t.Parallel()

	type methodKey string

	got := adapter.Fields([]any{map[methodKey]int{"status": 200}})

	want := []adapter.Field{
		{Key: "status", Value: 200},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_NilMap(t *testing.T) {
	t.Parallel()

	var values map[string]any

	got := adapter.Fields([]any{values})

	if len(got) != 0 {
		t.Errorf("expected no fields from nil map, got %#v", got)
	}
}

func TestFields_DuplicateKeys(t *testing.T) {
	t.Parallel()

	got := adapter.Fields([]any{"status", 200, "status", 500})

	if len(got) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(got))
	}

	want := []adapter.Field{
		{Key: "status", Value: 200},
		{Key: "status", Value: 500},
	}

	if !equalFields(got, want) {
		t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
	}
}

func TestFields_DuplicateKeyValues(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		args []any
		want []adapter.Field
	}{
		{
			name: "pairs",
			args: []any{"status", 200, "status", 200},
			want: []adapter.Field{
				{Key: "status", Value: 200},
				{Key: "status", Value: 200},
			},
		},
		{
			name: "pair_and_map",
			args: []any{"status", 200, map[string]int{"status": 200}},
			want: []adapter.Field{
				{Key: "status", Value: 200},
				{Key: "status", Value: 200},
			},
		},
		{
			name: "maps",
			args: []any{
				map[string]int{"status": 200},
				map[string]int{"status": 200},
			},
			want: []adapter.Field{
				{Key: "status", Value: 200},
				{Key: "status", Value: 200},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := adapter.Fields(tc.args)
			if len(got) != len(tc.want) {
				t.Fatalf("expected %d fields, got %d", len(tc.want), len(got))
			}

			if !equalFields(got, tc.want) {
				t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, tc.want)
			}
		})
	}
}

func TestFields_MapKeyTypes(t *testing.T) {
	t.Parallel()

	type methodKey string

	type id struct {
		N int
	}

	testCases := []struct {
		name  string
		arg   any
		key   string
		value any
	}{
		{name: "string", arg: map[string]string{"method": "GET"}, key: "method", value: "GET"},
		{name: "defined_string", arg: map[methodKey]int{"status": 200}, key: "status", value: 200},
		{name: "any_string", arg: map[any]string{"method": "GET"}, key: "method", value: "GET"},
		{name: "any_int", arg: map[any]string{42: "answer"}, key: "42", value: "answer"},
		{name: "bool", arg: map[bool]string{true: "yes"}, key: "true", value: "yes"},
		{name: "int", arg: map[int]string{1: "a"}, key: "1", value: "a"},
		{name: "int8", arg: map[int8]string{8: "v"}, key: "8", value: "v"},
		{name: "int16", arg: map[int16]string{16: "v"}, key: "16", value: "v"},
		{name: "int32", arg: map[int32]string{32: "v"}, key: "32", value: "v"},
		{name: "int64", arg: map[int64]string{64: "v"}, key: "64", value: "v"},
		{name: "uint", arg: map[uint]string{1: "v"}, key: "1", value: "v"},
		{name: "uint8", arg: map[uint8]string{8: "v"}, key: "8", value: "v"},
		{name: "uint16", arg: map[uint16]string{16: "v"}, key: "16", value: "v"},
		{name: "uint32", arg: map[uint32]string{32: "v"}, key: "32", value: "v"},
		{name: "uint64", arg: map[uint64]string{64: "v"}, key: "64", value: "v"},
		{name: "uintptr", arg: map[uintptr]string{8: "v"}, key: "8", value: "v"},
		{name: "float32", arg: map[float32]string{1.5: "v"}, key: "1.5", value: "v"},
		{name: "float64", arg: map[float64]string{1.5: "v"}, key: "1.5", value: "v"},
		{name: "complex64", arg: map[complex64]string{1 + 2i: "v"}, key: "(1+2i)", value: "v"},
		{name: "complex128", arg: map[complex128]string{3 + 4i: "v"}, key: "(3+4i)", value: "v"},
		{name: "array", arg: map[[2]int]string{{1, 2}: "v"}, key: "[1 2]", value: "v"},
		{name: "struct", arg: map[id]string{{N: 1}: "v"}, key: "{1}", value: "v"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := adapter.Fields([]any{tc.arg})
			if len(got) != 1 {
				t.Fatalf("expected 1 field, got %#v", got)
			}

			want := []adapter.Field{
				{Key: tc.key, Value: tc.value},
			}

			if !equalFields(got, want) {
				t.Errorf("unexpected fields:\nGot:  %#v\nWant: %#v", got, want)
			}
		})
	}
}

func equalFields(got, want []adapter.Field) bool {
	return reflect.DeepEqual(sortedFields(got), sortedFields(want))
}

func sortedFields(fields []adapter.Field) []adapter.Field {
	sorted := slices.Clone(fields)

	slices.SortStableFunc(sorted, func(left, right adapter.Field) int {
		if keyCmp := cmp.Compare(left.Key, right.Key); keyCmp != 0 {
			return keyCmp
		}

		return cmp.Compare(fmt.Sprintf("%#v", left.Value), fmt.Sprintf("%#v", right.Value))
	})

	return sorted
}
