// Package arg provides mechanisms for defining and configuring command line arguments.
package arg

import (
	"net"
	"time"
)

// Argable is a type constraint that defines any type capable of being parsed as a command line arg.
type Argable interface {
	// TODO(@FollowTheProcess): Slices of stuff
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
		time.Time |
		time.Duration |
		net.IP
}
