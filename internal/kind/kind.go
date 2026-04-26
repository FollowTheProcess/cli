// Package kind defines a compact type tag identifying the underlying
// concrete type of a Flag or Arg value.
//
// It exists so that hot paths do not have to do type switching which
// boosts performance and cuts allocations.
package kind

// Kind identifies the underlying concrete type of a Flag or Arg value.
type Kind uint8

// Concrete kinds for every type in the public flag.Flaggable / arg.Argable
// constraints.
const (
	Invalid Kind = iota
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	String
	Bool
	BytesHex
	Count
	Time
	Duration
	IP
	URL
	IntSlice
	Int8Slice
	Int16Slice
	Int32Slice
	Int64Slice
	UintSlice
	Uint16Slice
	Uint32Slice
	Uint64Slice
	Float32Slice
	Float64Slice
	StringSlice
)
