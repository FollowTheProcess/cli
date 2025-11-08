// Package arg provides a command line arg definition and parsing library.
//
// Arg is intentionally internal so the only interaction is via the Arg option on a command.
package arg

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"go.followtheprocess.codes/cli/arg"
	"go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/cli/internal/constraints"
)

// TODO(@FollowTheProcess): LOTS of duplicated stuff with internal/flag.
// Once we know this is the direction to go down, then we should combine all the shared
// stuff and use it from each package

const (
	_      = 4 << iota // Unused
	bits8              // 8 bit integer
	bits16             // 16 bit integer
	bits32             // 32 bit integer
	bits64             // 64 bit integer
)

const (
	typeInt      = "int"
	typeInt8     = "int8"
	typeInt16    = "int16"
	typeInt32    = "int32"
	typeInt64    = "int64"
	typeUint     = "uint"
	typeUint8    = "uint8"
	typeUint16   = "uint16"
	typeUint32   = "uint32"
	typeUint64   = "uint64"
	typeUintptr  = "uintptr"
	typeFloat32  = "float32"
	typeFloat64  = "float64"
	typeString   = "string"
	typeBool     = "bool"
	typeBytesHex = "bytesHex"
	typeTime     = "time"
	typeDuration = "duration"
	typeIP       = "ip"
)

var _ Value = Arg[string]{} // This will fail if we violate our Value interface

// Arg represents a single command line argument.
type Arg[T arg.Argable] struct {
	value *T     // The actual stored value
	name  string // Name of the argument as it appears on the command line
	usage string // One line description of the argument.
}

// New constructs and returns a new [Arg].
func New[T arg.Argable](p *T, name, usage string) (Arg[T], error) {
	if err := validateArgName(name); err != nil {
		return Arg[T]{}, fmt.Errorf("invalid arg name %q: %w", name, err)
	}

	if p == nil {
		p = new(T)
	}

	argument := Arg[T]{
		value: p,
		name:  name,
		usage: usage,
	}

	return argument, nil
}

// Name returns the name of the Arg.
func (a Arg[T]) Name() string {
	return a.name
}

// Usage returns the usage line of the Arg.
func (a Arg[T]) Usage() string {
	return a.usage
}

// String returns the string representation of the current value of the arg.
//
//nolint:cyclop // No other way of doing this realistically
func (a Arg[T]) String() string {
	if a.value == nil {
		return "<nil>"
	}

	switch typ := any(*a.value).(type) {
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
	default:
		return fmt.Sprintf("Arg.String: unsupported arg type: %T", typ)
	}
}

// Type returns a string representation of the type of the Arg.
//
//nolint:cyclop // No other way of doing this realistically
func (a Arg[T]) Type() string {
	if a.value == nil {
		return "<nil>"
	}

	switch typ := any(*a.value).(type) {
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
	default:
		return fmt.Sprintf("%T", typ)
	}
}

// Set sets an [Arg] value by parsing it's string value.
//
//nolint:cyclop // No other way of doing this realistically
func (a Arg[T]) Set(str string) error {
	if a.value == nil {
		return fmt.Errorf("cannot set value %s, arg.value was nil", str)
	}

	switch typ := any(*a.value).(type) {
	case int:
		val, err := parseInt[int](0)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case int8:
		val, err := parseInt[int8](bits8)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case int16:
		val, err := parseInt[int16](bits16)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case int32:
		val, err := parseInt[int32](bits32)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case int64:
		val, err := parseInt[int64](bits64)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case flag.Count:
		// We have to do a bit of custom stuff here as an increment is a read and write op
		// First read the current value of the flag and cast it to a Count so we
		// can increment it
		current, ok := any(*a.value).(flag.Count)
		if !ok {
			// This basically shouldn't ever happen but it's easy enough to handle nicely
			return errBadType(*a.value)
		}

		// Add the count and store it back, we still parse the given str rather
		// than just +1 every time as this allows people to do e.g. --verbosity=3
		// as well as -vvv
		val, err := parseUint[uint](0)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		newValue := current + flag.Count(val)
		*a.value = *cast[T](&newValue)

		return nil
	case uint:
		val, err := parseUint[uint](0)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case uint8:
		val, err := parseUint[uint8](bits8)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case uint16:
		val, err := parseUint[uint16](bits16)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case uint32:
		val, err := parseUint[uint32](bits32)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case uint64:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case uintptr:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case float32:
		val, err := parseFloat[float32](bits32)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case float64:
		val, err := parseFloat[float64](bits64)(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case string:
		val := str
		*a.value = *cast[T](&val)

		return nil
	case bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case []byte:
		val, err := hex.DecodeString(strings.TrimSpace(str))
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case time.Time:
		val, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case time.Duration:
		val, err := time.ParseDuration(str)
		if err != nil {
			return errParse(a.name, str, typ, err)
		}

		*a.value = *cast[T](&val)

		return nil
	case net.IP:
		val := net.ParseIP(str)
		if val == nil {
			return errParse(a.name, str, typ, errors.New("invalid IP address"))
		}

		*a.value = *cast[T](&val)

		return nil
	default:
		return fmt.Errorf("Arg.Set: unsupported arg type: %T", typ)
	}
}

// validateArgName ensures an argument name is valid, returning an error if it's not.
//
// Arg names must be all lower case ASCII letters, a hyphen separator is allowed e.g. "workspace-dir"
// but this must be in between letters, not leading or trailing.
func validateArgName(name string) error {
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

// formatInt is a generic helper to return a string representation of any signed integer.
func formatInt[T constraints.Signed](in T) string {
	return strconv.FormatInt(int64(in), 10)
}

// formatUint is a generic helper to return a string representation of any unsigned integer.
func formatUint[T constraints.Unsigned](in T) string {
	return strconv.FormatUint(uint64(in), 10)
}

// formatFloat is a generic helper to return a string representation of any floating point digit.
func formatFloat[T ~float32 | ~float64](bits int) func(T) string {
	return func(in T) string {
		return strconv.FormatFloat(float64(in), 'g', -1, bits)
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
func cast[T2, T1 any](v *T1) *T2 {
	return (*T2)(unsafe.Pointer(v))
}

// errParse is a helper to quickly return a consistent error in the face of flag
// value parsing errors.
func errParse[T flag.Flaggable](name, str string, typ T, err error) error {
	return fmt.Errorf(
		"flag %q received invalid value %q (expected %T), detail: %w",
		name,
		str,
		typ,
		err,
	)
}

// errBadType makes a consistent error in the face of a bad type
// assertion.
func errBadType[T flag.Flaggable](value T) error {
	return fmt.Errorf("bad value %v, could not cast to %T", value, value)
}

// parseInt is a generic helper to parse all signed integers, given a bit size.
//
// It returns the parsed value or an error.
func parseInt[T constraints.Signed](bits int) func(str string) (T, error) {
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
func parseUint[T constraints.Unsigned](bits int) func(str string) (T, error) {
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
