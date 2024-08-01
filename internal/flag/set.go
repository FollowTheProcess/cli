package flag

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Set is a set of command line flags.
type Set struct {
	flags      map[string]flagEntry // The actual stored flags, can lookup by name
	shorthands map[rune]flagEntry   // The flags by shorthand
	args       []string             // Arguments minus flags or flag values
}

// NewSet builds and returns a new set of flags.
func NewSet() *Set {
	return &Set{
		flags:      make(map[string]flagEntry),
		shorthands: make(map[rune]flagEntry),
	}
}

// AddToSet adds a flag to the given Set.
func AddToSet[T Flaggable](set *Set, flag Flag[T]) error {
	if set == nil {
		return errors.New("cannot add flag to a nil set")
	}
	// TODO: Would this be better as a method on Flag[T]?
	_, exists := set.flags[flag.name]
	if exists {
		return fmt.Errorf("flag %q already defined", flag.name)
	}
	var defaultValueNoArg string
	if flag.Type() == "bool" {
		// Boolean flags imply passing true, "--force" vs "--force true"
		defaultValueNoArg = "true"
	}
	entry := flagEntry{
		value:             flag,
		name:              flag.name,
		usage:             flag.usage,
		defaultValue:      flag.String(),
		defaultValueNoArg: defaultValueNoArg,
		shorthand:         flag.short,
	}
	set.flags[flag.name] = entry

	// Only add the shorthand if it wasn't opted out of
	if flag.short != NoShortHand {
		set.shorthands[flag.short] = entry
	}

	return nil
}

// Get gets a flag Value from the Set by name and a boolean to indicate
// whether it was present.
func (s *Set) Get(name string) (Value, bool) {
	if s == nil {
		return nil, false
	}
	entry, ok := s.flags[name]
	if !ok {
		return nil, false
	}
	return entry.value, true
}

// Parse parses flags and their values from the command line.
func (s *Set) Parse(args []string) (err error) {
	if s == nil {
		return errors.New("Parse called on a nil set")
	}
	for len(args) > 0 {
		arg := args[0]  // The argument we're currently inspecting
		args = args[1:] // Remainder

		switch {
		case strings.HasPrefix(arg, "--"):
			// This is a long flag e.g. --delete
			args, err = s.parseLongFlag(arg, args)
			if err != nil {
				return err
			}
		case strings.HasPrefix(arg, "-"):
			// Short flag e.g. -d
			args, err = s.parseShortFlag(arg, args)
			if err != nil {
				return err
			}
		default:
			// Regular positional argument
			s.args = append(s.args, arg)
		}
	}

	return nil
}

// flagEntry represents a single flag in the set.
type flagEntry struct {
	value             Value  // The actual Flag[T]
	name              string // The full name of the flag e.g. "delete"
	usage             string // The flag's usage message
	defaultValue      string // String representation of the default flag value
	defaultValueNoArg string // String representation of the default flag value if used without an arg, e.g. boolean flags "--force" implies "--force true"
	shorthand         rune   // The optional shorthand e.g. 'd' or [NoShortHand]
}

// parseLongFlag parses a single long flag e.g. --delete. It is passed
// the possible long flag and the rest of the argument list and returns
// the remaining arguments after it's done parsing to the caller.
//
// The forms it expects are --flag (boolean flags), --flag=value or --flag value.
func (s *Set) parseLongFlag(long string, rest []string) (remaining []string, err error) {
	// Could either be "flag" or "flag=value"
	name := strings.TrimPrefix(long, "--")

	// name will either be the entire string or the name before the "="
	name, value, containsEquals := strings.Cut(name, "=")
	if err := validateFlagName(name); err != nil {
		return nil, fmt.Errorf("invalid flag name %q: %w", name, err)
	}
	flag, exists := s.flags[name]
	if !exists {
		return nil, fmt.Errorf("unrecognised flag: --%s", name)
	}

	if containsEquals {
		// Must be "flag=value"
		err := flag.value.Set(value)
		if err != nil {
			return nil, err
		}

		// We're done, no need to cut anything from rest as this was a single arg
		return rest, nil
	}

	// Must now either be --flag (boolean) or --flag value
	switch {
	case flag.defaultValueNoArg != "":
		// --flag (boolean)
		err := flag.value.Set(flag.defaultValueNoArg)
		if err != nil {
			return nil, err
		}
		// Done, as above no need to cut anything
		return rest, nil
	case len(rest) > 0:
		// --flag value
		value := rest[0]
		err := flag.value.Set(value)
		if err != nil {
			return nil, err
		}
		// Done, cut value from args and return
		return rest[1:], nil
	default:
		// --flag (value was required)
		// Do we ever hit this? Idk when a flag would be required
		return nil, fmt.Errorf("flag --%s requires an argument", name)
	}
}

// parseShortFlag parses short flags from the command line. It is passed the possible
// short flag and the rest of the argument list and returns the remaining arguments
// after it's done parsing to the caller.
//
// The forms it expects are "-f", "-vfg", "-f value" and "-f=value".
func (s *Set) parseShortFlag(short string, rest []string) (remaining []string, err error) {
	// TODO: Refactor this to clean it up and reduce duplication

	// Could either be "f", "vfg" or "f=value"
	shorthands := strings.TrimPrefix(short, "-")

	// Is it e.g. f=value
	shorthand, value, containsEquals := strings.Cut(shorthands, "=")
	if err := validateFlagName(shorthand); err != nil {
		return nil, fmt.Errorf("invalid flag name %q: %w", shorthand, err)
	}
	if containsEquals {
		// Yes, it is. If the thing on the left of the equals is > 1 char it's an error
		if len(shorthand) != 1 {
			return nil, fmt.Errorf("invalid shorthand syntax: expected e.g. -f=<value> got %s", short)
		}

		char, _ := utf8.DecodeRuneInString(shorthand)
		if err := validateFlagShort(char); err != nil {
			return nil, fmt.Errorf("invalid flag shorthand %q: %w", string(char), err)
		}

		flag, exists := s.shorthands[char]
		if !exists {
			return nil, fmt.Errorf("unrecognised shorthand flag: -%s", string(char))
		}

		if err := flag.value.Set(value); err != nil {
			return nil, err
		}

		// We're done, nothing to trim off
		return rest, nil
	}

	// It's not "f=value" so must be one of "fvalue", "f value", or "vvv"
	// len("fvalue") is > 1 but len("f") isn't (value in that last case is the first arg in 'rest')
	if len(shorthands) > 1 {
		// It must be "fvalue", so extract "value"
		char, _ := utf8.DecodeRuneInString(shorthands)
		if err := validateFlagShort(char); err != nil {
			return nil, fmt.Errorf("invalid flag shorthand %q: %w", string(char), err)
		}
		value = shorthands[1:]
		flag, exists := s.shorthands[char]
		if !exists {
			return nil, fmt.Errorf("unrecognised shorthand flag: -%s", string(char))
		}
		if err := flag.value.Set(value); err != nil {
			return nil, err
		}

		// We're done, nothing to trim off
		return rest, nil
	}

	// Any arguments after the short flag?
	if len(rest) > 0 {
		// It must be "f value" and value is the next argument in rest
		char, _ := utf8.DecodeRuneInString(shorthands)
		if err := validateFlagShort(char); err != nil {
			return nil, fmt.Errorf("invalid flag shorthand %q: %w", string(char), err)
		}
		value = rest[0]
		flag, exists := s.shorthands[char]
		if !exists {
			return nil, fmt.Errorf("unrecognised shorthand flag: -%s", string(char))
		}
		if err := flag.value.Set(value); err != nil {
			return nil, err
		}

		// We've consumed "value" from rest so trim it off
		return rest[1:], nil
	}

	// If we get here, it must be the "-vvv" form
	for _, char := range shorthands {
		flag, exists := s.shorthands[char]
		if !exists {
			return nil, fmt.Errorf("unrecognised shorthand flag: %q in -%s", string(char), shorthands)
		}

		// -f (boolean flag)
		if flag.defaultValueNoArg != "" {
			err := flag.value.Set(flag.defaultValueNoArg)
			if err != nil {
				return nil, err
			}
			// Done, as above no need to cut anything
			return rest, nil
		}
	}

	return rest, nil
}
