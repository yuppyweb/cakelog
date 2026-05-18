package adapter

import (
	"errors"
)

// This file defines common options and constants for the logger adapters in this package.
const defaultArgsKey = "context"

var (
	ErrNilOptions   = errors.New("is nil options")
	ErrEmptyArgsKey = errors.New("is empty args key")
)

// Option represents a configuration option for the logger adapters.
// It is a function that takes a pointer to an Options struct and returns an error if the option is invalid.
type Option func(*Options) error

// Options represents the configuration options for the logger adapters.
type Options struct {
	// The key under which the context arguments will be stored in log entries.
	argsKey string
}

// DefaultOptions returns a new Options struct initialized with default values.
func DefaultOptions() *Options {
	return &Options{
		argsKey: defaultArgsKey,
	}
}

// WithArgsKey returns an Option that sets the key under which context arguments will be stored in log entries.
func WithArgsKey(key string) Option {
	return func(opts *Options) error {
		if opts == nil {
			return ErrNilOptions
		}

		if key == "" {
			return ErrEmptyArgsKey
		}

		opts.argsKey = key

		return nil
	}
}
