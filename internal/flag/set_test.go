package flag_test

import (
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/FollowTheProcess/cli/internal/flag"
	"github.com/FollowTheProcess/test"
)

var (
	debug  = goflag.Bool("debug", false, "Print debug output during tests")
	update = goflag.Bool("update", false, "Update golden files")
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string                            // The name of the test case
		errMsg  string                            // If we did get an error, what should it say
		newSet  func(t *testing.T) *flag.Set      // Function to create a set for testing
		test    func(t *testing.T, set *flag.Set) // Function to invoke to test the set's functionality
		args    []string                          // The arguments provided to parse
		wantErr bool                              // Whether we want a parse error
	}{
		{
			name: "empty set no args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("something")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "empty set args no flags",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("something")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(t, set.Args(), []string{"some", "args", "here", "no", "flags"}, slices.Equal)
			},
			args:    []string{"some", "args", "here", "no", "flags"},
			wantErr: false,
		},
		{
			name: "empty set args with terminator",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("something")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(t, set.Args(), []string{"some", "args", "here", "no", "flags", "extra", "args"}, slices.Equal)
				test.EqualFunc(t, set.ExtraArgs(), []string{"extra", "args"}, slices.Equal)
			},
			args:    []string{"some", "args", "here", "no", "flags", "--", "extra", "args"},
			wantErr: false,
		},
		{
			name: "empty set with flags",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("flag")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(t, set.Args(), []string{"some", "args", "here"}, slices.Equal)
			},
			args:    []string{"some", "args", "here", "--flag", "-s"},
			wantErr: true,
			errMsg:  "unrecognised flag: --flag",
		},
		{
			name: "undefined flag long",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--unknown"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-u"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -u`,
		},
		{
			name: "undefined flag long with value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--unknown", "value"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short with value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-u", "value"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -u`,
		},
		{
			name: "undefined flag long equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--unknown=value"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-u=value"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -u=value`,
		},
		{
			name: "undefined flag shortvalue",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-uvalue"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -uvalue`,
		},
		{
			name: "bad syntax short empty name",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-"},
			wantErr: true,
			errMsg:  `invalid flag name "": must not be empty`,
		},
		{
			name: "bad syntax long extra hyphen",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"---"},
			wantErr: true,
			errMsg:  `invalid flag name "-": trailing hyphen`,
		},
		{
			name: "bad syntax long leading whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-- delete"},
			wantErr: true,
			errMsg:  `invalid flag name " delete": cannot contain whitespace`,
		},
		{
			name: "bad syntax short leading whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"- d"},
			wantErr: true,
			errMsg:  `invalid flag shorthand " ": cannot contain whitespace`,
		},
		{
			name: "bad syntax long trailing whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--delete "},
			wantErr: true,
			errMsg:  `invalid flag name "delete ": cannot contain whitespace`,
		},
		{
			name: "bad syntax short trailing whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				set := flag.NewSet()

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			args:    []string{"-d "},
			wantErr: true,
			errMsg:  `invalid flag shorthand " ": cannot contain whitespace`,
		},
		{
			name: "bad syntax short more than 1 char equals",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-dfv=something"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "d" in -dfv=something`,
		},
		{
			name: "bad syntax short non utf8",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-Ê"},
			wantErr: true,
			errMsg:  `invalid flag shorthand "Ê": invalid character, must be a single ASCII letter, got "Ê"`,
		},
		{
			name: "bad syntax short non utf8 equals",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-Ê=something"},
			wantErr: true,
			errMsg:  `invalid flag shorthand "Ê": invalid character, must be a single ASCII letter, got "Ê"`,
		},
		{
			name: "bad syntax short multiple non utf8",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"-本¼語"},
			wantErr: true,
			errMsg:  `invalid flag shorthand "本": invalid character, must be a single ASCII letter, got "本"`,
		},
		{
			name: "bad syntax long internal whitespace",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			args:    []string{"--de lete"},
			wantErr: true,
			errMsg:  `invalid flag name "de lete": cannot contain whitespace`,
		},
		{
			name: "valid long",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), "true")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"--delete"},
			wantErr: false,
		},
		{
			name: "valid long with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), "true")

				test.EqualFunc(t, set.Args(), []string{"some", "subcommand", "extra", "args"}, slices.Equal)
				test.EqualFunc(t, set.ExtraArgs(), []string{"extra", "args"}, slices.Equal)
			},
			args:    []string{"some", "subcommand", "--delete", "--", "extra", "args"},
			wantErr: false,
		},
		{
			name: "valid short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), "true")

				// Get by short
				flag, exists = set.GetShort('d')
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), "true")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"-d"},
			wantErr: false,
		},
		{
			name: "valid shortvalue",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "number", 'n', 0, "Number of something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "42")

				// Get by short
				flag, exists = set.GetShort('n')
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "42")
			},
			args:    []string{"-n42"},
			wantErr: false,
		},
		{
			name: "valid short with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), "true")

				// Get by short
				flag, exists = set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), "true")

				test.EqualFunc(t, set.Args(), []string{"some", "arg"}, slices.Equal)
			},
			args:    []string{"some", "arg", "-d"},
			wantErr: false,
		},
		{
			name: "valid long value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"--count", "1"},
			wantErr: false,
		},
		{
			name: "valid long missing value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "0")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"--count"}, // Count needs an argument
			wantErr: true,
			errMsg:  "flag --count requires an argument",
		},
		{
			name: "valid short missing value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "0")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"-c"}, // Count needs an argument
			wantErr: true,
			errMsg:  `flag count needs an argument: "c" in -c`,
		},
		{
			name: "valid long value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"some", "arg", "more", "args"}, slices.Equal)
			},
			args:    []string{"some", "arg", "--count", "1", "more", "args"},
			wantErr: false,
		},
		{
			name: "invalid long value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', 0, "Uint")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "uint")
				test.Equal(t, flag.String(), "0") // Shouldn't have been set
			},
			args:    []string{"--number", "-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `flag "number" received invalid value "-8" (expected uint), detail: strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid short value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				// Get by short
				flag, exists = set.GetShort('c')
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"-c", "1"},
			wantErr: false,
		},
		{
			name: "valid short value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				// Get by short
				flag, exists = set.GetShort('c')
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"args", "more", "args"}, slices.Equal)
			},
			args:    []string{"args", "-c", "1", "more", "args"},
			wantErr: false,
		},
		{
			name: "invalid short value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', 0, "Uint")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "uint")
				test.Equal(t, flag.String(), "0") // Shouldn't have been set
			},
			args:    []string{"-n", "-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `flag "number" received invalid value "-8" (expected uint), detail: strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid long equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")
			},
			args:    []string{"--count=1"},
			wantErr: false,
		},
		{
			name: "valid long equals value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"args", "more", "args"}, slices.Equal)
			},
			args:    []string{"args", "--count=1", "more", "args"},
			wantErr: false,
		},
		{
			name: "invalid long equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', 0, "Uint")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "uint")
				test.Equal(t, flag.String(), "0") // Shouldn't have been set
			},
			args:    []string{"--number=-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `flag "number" received invalid value "-8" (expected uint), detail: strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid short equals value",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				// Get by short
				flag, exists = set.GetShort('c')
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")
			},
			args:    []string{"-c=1"},
			wantErr: false,
		},
		{
			name: "valid short equals value with args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				// Get by short
				flag, exists = set.GetShort('c')
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				test.EqualFunc(t, set.Args(), []string{"args", "more", "args"}, slices.Equal)
			},
			args:    []string{"args", "-c=1", "more", "args"},
			wantErr: false,
		},
		{
			name: "no shorthand use long",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", flag.NoShortHand, 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "1")

				// Get by short
				flag, exists = set.GetShort('c')
				test.False(t, exists) // Short shouldn't exist
				test.Equal(t, flag, nil)
			},
			args:    []string{"--count", "1"},
			wantErr: false,
		},
		{
			name: "no shorthand use short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", flag.NoShortHand, 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				// Get by name
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "int")
				test.Equal(t, flag.String(), "0")

				// Get by short
				flag, exists = set.GetShort('c')
				test.False(t, exists) // Short shouldn't exist
				test.Equal(t, flag, nil)
			},
			args:    []string{"-c", "1"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "c" in -c`,
		},
		{
			name: "valid count long",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(flag.Count), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "count")
				test.Equal(t, flag.String(), "2") // Should be incremented by 2
			},
			args:    []string{"--count", "--count"},
			wantErr: false,
		},
		{
			name: "valid count short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(flag.Count), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "count")
				test.Equal(t, flag.String(), "2") // Should be incremented by 2
			},
			args:    []string{"-c", "-c"},
			wantErr: false,
		},
		{
			name: "valid count super short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()
				f, err := flag.New(new(flag.Count), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "count")
				test.Equal(t, flag.String(), "3") // Should be incremented by 3
			},
			args:    []string{"-ccc"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := tt.newSet(t)
			err := set.Parse(tt.args)
			test.WantErr(t, err, tt.wantErr)

			if err != nil {
				test.Equal(t, err.Error(), tt.errMsg)
			}

			if tt.test != nil {
				tt.test(t, set)
			}
		})
	}
}

func TestFlagSet(t *testing.T) {
	tests := []struct {
		newSet func(t *testing.T) *flag.Set      // Function to build the flag set under test
		test   func(t *testing.T, set *flag.Set) // Function to test the set
		name   string                            // Name of the test case
	}{
		{
			name: "empty",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("missing")
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "empty short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.GetShort('d')
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe get",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.Get("missing")
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe get short",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, exists := set.GetShort('m')
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe add",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, err := flag.New(new(bool), "force", 'f', false, "Force something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), "cannot add flag to a nil set")
				}
			},
		},
		{
			name: "nil safe parse",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				err := set.Parse([]string{"args", "here", "doesn't", "matter"})
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), "Parse called on a nil set")
				}
			},
		},
		{
			name: "nil safe args",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
		},
		{
			name: "duplicate flag added",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				// Add the same flag again
				err = flag.AddToSet(set, f)
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), `flag "count" already defined`)
				}
			},
		},
		{
			name: "different flag same name added",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				f2, err := flag.New(new(uint), "count", 'c', 0, "Count something 2")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				// Add the different flag with the same na,e
				err = flag.AddToSet(set, f2)
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), `flag "count" already defined`)
				}
			},
		},
		{
			name: "different flag same short added",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				f, err := flag.New(new(int), "count", 'c', 0, "Count something")
				test.Ok(t, err)

				f2, err := flag.New(new(string), "config", 'c', "", "Choose a config file")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				// Add the different flag with the same na,e
				err = flag.AddToSet(set, f2)
				test.Err(t, err)
				if err != nil {
					test.Equal(t, err.Error(), `shorthand "c" already in use for flag "count"`)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := tt.newSet(t)
			tt.test(t, set)
		})
	}
}

func TestHelpVersion(t *testing.T) {
	tests := []struct {
		newSet func(t *testing.T) *flag.Set      // Function to build the flag set under test
		test   func(t *testing.T, set *flag.Set) // Function to test the set
		name   string                            // Name of the test case
	}{
		{
			name: "empty help",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				help, ok := set.Help()
				test.False(t, help)
				test.False(t, ok)
			},
		},
		{
			name: "empty version",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				version, ok := set.Version()
				test.False(t, version)
				test.False(t, ok)
			},
		},
		{
			name: "help false",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()

				f, err := flag.New(new(bool), "help", 'h', false, "Show help")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				help, ok := set.Help()
				test.True(t, ok)    // It should exist
				test.False(t, help) // But not true
			},
		},
		{
			name: "help non bool",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()

				f, err := flag.New(new(int), "help", 'h', 0, "Show help")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				help, ok := set.Help()
				test.False(t, ok)   // It should not exist
				test.False(t, help) // And be false because it's not a bool
			},
		},
		{
			name: "help true",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()

				f, err := flag.New(new(bool), "help", 'h', true, "Show help")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				help, ok := set.Help()
				test.True(t, ok)   // It should exist
				test.True(t, help) // And be true
			},
		},
		{
			name: "version false",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()

				f, err := flag.New(new(bool), "version", 'v', false, "Show version")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				version, ok := set.Version()
				test.True(t, ok)       // It should exist
				test.False(t, version) // But not true
			},
		},
		{
			name: "version non bool",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()

				f, err := flag.New(new(int), "version", 'v', 0, "Show version")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				version, ok := set.Version()
				test.False(t, ok)      // It should not exist
				test.False(t, version) // And be false because it's not a bool
			},
		},
		{
			name: "version true",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				set := flag.NewSet()

				f, err := flag.New(new(bool), "version", 'v', true, "Show version")
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				t.Helper()
				version, ok := set.Version()
				test.True(t, ok)      // It should exist
				test.True(t, version) // And be true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := tt.newSet(t)
			tt.test(t, set)
		})
	}
}

func TestUsage(t *testing.T) {
	tests := []struct {
		newSet func(t *testing.T) *flag.Set // Function to build the flag set under test
		name   string                       // Name of the test case
		golden string                       // Name of the file containing expected output
	}{
		{
			name: "simple",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				help, err := flag.New(new(bool), "help", 'h', false, "Show help for test")
				test.Ok(t, err)

				version, err := flag.New(new(bool), "version", 'V', false, "Show version info for test")
				test.Ok(t, err)

				set := flag.NewSet()

				err = flag.AddToSet(set, help)
				test.Ok(t, err)

				err = flag.AddToSet(set, version)
				test.Ok(t, err)

				return set
			},
			golden: "simple.txt",
		},
		{
			name: "no shorthand",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				help, err := flag.New(new(bool), "help", 'h', false, "Show help for test")
				test.Ok(t, err)

				version, err := flag.New(new(bool), "version", 'V', false, "Show version info for test")
				test.Ok(t, err)

				up, err := flag.New(new(bool), "update", flag.NoShortHand, false, "Update something")
				test.Ok(t, err)

				set := flag.NewSet()

				err = flag.AddToSet(set, help)
				test.Ok(t, err)

				err = flag.AddToSet(set, version)
				test.Ok(t, err)

				err = flag.AddToSet(set, up)
				test.Ok(t, err)

				return set
			},
			golden: "no-shorthand.txt",
		},
		{
			name: "full",
			newSet: func(t *testing.T) *flag.Set {
				t.Helper()
				help, err := flag.New(new(bool), "help", 'h', false, "Show help for test")
				test.Ok(t, err)

				version, err := flag.New(new(bool), "version", 'V', false, "Show version info for test")
				test.Ok(t, err)

				up, err := flag.New(new(bool), "update", flag.NoShortHand, false, "Update something")
				test.Ok(t, err)

				count, err := flag.New(new(int), "count", 'c', 0, "Count things")
				test.Ok(t, err)

				thing, err := flag.New(new(string), "thing", 't', "", "Name the thing")
				test.Ok(t, err)

				set := flag.NewSet()

				err = flag.AddToSet(set, help)
				test.Ok(t, err)

				err = flag.AddToSet(set, version)
				test.Ok(t, err)

				err = flag.AddToSet(set, up)
				test.Ok(t, err)

				err = flag.AddToSet(set, count)
				test.Ok(t, err)

				err = flag.AddToSet(set, thing)
				test.Ok(t, err)

				return set
			},
			golden: "full.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Force no colour in tests
			t.Setenv("NO_COLOR", "true")
			set := tt.newSet(t)
			golden := filepath.Join(test.Data(t), "TestUsage", tt.golden)

			got, err := set.Usage()
			test.Ok(t, err)

			if *debug {
				fmt.Printf("DEBUG (%s)\n_____\n\n%s\n", tt.name, got)
			}

			if *update {
				t.Logf("Updating %s\n", golden)
				err := os.WriteFile(golden, []byte(got), os.ModePerm)
				test.Ok(t, err)
			}

			test.File(t, got, golden)
		})
	}
}

func BenchmarkParse(b *testing.B) {
	set := flag.NewSet()
	f, err := flag.New(new(int), "count", 'c', 0, "Count something")
	test.Ok(b, err)

	f2, err := flag.New(new(bool), "delete", 'd', false, "Delete something")
	test.Ok(b, err)

	f3, err := flag.New(new(string), "name", 'n', "", "Name something")
	test.Ok(b, err)

	err = flag.AddToSet(set, f)
	test.Ok(b, err)

	err = flag.AddToSet(set, f2)
	test.Ok(b, err)

	err = flag.AddToSet(set, f3)
	test.Ok(b, err)

	args := []string{"some", "args", "here", "--delete", "--count", "10", "--name", "John"}

	b.ResetTimer()
	for range b.N {
		err := set.Parse(args)
		if err != nil {
			b.Fatalf("Parse returned an error: %v", err)
		}
	}
}
