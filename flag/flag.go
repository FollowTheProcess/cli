// Package flag provides mechanisms for defining and configuring command line flags.
package flag

import (
	"net"
	"time"
)

// NoShortHand should be passed as the "short" value if the desired flag
// should be the long hand version only e.g. --count, not -c/--count.
const NoShortHand = rune(-1)

// Count is a type used for a flag who's job is to increment a counter, e.g. a "verbosity"
// flag may be used like so "-vvv" which should increase the verbosity level to 3.
//
// Count flags may be used in the following ways:
//   - -vvv
//   - --verbose --verbose --verbose (although not sure why you'd do this)
//   - --verbose=3
//
// All have the same effect of increasing the verbosity level to 3.
//
// --verbose 3 however is not supported, this is due to an internal parsing
// implementation detail.
type Count uint

// Flaggable is a type constraint that defines any type capable of being parsed as a command line flag.
type Flaggable interface {
	int |
		int8 |
		int16 |
		int32 |
		int64 |
		uint |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		uintptr |
		float32 |
		float64 |
		string |
		bool |
		[]byte |
		Count |
		time.Time |
		time.Duration |
		net.IP |
		[]int |
		[]int8 |
		[]int16 |
		[]int32 |
		[]int64 |
		[]uint |
		[]uint16 |
		[]uint32 |
		[]uint64 |
		[]float32 |
		[]float64 |
		[]string
}
