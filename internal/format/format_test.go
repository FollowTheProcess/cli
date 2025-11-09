package format //nolint:testpackage // I need the base and bits values and don't want to export them.

import (
	"strconv"
	"testing"
	"testing/quick"
)

func TestInt(t *testing.T) {
	//nolint:gocritic // It wants me to "unlambda" this but it's generic so I can't
	test := func(n int) string {
		return Int(n)
	}

	reference := func(n int) string {
		return strconv.FormatInt(int64(n), base10)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestUint(t *testing.T) {
	//nolint:gocritic // It wants me to "unlambda" this but it's generic so I can't
	test := func(n uint) string {
		return Uint(n)
	}

	reference := func(n uint) string {
		return strconv.FormatUint(uint64(n), base10)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat32(t *testing.T) {
	test := Float32

	reference := func(f float32) string {
		return strconv.FormatFloat(float64(f), floatFmt, floatPrecision, bits32)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}

func TestFloat64(t *testing.T) {
	test := Float64

	reference := func(f float64) string {
		return strconv.FormatFloat(float64(f), floatFmt, floatPrecision, bits64)
	}

	if err := quick.CheckEqual(test, reference, nil); err != nil {
		t.Error(err)
	}
}
