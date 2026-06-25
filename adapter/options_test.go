package adapter_test

import (
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog/adapter"
)

// TestWithArgsKey_ErrorNilOptions verifies that WithArgsKey option returns an error when applied to nil Options.
func TestWithArgsKey_ErrorNilOptions(t *testing.T) {
	t.Parallel()

	opt := adapter.WithArgsKey("")

	err := opt(nil)
	if err == nil {
		t.Fatal("expected error when providing nil options, got nil")
	}

	if !errors.Is(err, adapter.ErrNilOptions) {
		t.Errorf("expected error to be ErrNilOptions, got %v", err)
	}
}

// TestWithArgsKey_ErrorEmptyArgsKey verifies that WithArgsKey option returns
// an error when provided with an empty args key.
func TestWithArgsKey_ErrorEmptyArgsKey(t *testing.T) {
	t.Parallel()

	opt := adapter.WithArgsKey("")

	err := opt(new(adapter.Options))
	if err == nil {
		t.Fatal("expected error when providing empty args key, got nil")
	}

	if !errors.Is(err, adapter.ErrEmptyArgsKey) {
		t.Errorf("expected error to be ErrEmptyArgsKey, got %v", err)
	}
}

// TestWithArgsKey_Success verifies that WithArgsKey option successfully sets the args key for various input strings.
func TestWithArgsKey_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		key  string
	}{
		{name: "numeric_key", key: "1"},
		{name: "single_character", key: "a"},
		{name: "two_characters", key: "ab"},
		{name: "short_string", key: "ctx"},
		{name: "medium_string", key: "args_context"},
		{name: "long_string", key: "this_is_a_very_long_key_with_many_characters"},
		{name: "only_numbers", key: "123456789"},
		{name: "only_lowercase", key: "abcdefghijklmnopqrstuvwxyz"},
		{name: "only_uppercase", key: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{name: "mixed_case", key: "ArGsKeY"},
		{name: "with_underscore", key: "args_key"},
		{name: "with_hyphen", key: "args-key"},
		{name: "with_dot", key: "args.key"},
		{name: "with_special_chars", key: "args!@#$%^&*()"},
		{name: "with_unicode", key: "ключ_args"},
		{name: "with_emoji", key: "key_🔑"},
		{name: "with_spaces", key: "args key"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opt := adapter.WithArgsKey(tc.key)
			opts := adapter.DefaultOptions()

			err := opt(opts)
			if err != nil {
				t.Errorf("expected no error for key '%s', got %v", tc.key, err)
			}
		})
	}
}
