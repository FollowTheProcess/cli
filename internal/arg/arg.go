// Package arg provides a command line arg definition and parsing library.
//
// Arg is intentionally internal so the only interaction is via the Arg option on a command.
package arg

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

	"go.followtheprocess.codes/cli/arg"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/cli/internal/kind"
	"go.followtheprocess.codes/cli/internal/parse"
)

var _ Value = Arg[string]{} // This will fail if we violate our Value interface

// Arg represents a single command line argument.
type Arg[T arg.Argable] struct {
	value   *T        // The actual stored value
	config  Config[T] // Additional configuration
	name    string    // Name of the argument as it appears on the command line
	usage   string    // One line description of the argument.
	typeStr string    // Cached result of Type()
	kind    kind.Kind // Cached concrete kind of T, set in New so hot paths skip any() boxing
}

// New constructs and returns a new [Arg].
func New[T arg.Argable](p *T, name, usage string, config Config[T]) (Arg[T], error) {
	if err := validateArgName(name); err != nil {
		return Arg[T]{}, fmt.Errorf("invalid arg name %q: %w", name, err)
	}

	if p == nil {
		p = new(T)
	}

	k, typeStr := typeInfo[T]()

	argument := Arg[T]{
		value:   p,
		name:    name,
		usage:   usage,
		config:  config,
		typeStr: typeStr,
		kind:    k,
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

// Default returns the default value of the argument as a string
// or "" if the argument is required.
func (a Arg[T]) Default() string {
	if a.config.DefaultValue == nil {
		// DefaultValue is nil, therefore this is a required arg
		return ""
	}

	return formatValue(a.kind, a.config.DefaultValue)
}

// String returns the string representation of the current value of the arg.
func (a Arg[T]) String() string {
	if a.value == nil {
		return format.Nil
	}

	return formatValue(a.kind, a.value)
}

// Type returns a string representation of the type of the Arg.
func (a Arg[T]) Type() string {
	if a.value == nil {
		return format.Nil
	}

	return a.typeStr
}

// Set sets an [Arg] value by parsing it's string value.
//
//nolint:gocognit,maintidx,cyclop // No other way of doing this realistically
func (a Arg[T]) Set(str string) error {
	if a.value == nil {
		return fmt.Errorf("cannot set value %s, arg.value was nil", str)
	}

	switch a.kind {
	case kind.Int:
		val, err := parse.Int(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Int8:
		val, err := parse.Int8(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Int16:
		val, err := parse.Int16(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Int32:
		val, err := parse.Int32(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Int64:
		val, err := parse.Int64(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint:
		val, err := parse.Uint(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint8:
		val, err := parse.Uint8(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint16:
		val, err := parse.Uint16(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint32:
		val, err := parse.Uint32(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Uint64, kind.Uintptr:
		val, err := parse.Uint64(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Float32:
		val, err := parse.Float32(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Float64:
		val, err := parse.Float64(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.String:
		val := str
		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.URL:
		val, err := url.ParseRequestURI(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.BytesHex:
		val, err := hex.DecodeString(strings.TrimSpace(str))
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Time:
		val, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.Duration:
		val, err := time.ParseDuration(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case kind.IP:
		val := net.ParseIP(str)
		if val == nil {
			return parse.Error(parse.KindArgument, a.name, str, *a.value, errors.New("invalid IP address"))
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	default:
		return fmt.Errorf("Arg.Set: unsupported arg type: %T", *a.value)
	}
}

// typeInfo computes the type-dependent metadata (kind, type string) for an
// arg of type T. It is called once per arg at construction so that hot paths
// (Set, String, Type) never have to type-switch on any(*a.value), which would
// box the value on every call.
//
//nolint:cyclop // No other way of doing this realistically
func typeInfo[T arg.Argable]() (kind.Kind, string) {
	var zero T

	switch typ := any(zero).(type) {
	case int:
		return kind.Int, format.TypeInt
	case int8:
		return kind.Int8, format.TypeInt8
	case int16:
		return kind.Int16, format.TypeInt16
	case int32:
		return kind.Int32, format.TypeInt32
	case int64:
		return kind.Int64, format.TypeInt64
	case uint:
		return kind.Uint, format.TypeUint
	case uint8:
		return kind.Uint8, format.TypeUint8
	case uint16:
		return kind.Uint16, format.TypeUint16
	case uint32:
		return kind.Uint32, format.TypeUint32
	case uint64:
		return kind.Uint64, format.TypeUint64
	case uintptr:
		return kind.Uintptr, format.TypeUintptr
	case float32:
		return kind.Float32, format.TypeFloat32
	case float64:
		return kind.Float64, format.TypeFloat64
	case string:
		return kind.String, format.TypeString
	case *url.URL:
		return kind.URL, format.TypeURL
	case bool:
		return kind.Bool, format.TypeBool
	case []byte:
		return kind.BytesHex, format.TypeBytesHex
	case time.Time:
		return kind.Time, format.TypeTime
	case time.Duration:
		return kind.Duration, format.TypeDuration
	case net.IP:
		return kind.IP, format.TypeIP
	default:
		return kind.Invalid, fmt.Sprintf("%T", typ)
	}
}

// formatValue renders the value pointed to by p as a string using the kind dispatch.
//
//nolint:cyclop // No other way of doing this realistically
func formatValue[T arg.Argable](k kind.Kind, p *T) string {
	switch k {
	case kind.Int:
		return format.Int(*parse.Cast[int](p))
	case kind.Int8:
		return format.Int(*parse.Cast[int8](p))
	case kind.Int16:
		return format.Int(*parse.Cast[int16](p))
	case kind.Int32:
		return format.Int(*parse.Cast[int32](p))
	case kind.Int64:
		return format.Int(*parse.Cast[int64](p))
	case kind.Uint:
		return format.Uint(*parse.Cast[uint](p))
	case kind.Uint8:
		return format.Uint(*parse.Cast[uint8](p))
	case kind.Uint16:
		return format.Uint(*parse.Cast[uint16](p))
	case kind.Uint32:
		return format.Uint(*parse.Cast[uint32](p))
	case kind.Uint64:
		return format.Uint(*parse.Cast[uint64](p))
	case kind.Uintptr:
		return format.Uint(*parse.Cast[uintptr](p))
	case kind.Float32:
		return format.Float32(*parse.Cast[float32](p))
	case kind.Float64:
		return format.Float64(*parse.Cast[float64](p))
	case kind.String:
		return *parse.Cast[string](p)
	case kind.URL:
		u := *parse.Cast[*url.URL](p)
		if u == nil {
			return format.Nil
		}

		return u.String()
	case kind.Bool:
		return strconv.FormatBool(*parse.Cast[bool](p))
	case kind.BytesHex:
		return hex.EncodeToString(*parse.Cast[[]byte](p))
	case kind.Time:
		return parse.Cast[time.Time](p).Format(time.RFC3339)
	case kind.Duration:
		return parse.Cast[time.Duration](p).String()
	case kind.IP:
		return parse.Cast[net.IP](p).String()
	default:
		return fmt.Sprintf("Arg.String: unsupported arg type: %T", *p)
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
