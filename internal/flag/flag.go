// Package flag provides a command line flag definition and parsing library.
//
// Flag is intentionally internal so the only interaction is via the Flag option on a command.
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

	"go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/cli/internal/parse"
)

var _ Value = Flag[string]{} // This will fail if we violate our Value interface

// Flag represents a single command line flag.
type Flag[T flag.Flaggable] struct {
	value *T     // The actual stored value
	name  string // The name of the flag as appears on the command line, e.g. "force" for a --force flag
	usage string // One line description of the flag, e.g. "Force deletion without confirmation"
	short rune   // Optional shorthand version of the flag, e.g. "f" for a -f flag
}

// New constructs and returns a new [Flag].
//
// The name should be as it appears on the command line, e.g. "force" for a --force flag. An optional
// shorthand can be created by setting short to a single letter value, e.g. "f" to also create a -f version of "force".
func New[T flag.Flaggable](p *T, name string, short rune, usage string, config Config[T]) (Flag[T], error) {
	if err := validateFlagName(name); err != nil {
		return Flag[T]{}, fmt.Errorf("invalid flag name %q: %w", name, err)
	}

	if err := validateFlagShort(short); err != nil {
		return Flag[T]{}, fmt.Errorf("invalid shorthand for flag %q: %w", name, err)
	}

	if p == nil {
		p = new(T)
	}

	*p = config.DefaultValue

	flag := Flag[T]{
		value: p,
		name:  name,
		usage: usage,
		short: short,
	}

	// TODO(@FollowTheProcess): This needs to live in command.go and we should iterate over the flagset
	// adding the values to the tabwriter as we go rather than relying on the flagset.Usage() method
	// to provide *all* the usage

	// If the default value is not the zero value for the type, it is treated as
	// significant and shown to the user
	flag.usage += "\t"
	if !isZeroIsh(*p) {
		// \t so that defaults get aligned by tabwriter when the command
		// dumps the flags
		flag.usage += fmt.Sprintf("[default: %s]", flag.String())
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
	case format.TypeBool:
		// Boolean flags imply passing true, "--force" vs "--force true"
		return format.True
	case format.TypeCount:
		// Count flags imply passing 1, "--count --count" or "-cc" should inc by 2
		return "1"
	default:
		return ""
	}
}

// String implements [fmt.Stringer] for a [Flag], and also implements the String
// part of [Value], allowing a flag to print itself.
//
//nolint:cyclop // No other way of doing this realistically
func (f Flag[T]) String() string {
	if f.value == nil {
		return "<nil>"
	}

	switch typ := any(*f.value).(type) {
	case int:
		return format.Int(typ)
	case int8:
		return format.Int(typ)
	case int16:
		return format.Int(typ)
	case int32:
		return format.Int(typ)
	case int64:
		return format.Int(typ)
	case flag.Count:
		return format.Uint(typ)
	case uint:
		return format.Uint(typ)
	case uint8:
		return format.Uint(typ)
	case uint16:
		return format.Uint(typ)
	case uint32:
		return format.Uint(typ)
	case uint64:
		return format.Uint(typ)
	case uintptr:
		return format.Uint(typ)
	case float32:
		return format.Float32(typ)
	case float64:
		return format.Float64(typ)
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
		return format.Slice(typ)
	case []int8:
		return format.Slice(typ)
	case []int16:
		return format.Slice(typ)
	case []int32:
		return format.Slice(typ)
	case []int64:
		return format.Slice(typ)
	case []uint:
		return format.Slice(typ)
	case []uint16:
		return format.Slice(typ)
	case []uint32:
		return format.Slice(typ)
	case []uint64:
		return format.Slice(typ)
	case []float32:
		return format.Slice(typ)
	case []float64:
		return format.Slice(typ)
	case []string:
		return format.Slice(typ)
	default:
		return fmt.Sprintf("Flag.String: unsupported flag type: %T", typ)
	}
}

// Type returns a string representation of the type of the Flag.
func (f Flag[T]) Type() string { //nolint:cyclop // No other way of doing this realistically
	if f.value == nil {
		return "<nil>"
	}

	switch typ := any(*f.value).(type) {
	case int:
		return format.TypeInt
	case int8:
		return format.TypeInt8
	case int16:
		return format.TypeInt16
	case int32:
		return format.TypeInt32
	case int64:
		return format.TypeInt64
	case flag.Count:
		return format.TypeCount
	case uint:
		return format.TypeUint
	case uint8:
		return format.TypeUint8
	case uint16:
		return format.TypeUint16
	case uint32:
		return format.TypeUint32
	case uint64:
		return format.TypeUint64
	case uintptr:
		return format.TypeUintptr
	case float32:
		return format.TypeFloat32
	case float64:
		return format.TypeFloat64
	case string:
		return format.TypeString
	case bool:
		return format.TypeBool
	case []byte:
		return format.TypeBytesHex
	case time.Time:
		return format.TypeTime
	case time.Duration:
		return format.TypeDuration
	case net.IP:
		return format.TypeIP
	case []int:
		return format.TypeIntSlice
	case []int8:
		return format.TypeInt8Slice
	case []int16:
		return format.TypeInt16Slice
	case []int32:
		return format.TypeInt32Slice
	case []int64:
		return format.TypeInt64Slice
	case []uint:
		return format.TypeUintSlice
	case []uint16:
		return format.TypeUint16Slice
	case []uint32:
		return format.TypeUint32Slice
	case []uint64:
		return format.TypeUint64Slice
	case []float32:
		return format.TypeFloat32Slice
	case []float64:
		return format.TypeFloat64Slice
	default:
		return fmt.Sprintf("%T", typ)
	}
}

// Set sets a [Flag] value based on string input, i.e. parsing from the command line.
//
//nolint:gocognit,maintidx // No other way of doing this realistically
func (f Flag[T]) Set(str string) error {
	if f.value == nil {
		return fmt.Errorf("cannot set value %s, flag.value was nil", str)
	}

	switch typ := any(*f.value).(type) {
	case int:
		val, err := parse.Int(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case int8:
		val, err := parse.Int8(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case int16:
		val, err := parse.Int16(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case int32:
		val, err := parse.Int32(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case int64:
		val, err := parse.Int64(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case flag.Count:
		// We have to do a bit of custom stuff here as an increment is a read and write op
		// First read the current value of the flag and cast it to a Count so we
		// can increment it
		current, ok := any(*f.value).(flag.Count)
		if !ok {
			// This basically shouldn't ever happen but it's easy enough to handle nicely
			return errBadType(*f.value)
		}

		// Add the count and store it back, we still parse the given str rather
		// than just +1 every time as this allows people to do e.g. --verbosity=3
		// as well as -vvv
		val, err := parse.Uint(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		newValue := current + flag.Count(val)
		*f.value = *parse.Cast[T](&newValue)

		return nil
	case uint:
		val, err := parse.Uint(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case uint8:
		val, err := parse.Uint8(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case uint16:
		val, err := parse.Uint16(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case uint32:
		val, err := parse.Uint32(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case uint64:
		val, err := parse.Uint64(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case uintptr:
		val, err := parse.Uint64(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case float32:
		val, err := parse.Float32(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case float64:
		val, err := parse.Float64(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case string:
		val := str
		*f.value = *parse.Cast[T](&val)

		return nil
	case bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case []byte:
		val, err := hex.DecodeString(strings.TrimSpace(str))
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case time.Time:
		val, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case time.Duration:
		val, err := time.ParseDuration(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case net.IP:
		val := net.ParseIP(str)
		if val == nil {
			return parse.Error(parse.KindFlag, f.name, str, typ, errors.New("invalid IP address"))
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case []int:
		// Like Count, a slice flag is a read/write op
		slice, ok := any(*f.value).([]int)
		if !ok {
			return errBadType(*f.value)
		}

		// Append the given value to the slice
		newValue, err := parse.Int(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []int8:
		slice, ok := any(*f.value).([]int8)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Int8(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []int16:
		slice, ok := any(*f.value).([]int16)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Int16(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []int32:
		slice, ok := any(*f.value).([]int32)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Int32(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []int64:
		slice, ok := any(*f.value).([]int64)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Int64(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil

	case []uint:
		slice, ok := any(*f.value).([]uint)
		if !ok {
			return errBadType(*f.value)
		}

		// Append the given value to the slice
		newValue, err := parse.Uint(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []uint16:
		slice, ok := any(*f.value).([]uint16)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Uint16(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []uint32:
		slice, ok := any(*f.value).([]uint32)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Uint32(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []uint64:
		slice, ok := any(*f.value).([]uint64)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Uint64(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []float32:
		slice, ok := any(*f.value).([]float32)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Float32(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []float64:
		slice, ok := any(*f.value).([]float64)
		if !ok {
			return errBadType(*f.value)
		}

		newValue, err := parse.Float64(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, typ, err)
		}

		slice = append(slice, newValue)
		*f.value = *parse.Cast[T](&slice)

		return nil
	case []string:
		slice, ok := any(*f.value).([]string)
		if !ok {
			return errBadType(*f.value)
		}

		// No parsing to do because a string is... well, a string
		slice = append(slice, str)
		*f.value = *parse.Cast[T](&slice)

		return nil
	default:
		return fmt.Errorf("Flag.Set: unsupported flag type: %T", typ)
	}
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
	if short == flag.NoShortHand {
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

// errBadType makes a consistent error in the face of a bad type
// assertion.
func errBadType[T flag.Flaggable](value T) error {
	return fmt.Errorf("bad value %v, could not cast to %T", value, value)
}

// isZeroIsh reports whether value is the zero value (ish) for it's type.
//
// "ish" means that empty slices will return true from isZeroIsh despite their official
// zero value being nil. The primary use of isZeroIsh is to determine whether or not
// a default value is worth displaying to the user in the help text, and an empty slice
// is probably not.
func isZeroIsh[T flag.Flaggable](value T) bool { //nolint:cyclop // Not much else we can do here
	// Note: all the slice values ([]T) are in their own separate branches because if you
	// combine them, the resulting value in the body of the case block is 'any' and
	// you cannot do len(any)
	switch typ := any(value).(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64:
		return typ == 0
	case flag.Count:
		return typ == flag.Count(0)
	case string:
		return typ == ""
	case bool:
		return !typ
	case []byte:
		return len(typ) == 0
	case net.IP:
		return len(typ) == 0
	case []int:
		return len(typ) == 0
	case []int8:
		return len(typ) == 0
	case []int16:
		return len(typ) == 0
	case []int32:
		return len(typ) == 0
	case []int64:
		return len(typ) == 0
	case []uint:
		return len(typ) == 0
	case []uint16:
		return len(typ) == 0
	case []uint32:
		return len(typ) == 0
	case []uint64:
		return len(typ) == 0
	case []float32:
		return len(typ) == 0
	case []float64:
		return len(typ) == 0
	case []string:
		return len(typ) == 0
	case time.Time:
		var zero time.Time
		return typ.Equal(zero)
	case time.Duration:
		var zero time.Duration
		return typ == zero
	default:
		return false
	}
}
