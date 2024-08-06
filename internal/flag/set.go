package flag

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/FollowTheProcess/cli/internal/colour"
	"github.com/FollowTheProcess/cli/internal/table"
)

// Set is a set of command line flags.
type Set struct {
	flags      map[string]Entry // The actual stored flags, can lookup by name
	shorthands map[rune]Entry   // The flags by shorthand
	args       []string         // Arguments minus flags or flag values
	extra      []string         // Arguments after "--" was hit
}

// NewSet builds and returns a new set of flags.
func NewSet() *Set {
	return &Set{
		flags:      make(map[string]Entry),
		shorthands: make(map[rune]Entry),
	}
}

// AddToSet adds a flag to the given Set.
func AddToSet[T Flaggable](set *Set, flag Flag[T]) error {
	if set == nil {
		return errors.New("cannot add flag to a nil set")
	}
	_, exists := set.flags[flag.name]
	if exists {
		return fmt.Errorf("flag %q already defined", flag.name)
	}
	var defaultValueNoArg string
	if flag.Type() == typeBool {
		// Boolean flags imply passing true, "--force" vs "--force true"
		defaultValueNoArg = boolTrue
	}
	entry := Entry{
		Value:             flag,
		Name:              flag.name,
		Usage:             flag.usage,
		DefaultValue:      flag.String(),
		DefaultValueNoArg: defaultValueNoArg,
		Shorthand:         flag.short,
	}
	set.flags[flag.name] = entry

	// Only add the shorthand if it wasn't opted out of
	if flag.short != NoShortHand {
		set.shorthands[flag.short] = entry
	}

	return nil
}

// Get gets a flag [Entry] from the Set by name and a boolean to indicate
// whether it was present.
func (s *Set) Get(name string) (Entry, bool) {
	if s == nil {
		return Entry{}, false
	}
	entry, ok := s.flags[name]
	if !ok {
		return Entry{}, false
	}
	return entry, true
}

// GetShort gets a flag [Entry] from the Set by it's shorthand and a boolean to indicate
// whether it was present.
func (s *Set) GetShort(short rune) (Entry, bool) {
	if s == nil {
		return Entry{}, false
	}
	entry, ok := s.shorthands[short]
	if !ok {
		return Entry{}, false
	}
	return entry, true
}

// Help returns whether the [Set] has a boolean flag named "help" and what the value
// of that flag is currently set to, it simplifies checking for --help.
func (s *Set) Help() (value, ok bool) {
	entry, exists := s.Get("help")
	if !exists {
		// No help defined
		return false, false
	}
	// Is it a bool flag?
	if entry.Value.Type() != typeBool {
		return false, false
	}
	// It is there, we can infer from the string value what it's set to
	// to avoid unnecessary type conversions
	return entry.Value.String() == boolTrue, true
}

// Version returns whether the [Set] has a boolean flag named "version" and what the value
// of that flag is currently set to, it simplifies checking for --version.
func (s *Set) Version() (value, ok bool) {
	entry, exists := s.Get("version")
	if !exists {
		// No help defined
		return false, false
	}
	// Is it a bool flag?
	if entry.Value.Type() != typeBool {
		return false, false
	}
	// It is there, we can infer from the string value what it's set to
	// to avoid unnecessary type conversions
	return entry.Value.String() == boolTrue, true
}

// Args returns a slice of all the non-flag arguments, including any
// following a "--" terminator.
func (s *Set) Args() []string {
	if s == nil {
		return nil
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

	// Flags should be sorted alphabetically
	names := make([]string, 0, len(s.flags))
	for name := range s.flags {
		names = append(names, name)
	}
	slices.Sort(names)

	tab := table.New(buf)

	for _, name := range names {
		entry := s.flags[name]
		var shorthand string
		if entry.Shorthand != NoShortHand {
			shorthand = fmt.Sprintf("-%s", string(entry.Shorthand))
		} else {
			shorthand = "N/A"
		}

		tab.Row("  %s\t--%s\t%s\t%s\n", colour.Bold(shorthand), colour.Bold(entry.Name), entry.Value.Type(), entry.Usage)
	}

	if err := tab.Flush(); err != nil {
		return "", fmt.Errorf("could not format flags: %w", err)
	}

	return buf.String(), nil
}

// Entry represents a single flag in the set, as stored.
type Entry struct {
	Value             Value  // The actual Flag[T]
	Name              string // The full name of the flag e.g. "delete"
	Usage             string // The flag's usage message
	DefaultValue      string // String representation of the default flag value
	DefaultValueNoArg string // String representation of the default flag value if used without an arg, e.g. boolean flags "--force" implies "--force true"
	Shorthand         rune   // The optional shorthand e.g. 'd' or [NoShortHand]
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
		err := flag.Value.Set(value)
		if err != nil {
			return nil, err
		}

		// We're done, no need to cut anything from rest as this was a single arg
		return rest, nil
	}

	// Must now either be --flag (boolean) or --flag value
	switch {
	case flag.DefaultValueNoArg != "":
		// --flag (boolean)
		err := flag.Value.Set(flag.DefaultValueNoArg)
		if err != nil {
			return nil, err
		}
		// Done, as above no need to cut anything
		return rest, nil
	case len(rest) > 0:
		// --flag value
		value := rest[0]
		err := flag.Value.Set(value)
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

	// go test inserts flags like "-test.testlogfile"
	if strings.HasPrefix(shorthands, "test.") {
		return rest, nil
	}

	// Is it e.g. f=value
	shorthand, value, containsEquals := strings.Cut(shorthands, "=")
	if err := validateFlagName(shorthand); err != nil {
		return nil, fmt.Errorf("invalid flag name %q: %w", shorthand, err)
	}

	switch {
	case containsEquals:
		// Yes, it is. If the thing on the left of the equals is > 1 char it's an error
		return s.parseShortFlagEqualsValue(shorthand, value, rest)

	case len(shorthands) > 1:
		// It must be e.g. -fvalue and 'shorthands' here refers to the string "fvalue"
		return s.parseShortFlagValue(shorthands, rest)

	case len(rest) > 0:
		// It must be "f value" and value is the next argument in rest
		return s.parseShortFlagArgValue(shorthands, rest)

	default:
		// If we get here, it must be the "-vvv" form
		for _, char := range shorthands {
			flag, exists := s.shorthands[char]
			if !exists {
				return nil, fmt.Errorf("unrecognised shorthand flag: %q in -%s", string(char), shorthands)
			}

			// -f (boolean flag)
			if flag.DefaultValueNoArg != "" {
				err := flag.Value.Set(flag.DefaultValueNoArg)
				if err != nil {
					return nil, err
				}
			}
		}
		return rest, nil
	}
}

// parseShortFlagEqualsValue parses a -f=value style short flag.
//
// It is passed the 'short' which is the left hand side of the equals sign (minus the '-'), the 'value' which is the right
// and the remaining args 'rest'. It returns the remaining args after it's finished parsing and any error.
func (s *Set) parseShortFlagEqualsValue(short string, value string, rest []string) ([]string, error) {
	// Yes, it is. If the thing on the left of the equals is > 1 char it's an error
	if utf8.RuneCountInString(short) > 1 {
		return nil, fmt.Errorf("invalid shorthand syntax: expected e.g. -f=<value> got -%s=%s", short, value)
	}

	char, _ := utf8.DecodeRuneInString(short)
	if err := validateFlagShort(char); err != nil {
		return nil, fmt.Errorf("invalid flag shorthand %q: %w", string(char), err)
	}

	flag, exists := s.shorthands[char]
	if !exists {
		return nil, fmt.Errorf("unrecognised shorthand flag: -%s", string(char))
	}

	if err := flag.Value.Set(value); err != nil {
		return nil, err
	}

	// We're done, nothing to trim off as "-f=value" is all 1 arg
	return rest, nil
}

// parseShortFlagValue parses a -fvalue style short flag.
//
// It is passed the 'shorthands' which is the entire string following the '-', so in our example it would be
// "fvalue", and the remaining args 'rest'. It returns the remaining args after it's finished parsing
// and any error.
func (s *Set) parseShortFlagValue(shorthands string, rest []string) ([]string, error) {
	char, _ := utf8.DecodeRuneInString(shorthands)
	if err := validateFlagShort(char); err != nil {
		return nil, fmt.Errorf("invalid flag shorthand %q: %w", string(char), err)
	}
	value := shorthands[1:]
	flag, exists := s.shorthands[char]
	if !exists {
		return nil, fmt.Errorf("unrecognised shorthand flag: -%s", string(char))
	}

	// TODO: The parsing "-vvv" logic also belongs here, the switch should be if the
	// flag needs a value or not. If the flag is a boolean or a count it's ok to do a
	// -vvv, otherwise the flag needs an argument so this would be an error

	if err := flag.Value.Set(value); err != nil {
		return nil, err
	}

	// We're done, nothing to trim off as "-fvalue" is all 1 arg
	return rest, nil
}

// parseShortFlagArgValue parses an -f value style short flag.
//
// It is passed the 'shorthands' which is the entire string following the '-', which for this flag style
// should just be a single char "f", and the remaining args 'rest'.
// It returns the remaining args after it's finished parsing and any error.
func (s *Set) parseShortFlagArgValue(shorthands string, rest []string) ([]string, error) {
	char, _ := utf8.DecodeRuneInString(shorthands)
	if err := validateFlagShort(char); err != nil {
		return nil, fmt.Errorf("invalid flag shorthand %q: %w", string(char), err)
	}
	value := rest[0]
	flag, exists := s.shorthands[char]
	if !exists {
		return nil, fmt.Errorf("unrecognised shorthand flag: -%s", string(char))
	}
	if err := flag.Value.Set(value); err != nil {
		return nil, err
	}

	// We've consumed "value" from rest so trim it off
	return rest[1:], nil
}
