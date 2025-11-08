// Package constraints provides generic constraints for cli.
//
// It is roughly like golang/x/exp/constraints except the types here are more strict (no ~) and
// only Signed and Unsigned are provided.
package constraints

import "go.followtheprocess.codes/cli/flag"

// Signed is the same as constraints.Signed but we don't have to depend
// on golang/x/exp.
type Signed interface {
	int | int8 | int16 | int32 | int64
}

// Unsigned is the same as constraints.Unsigned (with Count mixed in) but we don't have to depend
// on golang/x/exp.
type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64 | uintptr | flag.Count
}
