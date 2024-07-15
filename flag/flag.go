// Package flag implements CLI flag defining and parsing functionality.
package flag

import (
	"net"
	"time"
)

// Flaggable is a type constraint that defines any type capable of being parsed as a command line flag.
type Flaggable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string | ~bool | ~[]byte | time.Time | net.IPNet
}

// Flag represents a single command line flag.
type Flag[T Flaggable] struct {
	value T      // The actual stored value
	name  string // The name of the flag as appears on the command line, e.g. "force" for a --force flag
	usage string // One line description of the flag, e.g. "Force deletion without confirmation"
	short string // Optional shorthand version of the flag, e.g. "f" for a -f flag
}

// New constructs and returns a new [Flag].
//
// The name should be as it appears on the command line, e.g. "force" for a --force flag. An optional
// shorthand can be created by setting short to a single letter value, e.g. "f" to also create a -f version of "force".
//
// If you want the flag to be longhand only, pass "" for short.
//
//	flag.New("force", "f", false, "Force deletion without confirmation")
func New[T Flaggable](name string, short string, value T, usage string) *Flag[T] {
	// Default implementation
	flag := &Flag[T]{
		value: value,
		name:  name,
		usage: usage,
		short: short,
	}
	return flag
}

// Value returns the flags current value.
func (f *Flag[T]) Value() T {
	return f.value
}
