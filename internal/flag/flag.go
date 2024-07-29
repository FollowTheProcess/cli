// Package flag provides an experimental command line flag definition and parsing library.
//
// CLI currently uses [spf13/pflag] for flag parsing (like Cobra), this package is an attempt at defining
// a new approach with some of the tools we now have in modern Go. It is not intended to be backwards compatible
// with pflag or the std lib flag package.
//
// Note: I'm using [spf13/pflag] here underneath as a gateway for now as it provides a lot of helpful functionality whilst I
// figure out what I want this to look like. So for now Flag implements pflag.Value so it can be used as a drop in.
//
// Flag is intentionally internal so the only interraction is via the Flag option on a command.
//
// [spf13/pflag]: https://github.com/spf13/pflag
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

	"github.com/spf13/pflag"
)

const (
	_      = 4 << iota // Unused
	bits8              // 8 bit integer
	bits16             // 16 bit integer
	bits32             // 32 bit integer
	bits64             // 64 bit integer
)

var _ pflag.Value = Flag[string]{} // This will fail if we violate pflag.Value.

// Flaggable is a type constraint that defines any type capable of being parsed as a command line flag.
//
// It's worth noting that the complete set of supported types is wider than this constraint appears
// as e.g. a [time.Duration] is actually just an int64 underneath, likewise a [net.IP] is actually just []byte.
type Flaggable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string | ~bool | ~[]byte | time.Time
}

// Flag represents a single command line flag.
type Flag[T Flaggable] struct {
	value *T     // The actual stored value
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
func New[T Flaggable](p *T, name string, short string, value T, usage string) (Flag[T], error) {
	if err := validateFlagName(name); err != nil {
		return Flag[T]{}, fmt.Errorf("invalid flag name: %w", err)
	}
	if err := validateFlagShort(short); err != nil {
		return Flag[T]{}, fmt.Errorf("invalid shorthand for flag %q: %w", name, err)
	}

	if p == nil {
		p = new(T)
	}
	*p = value

	flag := Flag[T]{
		value: p,
		name:  name,
		usage: usage,
		short: short,
	}

	return flag, nil
}

// Get gets a [Flag] value.
func (f Flag[T]) Get() T {
	if f.value == nil {
		return *new(T)
	}
	return *f.value
}

// String implements [fmt.Stringer] for a [Flag], and also implements the String
// part of [pflag.Value], allowing a flag to print itself.
func (f Flag[T]) String() string { //nolint:gocyclo // No other way of doing this realistically
	if f.value == nil {
		return ""
	}

	switch typ := any(f.value).(type) {
	case *int:
		return formatInt(*typ)
	case *int8:
		return formatInt(*typ)
	case *int16:
		return formatInt(*typ)
	case *int32:
		return formatInt(*typ)
	case *int64:
		return formatInt(*typ)
	case *uint:
		return formatUint(*typ)
	case *uint8:
		return formatUint(*typ)
	case *uint16:
		return formatUint(*typ)
	case *uint32:
		return formatUint(*typ)
	case *uint64:
		return formatUint(*typ)
	case *uintptr:
		return formatUint(*typ)
	case *float32:
		return formatFloat[float32](bits32)(*typ)
	case *float64:
		return formatFloat[float64](bits64)(*typ)
	case *string:
		return *typ
	case *bool:
		return strconv.FormatBool(*typ)
	case *[]byte:
		return hex.EncodeToString(*typ)
	case *time.Time:
		return typ.Format(time.RFC3339)
	case *time.Duration:
		return typ.String()
	case *net.IP:
		return typ.String()
	case fmt.Stringer:
		return typ.String()
	default:
		return ""
	}
}

// Type returns a string representation of the type of the Flag.
func (f Flag[T]) Type() string { //nolint:gocyclo // No other way of doing this realistically
	if f.value == nil {
		return ""
	}
	switch typ := any(f.value).(type) {
	case *int:
		return "int"
	case *int8:
		return "int8"
	case *int16:
		return "int16"
	case *int32:
		return "int32"
	case *int64:
		return "int64"
	case *uint:
		return "uint"
	case *uint8:
		return "uint8"
	case *uint16:
		return "uint16"
	case *uint32:
		return "uint32"
	case *uint64:
		return "uint64"
	case *uintptr:
		return "uintptr"
	case *float32:
		return "float32"
	case *float64:
		return "float64"
	case *string:
		return "string"
	case *bool:
		return "bool"
	case *[]byte:
		return "bytesHex"
	case *time.Time:
		return "time"
	case *time.Duration:
		return "duration"
	case *net.IP:
		return "ip"
	default:
		return fmt.Sprintf("%T", typ)
	}
}

// Set sets a [Flag] value based on string input, i.e. parsing from the command line.
func (f Flag[T]) Set(str string) error { //nolint:gocyclo // No other way of doing this realistically
	if f.value == nil {
		return fmt.Errorf("cannot set value %s, flag.value was nil", str)
	}
	switch typ := any(f.value).(type) {
	case *int:
		val, err := parseInt[int](0)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *int8:
		val, err := parseInt[int8](bits8)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *int16:
		val, err := parseInt[int16](bits16)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *int32:
		val, err := parseInt[int32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *int64:
		val, err := parseInt[int64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *uint:
		val, err := parseUint[uint](0)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *uint8:
		val, err := parseUint[uint8](bits8)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *uint16:
		val, err := parseUint[uint16](bits16)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *uint32:
		val, err := parseUint[uint32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *uint64:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *uintptr:
		val, err := parseUint[uint64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *float32:
		val, err := parseFloat[float32](bits32)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *float64:
		val, err := parseFloat[float64](bits64)(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *string:
		val := str
		*f.value = *cast[T](&val)
		return nil
	case *bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *[]byte:
		val, err := hex.DecodeString(strings.TrimSpace(str))
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *time.Time:
		val, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *time.Duration:
		val, err := time.ParseDuration(str)
		if err != nil {
			return errParse(f.name, str, typ, err)
		}
		*f.value = *cast[T](&val)
		return nil
	case *net.IP:
		val := net.ParseIP(str)
		if val == nil {
			return errParse(f.name, str, typ, errors.New("invalid IP address"))
		}
		*f.value = *cast[T](&val)
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

// validateFlagName ensures a flag name is valid, returning an error if it's not.
//
// Flags names must be all lower case ASCII letters, a hypen separator is allowed e.g. "set-default"
// but this must be in between letters, not leading or trailing.
func validateFlagName(name string) error {
	if name == "" {
		return errors.New("must not be empty")
	}
	before, after, found := strings.Cut(name, "-")

	// Hyphen must be in between "words" like "set-default"
	// we can't have "-default" or "default-"
	if found && after == "" {
		return fmt.Errorf("trailing hyphen: %q", name)
	}

	if found && before == "" {
		return fmt.Errorf("leading hyphen: %q", name)
	}
	for _, char := range name {
		// Only ASCII characters allowed
		if char > unicode.MaxASCII {
			return fmt.Errorf("non ascii character: %q", string(char))
		}
		// Only non-letter character allowed is a hyphen
		if !unicode.IsLetter(char) && char != '-' {
			return fmt.Errorf("not ascii letter: %q", string(char))
		}
		// Any upper case letters are not allowed
		if unicode.IsLetter(char) && !unicode.IsLower(char) {
			return fmt.Errorf("upper case character %q", string(char))
		}
	}

	return nil
}

func validateFlagShort(short string) error {
	// len(short) > 1 means an error, shorthand must be a single character
	if length := utf8.RuneCountInString(short); length > 1 {
		return fmt.Errorf("must be a single ASCII letter, got %q which has %d letters", short, length)
	}

	if short != "" {
		// Shorthand must be a valid ASCII letter
		char, _ := utf8.DecodeRuneInString(short)
		if char == utf8.RuneError || char > unicode.MaxASCII || !unicode.IsLetter(char) {
			return fmt.Errorf("invalid character, must be a single ASCII letter, got %q", string(char))
		}
	}

	return nil
}
