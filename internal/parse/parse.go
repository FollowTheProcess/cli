// Package parse provides functions to parse strings into Go types and produce
// detailed, consistent errors.
//
// It is used across both internal/flag and internal/arg to provide consistency.
package parse

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

// Kind is the kind of parsing being done, either argument or flag.
type Kind string

// Err is a generic parse error.
//
// Errors returned from the [Error] function will match this in a call
// to [errors.Is].
var Err = errors.New("parse error")

const (
	// KindArgument is the [Kind] used for argument parsing.
	KindArgument Kind = "argument"

	// KindFlag is the [Kind] used for flag parsing.
	KindFlag Kind = "flag"
)

const (
	bits8  = 8 << iota // 8 bit integer
	bits16             // 16 bit integer
	bits32             // 32 bit integer
	bits64             // 64 bit integer
)

const base10 = 10

// Error produces a formatted parse error.
//
// The kind should must be [KindArgument] or [KindFlag], with name and str being the
// name of the arg/flag and the invalid text that triggered the error.
//
// The type T is the type we were parsing str into and err is any underlying
// error e.g. from strconv.
//
//	// Make a flag parse error
//	var force bool
//	return parse.Error(parse.KindArgument, "force", "faklse", force, strconv.ErrSyntax)
func Error[T any](kind Kind, name, str string, typ T, err error) error {
	// Ordinarily I wouldn't have a package like this concern itself with
	// details of other packages (like flag/arg) but given this package exists to produce consistent
	// behaviour and clear error messages in the narrow context of this cli framework, then
	// it makes sense the error is defined here too.
	return fmt.Errorf("%w: %s %q received invalid value %q (expected %T): %w", Err, kind, name, str, typ, err)
}

// Int parses an int from a string.
func Int(str string) (int, error) {
	val, err := strconv.ParseInt(str, base10, 0)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}

// Int8 parses an int8 from a string.
func Int8(str string) (int8, error) {
	val, err := strconv.ParseInt(str, base10, bits8)
	if err != nil {
		return 0, err
	}

	return int8(val), nil
}

// Int16 parses an int16 from a string.
func Int16(str string) (int16, error) {
	val, err := strconv.ParseInt(str, base10, bits16)
	if err != nil {
		return 0, err
	}

	return int16(val), nil
}

// Int32 parses an int32 from a string.
func Int32(str string) (int32, error) {
	val, err := strconv.ParseInt(str, base10, bits32)
	if err != nil {
		return 0, err
	}

	return int32(val), nil
}

// Int64 parses an int64 from a string.
func Int64(str string) (int64, error) {
	val, err := strconv.ParseInt(str, base10, bits64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Uint parses a uint from a string.
func Uint(str string) (uint, error) {
	val, err := strconv.ParseUint(str, base10, 0)
	if err != nil {
		return 0, err
	}

	return uint(val), nil
}

// Uint8 parses an uint8 from a string.
func Uint8(str string) (uint8, error) {
	val, err := strconv.ParseUint(str, base10, bits8)
	if err != nil {
		return 0, err
	}

	return uint8(val), nil
}

// Uint16 parses an uint16 from a string.
func Uint16(str string) (uint16, error) {
	val, err := strconv.ParseUint(str, base10, bits16)
	if err != nil {
		return 0, err
	}

	return uint16(val), nil
}

// Uint32 parses an uint32 from a string.
func Uint32(str string) (uint32, error) {
	val, err := strconv.ParseUint(str, base10, bits32)
	if err != nil {
		return 0, err
	}

	return uint32(val), nil
}

// Uint64 parses an uint64 from a string.
func Uint64(str string) (uint64, error) {
	val, err := strconv.ParseUint(str, base10, bits64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// Float32 parses a float32 from a string.
func Float32(str string) (float32, error) {
	val, err := strconv.ParseFloat(str, bits32)
	if err != nil {
		return 0, err
	}

	return float32(val), nil
}

// Float64 parses a float64 from a string.
func Float64(str string) (float64, error) {
	val, err := strconv.ParseFloat(str, bits64)
	if err != nil {
		return 0, err
	}

	return float64(val), nil
}

// Cast converts a *T1 to a *T2, we use it here when we know (via generics and compile time checks)
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
func Cast[T2, T1 any](v *T1) *T2 {
	return (*T2)(unsafe.Pointer(v))
}
