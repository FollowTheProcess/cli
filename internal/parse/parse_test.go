package parse //nolint:testpackage // I need the base and bits values and don't want to export them.

import (
	"strconv"
	"testing"
	"testing/quick"
)

// These are basically all just testing that I haven't broken anything
// by wrapping strconv and saves me having to write lots of test cases
// by hand.

func TestInt(t *testing.T) {
	test := Int

	reference := func(str string) (int, error) {
		val, err := strconv.ParseInt(str, base10, 0)
		return int(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt8(t *testing.T) {
	test := Int8

	reference := func(str string) (int8, error) {
		val, err := strconv.ParseInt(str, base10, bits8)
		return int8(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt16(t *testing.T) {
	test := Int16

	reference := func(str string) (int16, error) {
		val, err := strconv.ParseInt(str, base10, bits16)
		return int16(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt32(t *testing.T) {
	test := Int32

	reference := func(str string) (int32, error) {
		val, err := strconv.ParseInt(str, base10, bits32)
		return int32(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestInt64(t *testing.T) {
	test := Int64

	reference := func(str string) (int64, error) {
		val, err := strconv.ParseInt(str, base10, bits64)
		return val, err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint(t *testing.T) {
	test := Uint

	reference := func(str string) (uint, error) {
		val, err := strconv.ParseUint(str, base10, 0)
		return uint(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint8(t *testing.T) {
	test := Uint8

	reference := func(str string) (uint8, error) {
		val, err := strconv.ParseUint(str, base10, bits8)
		return uint8(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint16(t *testing.T) {
	test := Uint16

	reference := func(str string) (uint16, error) {
		val, err := strconv.ParseUint(str, base10, bits16)
		return uint16(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint32(t *testing.T) {
	test := Uint32

	reference := func(str string) (uint32, error) {
		val, err := strconv.ParseUint(str, base10, bits32)
		return uint32(val), err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint64(t *testing.T) {
	test := Uint64

	reference := func(str string) (uint64, error) {
		val, err := strconv.ParseUint(str, base10, bits64)
		return val, err
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}
