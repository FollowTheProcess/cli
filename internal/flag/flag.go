// Package flag provides a command line flag definition and parsing library.
//
// Flag is intentionally internal so the only interraction is via the Flag option on a command.
package flag

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

const (
	_      = 4 << iota // Unused
	bits8              // 8 bit integer
	bits16             // 16 bit integer
	bits32             // 32 bit integer
	bits64             // 64 bit integer
)

const (
	typeInt         = "int"
	typeInt8        = "int8"
	typeInt16       = "int16"
	typeInt32       = "int32"
	typeInt64       = "int64"
	typeCount       = "count"
	typeUint        = "uint"
	typeUint8       = "uint8"
	typeUint16      = "uint16"
	typeUint32      = "uint32"
	typeUint64      = "uint64"
	typeUintptr     = "uintptr"
	typeFloat32     = "float32"
	typeFloat64     = "float64"
	typeString      = "string"
	typeBool        = "bool"
	typeBytesHex    = "bytesHex"
	typeTime        = "time"
	typeDuration    = "duration"
	typeIP          = "ip"
	typeIntSlice    = "[]int"
	typeInt8Slice   = "[]int8"
	typeInt16Slice  = "[]int16"
	typeInt32Slice  = "[]int32"
	typeInt64Slice  = "[]int64"
	typeUintSlice   = "[]uint"
	typeUint16Slice = "[]uint16"
	typeUint32Slice = "[]uint32"
	typeUint64Slice = "[]uint64"
)

const (
	boolTrue  = "true"
	boolFalse = "false"
)

// NoShortHand should be passed as the "short" argument to [New] if the desired flag
// should be the long hand version only e.g. --count, not -c/--count.
const NoShortHand = rune(-1)

var _ Value = Flag[string]{} // This will fail if we violate our Value interface

// Count is a type used for a flag who's job is to increment a counter, e.g. a "verbosity"
// flag may be passed "-vvv" which should increase the verbosity level to 3.
type Count uint

// Flag represents a single command line flag.
type Flag[T Flaggable] struct {
	value *T     // The actual stored value
	name  string // The name of the flag as appears on the command line, e.g. "force" for a --force flag
	usage string // One line description of the flag, e.g. "Force deletion without confirmation"
	short rune   // Optional shorthand version of the flag, e.g. "f" for a -f flag
}

// New constructs and returns a new [Flag].
//
// The name should be as it appears on the command line, e.g. "force" for a --force flag. An optional
// shorthand can be created by setting short to a single letter value, e.g. "f" to also create a -f version of "force".
//
// If you want the flag to be longhand only, pass "" for short.
//
//	var force bool
//	flag.New(&force, "force", 'f', false, "Force deletion without confirmation")
func New[T Flaggable](p *T, name string, short rune, value T, usage string) (Flag[T], error) {
	if err := validateFlagName(name); err != nil {
		return Flag[T]{}, fmt.Errorf("invalid flag name %q: %w", name, err)
	}

	if err := validateFlagShort(short); err != nil {
		return Flag[T]{}, fmt.Errorf("invalid shorthand for flag %q: %w", name, err)
	}

	if p == nil {
		p = new(T)
	}

	*p = value

	// If the default value is not the zero value for the type, it is treated as
	// significant and shown to the user
	if !isZero(value) {
		usage += fmt.Sprintf(" [default: %v]", value)
	}

	flag := Flag[T]{
		value: p,
		name:  name,
		usage: usage,
		short: short,
	}

	return flag, nil
}

// Name returns the name of the [Flag].
func (f Flag[T]) Name() string {
	return f.name
}

// Short returns the shorthand registered for the flag (e.g. -d for --delete), or
// NoShortHand if the flag should be long only.
func (f Flag[T]) Short() rune {
	return f.short
}

// Usage returns the usage line for the flag.
func (f Flag[T]) Usage() string {
	return f.usage
}

// NoArgValue returns a string representation of value the flag should hold
// when it is given no arguments on the command line. For example a boolean flag
// --delete, when passed without arguments implies --delete true.
func (f Flag[T]) NoArgValue() string {
	switch f.Type() {
	case typeBool:
		// Boolean flags imply passing true, "--force" vs "--force true"
		return boolTrue
	case typeCount:
		// Count flags imply passing 1, "--count --count" or "-cc" should inc by 2
		return "1"
	default:
		return ""
	}
}

// String implements [fmt.Stringer] for a [Flag], and also implements the String
// part of [Value], allowing a flag to print itself.
func (f Flag[T]) String() string { //nolint:cyclop // No other way of doing this realistically
	if f.value == nil {
		return "<nil>"
	}

	switch typ := any(*f.value).(type) {
	case int:
		return formatInt(typ)
	case int8:
		return formatInt(typ)
	case int16:
		return formatInt(typ)
	case int32:
		return formatInt(typ)
	case int64:
		return formatInt(typ)
	case Count:
		return formatUint(typ)
	case uint:
		return formatUint(typ)
	case uint8:
		return formatUint(typ)
	case uint16:
		return formatUint(typ)
	case uint32:
		return formatUint(typ)
	case uint64:
		return formatUint(typ)
	case uintptr:
		return formatUint(typ)
	case float32:
		return formatFloat[float32](bits32)(typ)
	case float64:
		return formatFloat[float64](bits64)(typ)
	case string:
		return typ
	case bool:
		return strconv.FormatBool(typ)
	case []byte:
		return hex.EncodeToString(typ)
	case time.Time:
		return typ.Format(time.RFC3339)
	case time.Duration:
		return typ.String()
	case net.IP:
		return typ.String()
	case []int:
		return formatSlice(typ)
	case []int8:
		return formatSlice(typ)
	case []int16:
		return formatSlice(typ)
	case []int32:
		return formatSlice(typ)
	case []int64:
		return formatSlice(typ)
	case []uint:
		return formatSlice(typ)
	case []uint16:
		return formatSlice(typ)
	case []uint32:
		return formatSlice(typ)
	case []uint64:
		return formatSlice(typ)
	case fmt.Stringer:
		return typ.String()
	default:
		return fmt.Sprintf("ERROR: unhandled type %T", typ)
	}
}

// Type returns a string representation of the type of the Flag.
func (f Flag[T]) Type() string { //nolint:cyclop // No other way of doing this realistically
	if f.value == nil {
		return ""
	}

	switch typ := any(*f.value).(type) {
	case int:
		return typeInt
	case int8:
		return typeInt8
	case int16:
		return typeInt16
	case int32:
		return typeInt32
	case int64:
		return typeInt64
	case Count:
		return typeCount
	case uint:
		return typeUint
	case uint8:
		return typeUint8
	case uint16:
		return typeUint16
	case uint32:
		return typeUint32
	case uint64:
		return typeUint64
	case uintptr:
		return typeUintptr
	case float32:
		return typeFloat32
	case float64:
		return typeFloat64
	case string:
		return typeString
	case bool:
		return typeBool
	case []byte:
		return typeBytesHex
	case time.Time:
		return typeTime
	case time.Duration:
		return typeDuration
	case net.IP:
		return typeIP
	case []int:
		return typeIntSlice
	case []int8:
		return typeInt8Slice
	case []int16:
		return typeInt16Slice
	case []int32:
		return typeInt32Slice
	case []int64:
		return typeInt64Slice
	case []uint:
		return typeUintSlice
	case []uint16:
		return typeUint16Slice
	case []uint32:
		return typeUint32Slice
	case []uint64:
		return typeUint64Slice
	default:
		return fmt.Sprintf("%T", typ)
	}
}

// Set sets a [Flag] value based on string input, i.e. parsing from the command line.
func (f Flag[T]) Set(str string) error { //nolint:gocognit,cyclop // No other way of doing this realistically
	if f.value == nil {
		return fmt.Errorf("cannot set value %s, flag.value was nil", str)
	}

	switch typ := any(*f.value).(type) {
	case int:
		val, err := parseInt[int](0)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case int8:
		val, err := parseInt[int8](bits8)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case int16:
		val, err := parseInt[int16](bits16)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case int32:
		val, err := parseInt[int32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case int64:
		val, err := parseInt[int64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case Count:
		// We have to do a bit of custom stuff here as an incremement is a read and write op
		// First read the current value of the flag and cast it to a Count so we
		// can increment it
		current, ok := any(*f.value).(Count)
		if !ok {
			// This basically shouldn't ever happen but it's easy enough to handle nicely
			return errBadType(*f.value)
		}

		// Increment the count and store it back
		newValue := current + 1
		*f.value = *cast[T](&newValue)

		return nil
	case uint:
		val, err := parseUint[uint](0)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case uint8:
		val, err := parseUint[uint8](bits8)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case uint16:
		val, err := parseUint[uint16](bits16)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case uint32:
		val, err := parseUint[uint32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case uint64:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case uintptr:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case float32:
		val, err := parseFloat[float32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case float64:
		val, err := parseFloat[float64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case string:
		val := str
		*f.value = *cast[T](&val)

		return nil
	case bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case []byte:
		val, err := hex.DecodeString(strings.TrimSpace(str))
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case time.Time:
		val, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case time.Duration:
		val, err := time.ParseDuration(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}

		*f.value = *cast[T](&val)

		return nil
	case net.IP:
		val := net.ParseIP(str)
		if val == nil {
			return errParse(f.name, str, typ, errors.New("invalid IP address"))
		}

		*f.value = *cast[T](&val)

		return nil
	case []int:
		// Like Count, a slice flag is a read/write op
		slice, ok := any(*f.value).([]int)
		if !ok {
			return errBadType(*f.value)
		}

		// Append the given value to the slice
		newValue, err := parseInt[int](0)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	case []int8:
		slice, ok := any(*f.value).([]int8)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parseInt[int8](bits8)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	case []int16:
		slice, ok := any(*f.value).([]int16)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parseInt[int16](bits16)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	case []int32:
		slice, ok := any(*f.value).([]int32)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parseInt[int32](bits32)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	case []int64:
		slice, ok := any(*f.value).([]int64)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parseInt[int64](bits64)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil

	case []uint:
		slice, ok := any(*f.value).([]uint)
		if !ok {
			return errBadType(*f.value)
		}

		// Append the given value to the slice
		newValue, err := parseUint[uint](0)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	case []uint16:
		slice, ok := any(*f.value).([]uint16)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parseUint[uint16](bits16)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	case []uint32:
		slice, ok := any(*f.value).([]uint32)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parseUint[uint32](bits32)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	case []uint64:
		slice, ok := any(*f.value).([]uint64)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParseSlice(f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *cast[T](&slice)

		return nil
	default:
		return fmt.Errorf("unsupported flag type: %T", typ)
	}
}

// signed is the same as constraints.Signed but we don't have to depend
// on golang/x/exp.
type signed interface {
	int | int8 | int16 | int32 | int64
}

// unsigned is the same as constraints.Unsigned (with Count mixed in) but we don't have to depend
// on golang/x/exp.
type unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64 | uintptr | Count
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
func cast[T2, T1 any](v *T1) *T2 {
	return (*T2)(unsafe.Pointer(v))
}

// validateFlagName ensures a flag name is valid, returning an error if it's not.
//
// Flags names must be all lower case ASCII letters, a hyphen separator is allowed e.g. "set-default"
// but this must be in between letters, not leading or trailing.
func validateFlagName(name string) error {
	if name == "" {
		return errors.New("must not be empty")
	}

	before, after, found := strings.Cut(name, "-")

	// Hyphen must be in between "words" like "set-default"
	// we can't have "-default" or "default-"
	if found && after == "" {
		return errors.New("trailing hyphen")
	}

	if found && before == "" {
		return errors.New("leading hyphen")
	}

	for _, char := range name {
		// No whitespace
		if unicode.IsSpace(char) {
			return errors.New("cannot contain whitespace")
		}
		// Only ASCII characters allowed
		if char > unicode.MaxASCII {
			return fmt.Errorf("contains non ascii character: %q", string(char))
		}
		// Only non-letter character allowed is a hyphen
		if !unicode.IsLetter(char) && char != '-' {
			return fmt.Errorf("contains non ascii letter: %q", string(char))
		}
		// Any upper case letters are not allowed
		if unicode.IsLetter(char) && !unicode.IsLower(char) {
			return fmt.Errorf("contains upper case character %q", string(char))
		}
	}

	return nil
}

// validateFlagShort ensures a flag's shorthand is valid, returning an error if it's not.
//
// Flag shorthands must be a single ASCII letter, we use rune as the type here rather than string because
// it enforces only a single character, so all we have to do is make sure it's a valid ASCII character.
func validateFlagShort(short rune) error {
	// If it's the marker for long hand only, this is fine
	if short == NoShortHand {
		return nil
	}

	if unicode.IsSpace(short) {
		return errors.New("cannot contain whitespace")
	}
	// Shorthand must be a valid ASCII letter
	if short == utf8.RuneError || short > unicode.MaxASCII || !unicode.IsLetter(short) {
		return fmt.Errorf("invalid character, must be a single ASCII letter, got %q", string(short))
	}

	return nil
}

// errParse is a helper to quickly return a consistent error in the face of flag
// value parsing errors.
func errParse[T Flaggable](name, str string, typ T, err error) error {
	return fmt.Errorf(
		"flag %q received invalid value %q (expected %T), detail: %w",
		name,
		str,
		typ,
		err,
	)
}

// errParseSlice is like errParse but for []T flags where the error message
// needs to be a bit more specific.
func errParseSlice[T Flaggable](name, str string, typ T, err error) error {
	return fmt.Errorf(
		"flag %q (type %T) cannot append element %q: %w",
		name,
		typ,
		str,
		err,
	)
}

// errBadType makes a consistent error in the face of a bad type
// assertion.
func errBadType[T Flaggable](value T) error {
	return fmt.Errorf("bad value %v, could not cast to %T", value, value)
}

// parseInt is a generic helper to parse all signed integers, given a bit size.
//
// It returns the parsed value or an error.
func parseInt[T signed](bits int) func(str string) (T, error) {
	return func(str string) (T, error) {
		val, err := strconv.ParseInt(str, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

// parseUint is a generic helper to parse all signed integers, given a bit size.
//
// It returns the parsed value or an error.
func parseUint[T unsigned](bits int) func(str string) (T, error) {
	return func(str string) (T, error) {
		val, err := strconv.ParseUint(str, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

// parseFloat is a generic helper to parse floating point numbers, given a bit size.
//
// It returns the parsed value or an error.
func parseFloat[T ~float32 | ~float64](bits int) func(str string) (T, error) {
	return func(str string) (T, error) {
		val, err := strconv.ParseFloat(str, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

// formatInt is a generic helper to return a string representation of any signed integer.
func formatInt[T signed](in T) string {
	return strconv.FormatInt(int64(in), 10)
}

// formatUint is a generic helper to return a string representation of any unsigned integer.
func formatUint[T unsigned](in T) string {
	return strconv.FormatUint(uint64(in), 10)
}

// formatFloat is a generic helper to return a string representation of any floating point digit.
func formatFloat[T ~float32 | ~float64](bits int) func(T) string {
	return func(in T) string {
		return strconv.FormatFloat(float64(in), 'g', -1, bits)
	}
}

// formatSlice is a generic helper to return a string representation of a slice.
func formatSlice[T any](slice []T) string {
	return fmt.Sprintf("%v", slice)
}

// isZero reports whether value is the zero value for it's type.
func isZero[T Flaggable](value T) bool {
	switch typ := any(value).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return typ == 0
	case Count:
		return typ == Count(0)
	case string:
		return typ == ""
	case bool:
		return !typ
	case []byte:
		return typ == nil
	case time.Time:
		var zero time.Time
		return typ.Equal(zero)
	case time.Duration:
		var zero time.Duration
		return typ == zero
	case net.IP:
		return typ == nil
	case []int:
		return typ == nil
	default:
		return false
	}
}
