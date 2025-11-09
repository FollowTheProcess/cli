package flag

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"

	"go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/cli/internal/style"
	"go.followtheprocess.codes/hue/tabwriter"
)

// usageBufferSize is sufficient to hold most commands flag usage text.
const usageBufferSize = 256

// Set is a set of command line flags.
type Set struct {
	flags      map[string]Value // The actual stored flags, can lookup by name
	shorthands map[rune]Value   // The flags by shorthand
	args       []string         // Arguments minus flags or flag values
	extra      []string         // Arguments after "--" was hit
}

// NewSet builds and returns a new set of flags.
func NewSet() *Set {
	return &Set{
		flags:      make(map[string]Value),
		shorthands: make(map[rune]Value),
	}
}

// AddToSet adds a flag to the given Set.
func AddToSet[T flag.Flaggable](set *Set, f Flag[T]) error {
	if set == nil {
		return errors.New("cannot add flag to a nil set")
	}

	name := f.Name()
	short := f.Short()

	_, exists := set.flags[name]
	if exists {
		return fmt.Errorf("flag %q already defined", name)
	}

	if short != flag.NoShortHand {
		existingFlag, exists := set.shorthands[short]
		if exists {
			return fmt.Errorf("shorthand %q already in use for flag %q", string(short), existingFlag.Name())
		}
	}

	set.flags[name] = f

	// Only add the shorthand if it wasn't opted out of
	if short != flag.NoShortHand {
		set.shorthands[short] = f
	}

	return nil
}

// Get gets a flag from the Set by name and a boolean to indicate
// whether it was present.
func (s *Set) Get(name string) (Value, bool) {
	if s == nil {
		return nil, false
	}

	flag, ok := s.flags[name]
	if !ok {
		return nil, false
	}

	return flag, true
}

// GetShort gets a flag from the Set by it's shorthand and a boolean to indicate
// whether it was present.
func (s *Set) GetShort(short rune) (Value, bool) {
	if s == nil {
		return nil, false
	}

	flag, ok := s.shorthands[short]
	if !ok {
		return nil, false
	}

	return flag, true
}

// Help returns whether the [Set] has a boolean flag named "help" and what the value
// of that flag is currently set to, it simplifies checking for --help.
func (s *Set) Help() (value, ok bool) {
	flag, exists := s.Get("help")
	if !exists {
		// No help defined
		return false, false
	}
	// Is it a bool flag?
	if flag.Type() != format.TypeBool {
		return false, false
	}
	// It is there, we can infer from the string value what it's set to
	// avoid unnecessary type conversions
	return flag.String() == "true", true
}

// Version returns whether the [Set] has a boolean flag named "version" and what the value
// of that flag is currently set to, it simplifies checking for --version.
func (s *Set) Version() (value, ok bool) {
	flag, exists := s.Get("version")
	if !exists {
		// No version defined
		return false, false
	}
	// Is it a bool flag?
	if flag.Type() != format.TypeBool {
		return false, false
	}
	// It is there, we can infer from the string value what it's set to
	// avoid unnecessary type conversions
	return flag.String() == "true", true
}

// Args returns a slice of all the non-flag arguments, including any
// following a "--" terminator.
func (s *Set) Args() []string {
	if s == nil {
		return []string{}
	}

	return s.args
}

// ExtraArgs returns any arguments after a "--" was encountered, or nil
// if there were none.
func (s *Set) ExtraArgs() []string {
	if s == nil {
		return nil
	}

	return s.extra
}

// Parse parses flags and their values from the command line.
func (s *Set) Parse(args []string) (err error) {
	if s == nil {
		return errors.New("Parse called on a nil set")
	}

	for len(args) > 0 {
		arg := args[0]  // The argument we're currently inspecting
		args = args[1:] // Remainder

		if arg == "--" {
			// "--" terminates the flags
			terminatorIndex := len(s.args)
			s.args = append(s.args, args...)
			s.extra = s.args[terminatorIndex:]

			return nil
		}

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

// Usage returns a string containing the usage info of all flags in the set.
func (s *Set) Usage() (string, error) {
	buf := &bytes.Buffer{}
	buf.Grow(usageBufferSize)

	// Flags should be sorted alphabetically
	names := make([]string, 0, len(s.flags))
	for name := range s.flags {
		names = append(names, name)
	}

	slices.Sort(names)

	tw := tabwriter.NewWriter(buf, style.MinWidth, style.TabWidth, style.Padding, style.PadChar, style.Flags)

	for _, name := range names {
		f := s.flags[name]
		if f == nil {
			return "", fmt.Errorf("Value stored against key %s was nil", name) // Should never happen
		}

		var shorthand string
		if f.Short() != flag.NoShortHand {
			shorthand = "-" + string(f.Short())
		} else {
			shorthand = "N/A"
		}

		fmt.Fprintf(
			tw,
			"  %s\t--%s\t%s\t%s\n",
			style.Bold.Text(shorthand),
			style.Bold.Text(name),
			f.Type(),
			f.Usage(),
		)
	}

	if err := tw.Flush(); err != nil {
		return "", fmt.Errorf("could not format flags: %w", err)
	}

	return buf.String(), nil
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
		err := flag.Set(value)
		if err != nil {
			return nil, err
		}

		// We're done, no need to cut anything from rest as this was a single arg
		return rest, nil
	}

	// Must now either be --flag (boolean) or --flag value
	switch {
	case flag.NoArgValue() != "":
		// --flag (boolean)
		err := flag.Set(flag.NoArgValue())
		if err != nil {
			return nil, err
		}
		// Done, as above no need to cut anything
		return rest, nil
	case len(rest) > 0:
		// --flag value
		value := rest[0]

		err := flag.Set(value)
		if err != nil {
			return nil, err
		}
		// Done, cut value from args and return
		return rest[1:], nil
	default:
		// --flag (value was required)
		return nil, fmt.Errorf("flag --%s requires an argument", name)
	}
}

// parseShortFlag parses short flags from the command line. It is passed the possible
// short flag and the rest of the argument list and returns the remaining arguments
// after it's done parsing to the caller.
//
// The forms it expects are "-f", "-vfg", "-f value" and "-f=value".
func (s *Set) parseShortFlag(short string, rest []string) (remaining []string, err error) {
	// Could either be "f", "vfg" or "f=value"
	shorthands := strings.TrimPrefix(short, "-")

	if len(shorthands) == 0 {
		return nil, fmt.Errorf("invalid flag name %q: must not be empty", shorthands)
	}

	// go test inserts flags like "-test.testlogfile"
	if strings.HasPrefix(shorthands, "test.") {
		return rest, nil
	}

	for len(shorthands) > 0 {
		shorthands, rest, err = s.parseSingleShortFlag(shorthands, rest)
		if err != nil {
			return nil, err
		}
	}

	return rest, nil
}

// parseSingleShortFlag parses a single short flag entry.
func (s *Set) parseSingleShortFlag(shorthands string, rest []string) (string, []string, error) {
	for _, char := range shorthands {
		if err := validateFlagShort(char); err != nil {
			return "", nil, fmt.Errorf("invalid flag shorthand %q: %w", string(char), err)
		}

		flag, exists := s.shorthands[char]
		if !exists {
			return "", nil, fmt.Errorf("unrecognised shorthand flag: %q in -%s", string(char), shorthands)
		}

		switch {
		case len(shorthands) > 2 && shorthands[1] == '=':
			// '-f=value'
			value := shorthands[2:]

			err := flag.Set(value)
			if err != nil {
				return "", nil, err
			}
			// No more shorthands to parse as we got given a value
			// Nothing to trim off the arguments as "-f=value" is all 1 arg
			return "", rest, nil

		case flag.NoArgValue() != "":
			// -f with implied value e.g. boolean or count
			err := flag.Set(flag.NoArgValue())
			if err != nil {
				return "", nil, err
			}

			// We've consumed a single short from the string so trim that off
			return shorthands[1:], rest, nil

		case len(shorthands) > 1:
			// '-fvalue'
			value := shorthands[1:]

			err := flag.Set(value)
			if err != nil {
				return "", nil, err
			}

			// No more shorthands to parse as we got given a value
			// Nothing to trim off as "-fvalue" is all 1 arg
			return "", rest, nil

		case len(rest) > 0:
			// '-f value'
			value := rest[0]

			err := flag.Set(value)
			if err != nil {
				return "", nil, err
			}

			// We've consumed an argument from rest as it was the value to this flag
			// as well as a shorthand
			return shorthands[1:], rest[1:], nil

		default:
			// '-f' with required value
			return "", nil, fmt.Errorf("flag %s needs an argument: %q in -%s", flag.Name(), string(char), shorthands)
		}
	}

	// Didn't match any of our rules, must be invalid short flag syntax
	return "", nil, fmt.Errorf("invalid short flag syntax: %s", shorthands)
}
