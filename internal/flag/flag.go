// Package flag provides a command line flag definition and parsing library.
//
// Flag is intentionally internal so the only interaction is via the Flag option on a command.
package flag

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/cli/internal/kind"
	"go.followtheprocess.codes/cli/internal/parse"
)

var _ Value = &Flag[string]{} // This will fail if we violate our Value interface

// Flag represents a single command line flag.
type Flag[T flag.Flaggable] struct {
	value      *T        // The actual stored value
	name       string    // The name of the flag as appears on the command line, e.g. "force" for a --force flag
	usage      string    // one line description of the flag, e.g. "Force deletion without confirmation"
	envVar     string    // Name of an environment variable that may set this flag's value if the flag is not explicitly provided on the command line
	typeStr    string    // Cached result of Type()
	noArgValue string    // Cached result of NoArgValue()
	short      rune      // Optional shorthand version of the flag, e.g. "f" for a -f flag
	kind       kind.Kind // Cached concrete kind of T
	isSlice    bool      // Cached result of IsSlice()
}

// New constructs and returns a new [Flag].
//
// The name should be as it appears on the command line, e.g. "force" for a --force flag. An optional
// shorthand can be created by setting short to a single letter value, e.g. "f" to also create a -f version of "force".
func New[T flag.Flaggable](p *T, name string, short rune, usage string, config Config[T]) (*Flag[T], error) {
	if err := validateFlagName(name); err != nil {
		return nil, fmt.Errorf("invalid flag name %q: %w", name, err)
	}

	if err := validateFlagShort(short); err != nil {
		return nil, fmt.Errorf("invalid shorthand for flag %q: %w", name, err)
	}

	if p == nil {
		return nil, fmt.Errorf("flag %q: target pointer must not be nil", name)
	}

	*p = config.DefaultValue

	info := typeInfo[T]()

	return &Flag[T]{
		value:      p,
		name:       name,
		usage:      usage,
		short:      short,
		envVar:     config.EnvVar,
		typeStr:    info.typeStr,
		noArgValue: info.noArgValue,
		kind:       info.kind,
		isSlice:    info.isSlice,
	}, nil
}

// Name returns the name of the [Flag].
func (f *Flag[T]) Name() string {
	return f.name
}

// Short returns the shorthand registered for the flag (e.g. -d for --delete), or
// NoShortHand if the flag should be long only.
func (f *Flag[T]) Short() rune {
	return f.short
}

// Usage returns the usage line for the flag.
func (f *Flag[T]) Usage() string {
	return f.usage
}

// Default returns the default value for the flag, as a string.
//
// If the flag's default is unset (i.e. the zero value for its type),
// an empty string is returned.
func (f *Flag[T]) Default() string {
	// Special case a --help flag, because if we didn't, when you call --help
	// it would show up with a default of true because you've passed it
	// so it's value is true here
	if f.isZeroIsh() || f.name == "help" {
		return ""
	}

	return f.String()
}

// EnvVar returns the name of the environment variable associated with this flag,
// or an empty string if none was configured.
func (f *Flag[T]) EnvVar() string {
	return f.envVar
}

// IsSlice reports whether the flag holds a slice value that accumulates repeated
// calls to Set. Returns false for []byte and net.IP, which are parsed atomically.
func (f *Flag[T]) IsSlice() bool {
	return f.isSlice
}

// NoArgValue returns a string representation of value the flag should hold
// when it is given no arguments on the command line. For example a boolean flag
// --delete, when passed without arguments implies --delete true.
func (f *Flag[T]) NoArgValue() string {
	return f.noArgValue
}

// Type returns a string representation of the type of the Flag.
func (f *Flag[T]) Type() string {
	if f.value == nil {
		// Nil-safe behaviour for zero-value composite-literal flags.
		return format.Nil
	}

	return f.typeStr
}

// String implements [fmt.Stringer] for a [Flag], and also implements the String
// part of [Value], allowing a flag to print itself.
//
//nolint:cyclop // No other way of doing this realistically
func (f *Flag[T]) String() string {
	if f.value == nil {
		return format.Nil
	}

	switch f.kind {
	case kind.Int:
		return format.Int(*parse.Cast[int](f.value))
	case kind.Int8:
		return format.Int(*parse.Cast[int8](f.value))
	case kind.Int16:
		return format.Int(*parse.Cast[int16](f.value))
	case kind.Int32:
		return format.Int(*parse.Cast[int32](f.value))
	case kind.Int64:
		return format.Int(*parse.Cast[int64](f.value))
	case kind.Count:
		return format.Uint(*parse.Cast[flag.Count](f.value))
	case kind.Uint:
		return format.Uint(*parse.Cast[uint](f.value))
	case kind.Uint8:
		return format.Uint(*parse.Cast[uint8](f.value))
	case kind.Uint16:
		return format.Uint(*parse.Cast[uint16](f.value))
	case kind.Uint32:
		return format.Uint(*parse.Cast[uint32](f.value))
	case kind.Uint64:
		return format.Uint(*parse.Cast[uint64](f.value))
	case kind.Uintptr:
		return format.Uint(*parse.Cast[uintptr](f.value))
	case kind.Float32:
		return format.Float32(*parse.Cast[float32](f.value))
	case kind.Float64:
		return format.Float64(*parse.Cast[float64](f.value))
	case kind.String:
		return *parse.Cast[string](f.value)
	case kind.Bool:
		return strconv.FormatBool(*parse.Cast[bool](f.value))
	case kind.BytesHex:
		return hex.EncodeToString(*parse.Cast[[]byte](f.value))
	case kind.Time:
		return parse.Cast[time.Time](f.value).Format(time.RFC3339)
	case kind.Duration:
		return parse.Cast[time.Duration](f.value).String()
	case kind.IP:
		return parse.Cast[net.IP](f.value).String()
	case kind.URL:
		u := *parse.Cast[*url.URL](f.value)
		if u == nil {
			return format.Nil
		}

		return u.String()
	case kind.IntSlice:
		return format.Slice(*parse.Cast[[]int](f.value))
	case kind.Int8Slice:
		return format.Slice(*parse.Cast[[]int8](f.value))
	case kind.Int16Slice:
		return format.Slice(*parse.Cast[[]int16](f.value))
	case kind.Int32Slice:
		return format.Slice(*parse.Cast[[]int32](f.value))
	case kind.Int64Slice:
		return format.Slice(*parse.Cast[[]int64](f.value))
	case kind.UintSlice:
		return format.Slice(*parse.Cast[[]uint](f.value))
	case kind.Uint16Slice:
		return format.Slice(*parse.Cast[[]uint16](f.value))
	case kind.Uint32Slice:
		return format.Slice(*parse.Cast[[]uint32](f.value))
	case kind.Uint64Slice:
		return format.Slice(*parse.Cast[[]uint64](f.value))
	case kind.Float32Slice:
		return format.Slice(*parse.Cast[[]float32](f.value))
	case kind.Float64Slice:
		return format.Slice(*parse.Cast[[]float64](f.value))
	case kind.StringSlice:
		return format.Slice(*parse.Cast[[]string](f.value))
	default:
		return fmt.Sprintf("Flag.String: unsupported flag type: %T", *f.value)
	}
}

// info bundles the cacheable, type-dependent metadata for a Flag of a given T.
type info struct {
	typeStr    string
	noArgValue string
	kind       kind.Kind
	isSlice    bool
}

// typeInfo computes the type-dependent metadata (kind, type string, no-arg
// value, isSlice) for a flag of type T. It is called once per flag at
// construction so that the hot path of Parse never has to type-switch on
// any(*f.value), which would otherwise box the value on every call.
func typeInfo[T flag.Flaggable]() info { //nolint:cyclop // No other way of doing this realistically
	var zero T

	switch typ := any(zero).(type) {
	case int:
		return info{kind: kind.Int, typeStr: format.TypeInt}
	case int8:
		return info{kind: kind.Int8, typeStr: format.TypeInt8}
	case int16:
		return info{kind: kind.Int16, typeStr: format.TypeInt16}
	case int32:
		return info{kind: kind.Int32, typeStr: format.TypeInt32}
	case int64:
		return info{kind: kind.Int64, typeStr: format.TypeInt64}
	case flag.Count:
		return info{kind: kind.Count, typeStr: format.TypeCount, noArgValue: "1"}
	case uint:
		return info{kind: kind.Uint, typeStr: format.TypeUint}
	case uint8:
		return info{kind: kind.Uint8, typeStr: format.TypeUint8}
	case uint16:
		return info{kind: kind.Uint16, typeStr: format.TypeUint16}
	case uint32:
		return info{kind: kind.Uint32, typeStr: format.TypeUint32}
	case uint64:
		return info{kind: kind.Uint64, typeStr: format.TypeUint64}
	case uintptr:
		return info{kind: kind.Uintptr, typeStr: format.TypeUintptr}
	case float32:
		return info{kind: kind.Float32, typeStr: format.TypeFloat32}
	case float64:
		return info{kind: kind.Float64, typeStr: format.TypeFloat64}
	case string:
		return info{kind: kind.String, typeStr: format.TypeString}
	case bool:
		return info{kind: kind.Bool, typeStr: format.TypeBool, noArgValue: format.True}
	case []byte:
		return info{kind: kind.BytesHex, typeStr: format.TypeBytesHex}
	case time.Time:
		return info{kind: kind.Time, typeStr: format.TypeTime}
	case time.Duration:
		return info{kind: kind.Duration, typeStr: format.TypeDuration}
	case net.IP:
		return info{kind: kind.IP, typeStr: format.TypeIP}
	case *url.URL:
		return info{kind: kind.URL, typeStr: format.TypeURL}
	case []int:
		return info{kind: kind.IntSlice, typeStr: format.TypeIntSlice, isSlice: true}
	case []int8:
		return info{kind: kind.Int8Slice, typeStr: format.TypeInt8Slice, isSlice: true}
	case []int16:
		return info{kind: kind.Int16Slice, typeStr: format.TypeInt16Slice, isSlice: true}
	case []int32:
		return info{kind: kind.Int32Slice, typeStr: format.TypeInt32Slice, isSlice: true}
	case []int64:
		return info{kind: kind.Int64Slice, typeStr: format.TypeInt64Slice, isSlice: true}
	case []uint:
		return info{kind: kind.UintSlice, typeStr: format.TypeUintSlice, isSlice: true}
	case []uint16:
		return info{kind: kind.Uint16Slice, typeStr: format.TypeUint16Slice, isSlice: true}
	case []uint32:
		return info{kind: kind.Uint32Slice, typeStr: format.TypeUint32Slice, isSlice: true}
	case []uint64:
		return info{kind: kind.Uint64Slice, typeStr: format.TypeUint64Slice, isSlice: true}
	case []float32:
		return info{kind: kind.Float32Slice, typeStr: format.TypeFloat32Slice, isSlice: true}
	case []float64:
		return info{kind: kind.Float64Slice, typeStr: format.TypeFloat64Slice, isSlice: true}
	case []string:
		return info{kind: kind.StringSlice, typeStr: format.TypeStringSlice, isSlice: true}
	default:
		return info{kind: kind.Invalid, typeStr: fmt.Sprintf("%T", typ)}
	}
}

// Set sets a [Flag] value based on string input, i.e. parsing from the command line.
//
//nolint:gocognit,maintidx,cyclop // No other way of doing this realistically
func (f *Flag[T]) Set(str string) error {
	if f.value == nil {
		return fmt.Errorf("cannot set value %s, flag.value was nil", str)
	}

	switch f.kind {
	case kind.Int:
		val, err := parse.Int(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Int8:
		val, err := parse.Int8(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Int16:
		val, err := parse.Int16(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Int32:
		val, err := parse.Int32(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Int64:
		val, err := parse.Int64(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Count:
		// Add the count and store it back, we still parse the given str rather
		// than just +1 every time as this allows people to do e.g. --verbosity=3
		// as well as -vvv
		val, err := parse.Uint(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		newValue := *parse.Cast[flag.Count](f.value) + flag.Count(val)
		*f.value = *parse.Cast[T](&newValue)

		return nil
	case kind.Uint:
		val, err := parse.Uint(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint8:
		val, err := parse.Uint8(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint16:
		val, err := parse.Uint16(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint32:
		val, err := parse.Uint32(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint64, kind.Uintptr:
		val, err := parse.Uint64(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Float32:
		val, err := parse.Float32(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Float64:
		val, err := parse.Float64(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.String:
		val := str
		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.BytesHex:
		val, err := hex.DecodeString(strings.TrimSpace(str))
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Time:
		val, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.Duration:
		val, err := time.ParseDuration(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.IP:
		val := net.ParseIP(str)
		if val == nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, errors.New("invalid IP address"))
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.URL:
		val, err := url.ParseRequestURI(str)
		if err != nil {
			return parse.Error(parse.KindFlag, f.name, str, *f.value, err)
		}

		*f.value = *parse.Cast[T](&val)

		return nil
	case kind.IntSlice:
		// Like Count, a slice flag is a read/write op
		newValue, err := parse.Int(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]int](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Int8Slice:
		newValue, err := parse.Int8(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]int8](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Int16Slice:
		newValue, err := parse.Int16(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]int16](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Int32Slice:
		newValue, err := parse.Int32(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]int32](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Int64Slice:
		newValue, err := parse.Int64(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]int64](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.UintSlice:
		newValue, err := parse.Uint(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]uint](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Uint16Slice:
		newValue, err := parse.Uint16(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]uint16](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Uint32Slice:
		newValue, err := parse.Uint32(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]uint32](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Uint64Slice:
		newValue, err := parse.Uint64(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]uint64](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Float32Slice:
		newValue, err := parse.Float32(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]float32](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.Float64Slice:
		newValue, err := parse.Float64(str)
		if err != nil {
			return parse.ErrorSlice(parse.KindFlag, f.name, str, *f.value, err)
		}

		typ := append(*parse.Cast[[]float64](f.value), newValue)
		*f.value = *parse.Cast[T](&typ)

		return nil
	case kind.StringSlice:
		typ := *parse.Cast[[]string](f.value)
		typ = append(typ, str)
		*f.value = *parse.Cast[T](&typ)

		return nil
	default:
		return fmt.Errorf("Flag.Set: unsupported flag type: %T", *f.value)
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

// isZeroIsh reports whether the flag's value is the zero value (ish) for it's type.
//
// "ish" means that empty slices will return true despite their official zero
// value being nil. The primary use is to determine whether a default value is
// worth displaying to the user in the help text, an empty slice is probably
// not.
//
//nolint:cyclop // Not much else we can do here
func (f *Flag[T]) isZeroIsh() bool {
	switch f.kind {
	case kind.Int:
		return *parse.Cast[int](f.value) == 0
	case kind.Int8:
		return *parse.Cast[int8](f.value) == 0
	case kind.Int16:
		return *parse.Cast[int16](f.value) == 0
	case kind.Int32:
		return *parse.Cast[int32](f.value) == 0
	case kind.Int64:
		return *parse.Cast[int64](f.value) == 0
	case kind.Uint:
		return *parse.Cast[uint](f.value) == 0
	case kind.Uint8:
		return *parse.Cast[uint8](f.value) == 0
	case kind.Uint16:
		return *parse.Cast[uint16](f.value) == 0
	case kind.Uint32:
		return *parse.Cast[uint32](f.value) == 0
	case kind.Uint64:
		return *parse.Cast[uint64](f.value) == 0
	case kind.Uintptr:
		return *parse.Cast[uintptr](f.value) == 0
	case kind.Float32:
		return *parse.Cast[float32](f.value) == 0
	case kind.Float64:
		return *parse.Cast[float64](f.value) == 0
	case kind.Count:
		return *parse.Cast[flag.Count](f.value) == 0
	case kind.String:
		return *parse.Cast[string](f.value) == ""
	case kind.Bool:
		return !*parse.Cast[bool](f.value)
	case kind.BytesHex:
		return len(*parse.Cast[[]byte](f.value)) == 0
	case kind.IP:
		return len(*parse.Cast[net.IP](f.value)) == 0
	case kind.URL:
		return *parse.Cast[*url.URL](f.value) == nil
	case kind.IntSlice:
		return len(*parse.Cast[[]int](f.value)) == 0
	case kind.Int8Slice:
		return len(*parse.Cast[[]int8](f.value)) == 0
	case kind.Int16Slice:
		return len(*parse.Cast[[]int16](f.value)) == 0
	case kind.Int32Slice:
		return len(*parse.Cast[[]int32](f.value)) == 0
	case kind.Int64Slice:
		return len(*parse.Cast[[]int64](f.value)) == 0
	case kind.UintSlice:
		return len(*parse.Cast[[]uint](f.value)) == 0
	case kind.Uint16Slice:
		return len(*parse.Cast[[]uint16](f.value)) == 0
	case kind.Uint32Slice:
		return len(*parse.Cast[[]uint32](f.value)) == 0
	case kind.Uint64Slice:
		return len(*parse.Cast[[]uint64](f.value)) == 0
	case kind.Float32Slice:
		return len(*parse.Cast[[]float32](f.value)) == 0
	case kind.Float64Slice:
		return len(*parse.Cast[[]float64](f.value)) == 0
	case kind.StringSlice:
		return len(*parse.Cast[[]string](f.value)) == 0
	case kind.Time:
		var zero time.Time
		return parse.Cast[time.Time](f.value).Equal(zero)
	case kind.Duration:
		return *parse.Cast[time.Duration](f.value) == 0
	default:
		return false
	}
}
