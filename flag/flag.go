// Package flag provides an experimental command line flag definition and parsing library.
//
// CLI currently uses [spf13/pflag] for flag parsing (like Cobra), this package is an attempt at defining
// a new approach with some of the tools we now have in modern Go. It is not intended to be backwards compatible
// with pflag or the std lib flag package.
//
// [spf13/pflag]: https://github.com/spf13/pflag
package flag

import (
	"fmt"
	"net"
	"time"
	"unsafe"
)

const (
	_      = 4 << iota // Unused
	bits8              // 8 bit integer
	bits16             // 16 bit integer
	bits32             // 32 bit integer
	bits64             // 64 bit integer
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
//	var force bool
//	flag.New(&force, "force", "f", false, "Force deletion without confirmation")
func New[T Flaggable](p *T, name string, short string, value T, usage string) *Flag[T] {
	if p == nil {
		p = new(T)
	}
	*p = value
	flag := &Flag[T]{
		value: value,
		name:  name,
		usage: usage,
		short: short,
	}
	return flag
}

// Get gets a [Flag] value.
func (f *Flag[T]) Get() T {
	return f.value
}

// Set sets a [Flag] value based on string input, i.e. parsing from the command line.
func (f *Flag[T]) Set( //nolint:gocyclo // No other way of doing this realistically
	str string,
) error {
	switch typ := any(f.value).(type) {
	case int:
		val, err := parseInt[int](0)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case int8:
		val, err := parseInt[int8](bits8)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case int16:
		val, err := parseInt[int16](bits16)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case int32:
		val, err := parseInt[int32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case int64:
		val, err := parseInt[int64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case uint:
		val, err := parseUint[uint](0)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case uint8:
		val, err := parseUint[uint8](bits8)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case uint16:
		val, err := parseUint[uint16](bits16)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case uint32:
		val, err := parseUint[uint32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case uint64:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case uintptr:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case float32:
		val, err := parseFloat[float32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case float64:
		val, err := parseFloat[float64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ)
		}
		f.value = *cast[T](&val)
		return nil
	case string:
		val := str
		f.value = *cast[T](&val)
		return nil
	default:
		return fmt.Errorf("unsupported flag type: %T", typ)
	}
}

// cast converts a *T1 to a *T2, we use it here when we know (via generics and compile time checks)
// that e.g. the Flag.value is a string, but we can't directly do Flag.value = "value" because
// we can't assign a string to a generic 'T', but we *know* that the value *is* a string because when
// instantiating a Flag[T], you have to provide (or compiler has to infer) Flag[string].
//
// # Safety
//
// This function uses [unsafe.Pointer] underneath to reassign the types but we know this is safe to do
// based on the compile time checks provided by generics. Further, it fits the following valid pattern
// specified in the docs for [unsafe.Pointer].
//
// Conversion of a *T1 to Pointer to *T2
//
// Provided that T2 is no larger than T1 and that the two share an equivalent
// memory layout, this conversion allows reinterpreting data of one type as
// data of another type.
//
// This describes our use case as we're converting a *T to e.g a *string but *only* when we know
// that a Flag[T] is actually Flag[string], so the memory layout and size is guaranteed by the
// compiler to be equivalent.
func cast[T2 any, T1 any](v *T1) *T2 {
	return (*T2)(unsafe.Pointer(v))
}
