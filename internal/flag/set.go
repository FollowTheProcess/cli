package flag

import (
	"errors"
	"fmt"
	"strings"
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

// Add adds a flag to the Set.
func (s *Set) Add(name string, shorthand rune, usage string, flag Value) error {
	_, exists := s.flags[name]
	if exists {
		return fmt.Errorf("flag %q already defined", name)
	}
	var defaultValueNoArg string
	if flag.Type() == "bool" {
		// Boolean flags imply passing true, "--force" vs "--force true"
		defaultValueNoArg = "true"
	}
	entry := flagEntry{
		value:             flag,
		name:              name,
		usage:             usage,
		defaultValue:      flag.String(),
		defaultValueNoArg: defaultValueNoArg,
		shorthand:         shorthand,
	}
	s.flags[name] = entry

	// Only add the shorthand if it wasn't opted out of
	if shorthand != NoShortHand {
		s.shorthands[shorthand] = entry
	}

	return nil
}

// Get gets a flag Value from the Set by name and a boolean to indicate
// whether it was present.
func (s *Set) Get(name string) (Value, bool) {
	entry, ok := s.flags[name]
	if !ok {
		return nil, false
	}
	return entry.value, true
}

// Parse parses flags and their values from the command line.
func (s *Set) Parse(args []string) (err error) {
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
			return errors.New("TODO")
		default:
			// Regular positional argument
			s.args = append(s.args, arg)
		}
	}

	return nil
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
	name, equalsValue, containsEquals := strings.Cut(name, "=")
	if err := validateFlagName(name); err != nil {
		return nil, fmt.Errorf("invalid flag name %q: %w", name, err)
	}
	flag, exists := s.flags[name]
	if !exists {
		return nil, fmt.Errorf("unrecognised flag: --%s", name)
	}

	if containsEquals {
		// Must be "flag=value"
		err := flag.value.Set(equalsValue)
		if err != nil {
			return nil, fmt.Errorf("failed to set value %s for flag --%s: %w", equalsValue, name, err)
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
			return nil, fmt.Errorf("failed to set value %s for flag --%s: %w", flag.defaultValueNoArg, name, err)
		}
		// Done, as above no need to cut anything
		return rest, nil
	case len(rest) > 0:
		// --flag value
		value := rest[0]
		err := flag.value.Set(value)
		if err != nil {
			return nil, fmt.Errorf("failed to set value %s for flag --%s: %w", value, name, err)
		}
		// Done, cut value from args and return
		return rest[1:], nil
	default:
		// --flag (value was required)
		return nil, fmt.Errorf("flag --%s requires an argument", name)
	}
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
