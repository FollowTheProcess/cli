package flag

import (
	"net"
	"time"
)

// TODO(@FollowTheProcess): Support []<everything below> where you can use the flag
// multiple times to append values e.g. --slice "one" --slice "two" --slice "three"
// produces []string{"one", "two", "three"}

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
		[]int8
}
