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

	"go.followtheprocess.codes/cli/arg"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/cli/internal/parse"
)

// TODO(@FollowTheProcess): LOTS of duplicated stuff with internal/flag.
// Once we know this is the direction to go down, then we should combine all the shared
// stuff and use it from each package

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
	value  *T        // The actual stored value
	config Config[T] // Additional configuration
	name   string    // Name of the argument as it appears on the command line
	usage  string    // One line description of the argument.
}

// New constructs and returns a new [Arg].
func New[T arg.Argable](p *T, name, usage string, config Config[T]) (Arg[T], error) {
	if err := validateArgName(name); err != nil {
		return Arg[T]{}, fmt.Errorf("invalid arg name %q: %w", name, err)
	}

	if p == nil {
		p = new(T)
	}

	argument := Arg[T]{
		value:  p,
		name:   name,
		usage:  usage,
		config: config,
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
//
//nolint:cyclop // No other way of doing this
func (a Arg[T]) Default() string {
	if a.config.DefaultValue == nil {
		// DefaultValue is nil, therefore this is a required arg
		return ""
	}

	switch typ := any(*a.config.DefaultValue).(type) {
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
	default:
		return fmt.Sprintf("Arg.String: unsupported arg type: %T", typ)
	}
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
		return format.Int(typ)
	case int8:
		return format.Int(typ)
	case int16:
		return format.Int(typ)
	case int32:
		return format.Int(typ)
	case int64:
		return format.Int(typ)
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
		val, err := parse.Int(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case int8:
		val, err := parse.Int8(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case int16:
		val, err := parse.Int16(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case int32:
		val, err := parse.Int32(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case int64:
		val, err := parse.Int64(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case uint:
		val, err := parse.Uint(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case uint8:
		val, err := parse.Uint8(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case uint16:
		val, err := parse.Uint16(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case uint32:
		val, err := parse.Uint32(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case uint64:
		val, err := parse.Uint64(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case uintptr:
		val, err := parse.Uint64(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case float32:
		val, err := parse.Float32(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case float64:
		val, err := parse.Float64(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case string:
		val := str
		*a.value = *parse.Cast[T](&val)

		return nil
	case bool:
		val, err := strconv.ParseBool(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case []byte:
		val, err := hex.DecodeString(strings.TrimSpace(str))
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case time.Time:
		val, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case time.Duration:
		val, err := time.ParseDuration(str)
		if err != nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, err)
		}

		*a.value = *parse.Cast[T](&val)

		return nil
	case net.IP:
		val := net.ParseIP(str)
		if val == nil {
			return parse.Error(parse.KindArgument, a.name, str, typ, errors.New("invalid IP address"))
		}

		*a.value = *parse.Cast[T](&val)

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
