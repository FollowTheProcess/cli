package flag

import "fmt"

// Set is a set of command line flags.
type Set struct {
	flags      map[string]flagEntry // The actual stored flags, can lookup by name
	shorthands map[rune]flagEntry   // The flags by shorthand
}

// NewSet builds and returns a new set of flags.
func NewSet() *Set {
	return &Set{
		flags:      make(map[string]flagEntry),
		shorthands: make(map[rune]flagEntry),
	}
}

// Add adds a flag to the Set.
func (s *Set) Add(name, usage string, shorthand rune, flag Value) error {
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

// flagEntry represents a single flag in the set.
type flagEntry struct {
	value             Value  // The actual Flag[T]
	name              string // The full name of the flag e.g. "delete"
	usage             string // The flag's usage message
	defaultValue      string // String representation of the default flag value
	defaultValueNoArg string // String representation of the default flag value if used without an arg, e.g. boolean flags "--force" implies "--force true"
	shorthand         rune   // The optional shorthand e.g. 'd' or [NoShortHand]
}
