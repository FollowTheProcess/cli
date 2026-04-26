package flag_test

import (
	"iter"
	"maps"
	"net"
	"slices"
	"testing"
	"time"

	publicflag "go.followtheprocess.codes/cli/flag"
	"go.followtheprocess.codes/cli/internal/flag"
	"go.followtheprocess.codes/cli/internal/format"
	"go.followtheprocess.codes/test"
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
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
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
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
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
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("something")
				test.False(t, exists)
				test.Equal(t, f, nil)

				test.EqualFunc(
					t,
					set.Args(),
					[]string{"some", "args", "here", "no", "flags", "extra", "args"},
					slices.Equal,
				)
				test.EqualFunc(t, set.ExtraArgs(), []string{"extra", "args"}, slices.Equal)
			},
			args:    []string{"some", "args", "here", "no", "flags", "--", "extra", "args"},
			wantErr: false,
		},
		{
			name: "empty set with flags",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
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
				return flag.NewSet()
			},
			args:    []string{"--unknown"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"-u"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -u`,
		},
		{
			name: "undefined flag long with value",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"--unknown", "value"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short with value",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"-u", "value"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -u`,
		},
		{
			name: "undefined flag long equals value",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"--unknown=value"},
			wantErr: true,
			errMsg:  "unrecognised flag: --unknown",
		},
		{
			name: "undefined flag short equals value",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"-u=value"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -u=value`,
		},
		{
			name: "undefined flag shortvalue",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"-uvalue"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "u" in -uvalue`,
		},
		{
			name: "bad syntax short empty name",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"-"},
			wantErr: true,
			errMsg:  `invalid flag name "": must not be empty`,
		},
		{
			name: "bad syntax short more than 1 char equals",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			args:    []string{"-dfv=something"},
			wantErr: true,
			errMsg:  `unrecognised shorthand flag: "d" in -dfv=something`,
		},
		{
			name: "valid long",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), format.True)

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"--delete"},
			wantErr: false,
		},
		{
			name: "valid long with args",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), format.True)

				test.EqualFunc(t, set.Args(), []string{"some", "subcommand", "extra", "args"}, slices.Equal)
				test.EqualFunc(t, set.ExtraArgs(), []string{"extra", "args"}, slices.Equal)
			},
			args:    []string{"some", "subcommand", "--delete", "--", "extra", "args"},
			wantErr: false,
		},
		{
			name: "valid short",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				// Get by name
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), format.True)

				// Get by short
				flag, exists = set.GetShort('d')
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), format.True)

				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
			args:    []string{"-d"},
			wantErr: false,
		},
		{
			name: "valid shortvalue",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "number", 'n', "Number of something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				// Get by name
				flag, exists := set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), format.True)

				// Get by short
				flag, exists = set.Get("delete")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "bool")
				test.Equal(t, flag.String(), format.True)

				test.EqualFunc(t, set.Args(), []string{"some", "arg"}, slices.Equal)
			},
			args:    []string{"some", "arg", "-d"},
			wantErr: false,
		},
		{
			name: "valid long value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', "Uint", flag.Config[uint]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				flag, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "uint")
				test.Equal(t, flag.String(), "0") // Shouldn't have been set
			},
			args:    []string{"--number", "-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `parse error: flag "number" received invalid value "-8" (expected uint): strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid short value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', "Uint", flag.Config[uint]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				flag, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "uint")
				test.Equal(t, flag.String(), "0") // Shouldn't have been set
			},
			args:    []string{"-n", "-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `parse error: flag "number" received invalid value "-8" (expected uint): strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid long equals value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(uint), "number", 'n', "Uint", flag.Config[uint]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				flag, exists := set.Get("number")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "uint")
				test.Equal(t, flag.String(), "0") // Shouldn't have been set
			},
			args:    []string{"--number=-8"}, // Trying to set a uint flag to negative number
			wantErr: true,
			errMsg:  `parse error: flag "number" received invalid value "-8" (expected uint): strconv.ParseUint: parsing "-8": invalid syntax`,
		},
		{
			name: "valid short equals value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", publicflag.NoShortHand, "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", publicflag.NoShortHand, "Count something", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(publicflag.Count), "count", 'c', "Count something", flag.Config[publicflag.Count]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(publicflag.Count), "count", 'c', "Count something", flag.Config[publicflag.Count]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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
				set := flag.NewSet()
				f, err := flag.New(new(publicflag.Count), "count", 'c', "Count something", flag.Config[publicflag.Count]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				flag, exists := set.Get("count")
				test.True(t, exists)

				test.Equal(t, flag.Type(), "count")
				test.Equal(t, flag.String(), "3") // Should be incremented by 3
			},
			args:    []string{"-ccc"},
			wantErr: false,
		},
		{
			name: "env var applied when flag not set on CLI",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_DELETE", "true")

				var val bool

				f, err := flag.New(&val, "delete", 'd', "Delete something", flag.Config[bool]{EnvVar: "MYTOOL_DELETE"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("delete")
				test.True(t, exists)
				test.Equal(t, f.String(), format.True)
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "CLI long flag overrides env var",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_DELETE", "true")

				var val bool

				f, err := flag.New(&val, "delete", 'd', "Delete something", flag.Config[bool]{EnvVar: "MYTOOL_DELETE"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("delete")
				test.True(t, exists)
				test.Equal(t, f.String(), "false")
			},
			args:    []string{"--delete=false"},
			wantErr: false,
		},
		{
			name: "CLI short flag overrides env var",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_DELETE", "true")

				var val bool

				f, err := flag.New(&val, "delete", 'd', "Delete something", flag.Config[bool]{EnvVar: "MYTOOL_DELETE"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("delete")
				test.True(t, exists)
				test.Equal(t, f.String(), "false")
			},
			args:    []string{"-d=false"},
			wantErr: false,
		},
		{
			name: "invalid env var value returns error",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_COUNT", "notanumber")

				var val int

				f, err := flag.New(&val, "count", 'c', "Count", flag.Config[int]{EnvVar: "MYTOOL_COUNT"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			args:    []string{},
			wantErr: true,
			errMsg:  `could not set flag from env: env var MYTOOL_COUNT: parse error: flag "count" received invalid value "notanumber" (expected int): strconv.ParseInt: parsing "notanumber": invalid syntax`,
		},
		{
			name: "env var not set leaves flag at default",
			newSet: func(t *testing.T) *flag.Set {
				// MYTOOL_DELETE is deliberately NOT set in the environment
				var val bool

				f, err := flag.New(&val, "delete", 'd', "Delete something", flag.Config[bool]{EnvVar: "MYTOOL_DELETE"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("delete")
				test.True(t, exists)
				test.Equal(t, f.String(), "false")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "env var applied when terminator present",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_DELETE", "true")

				var val bool

				f, err := flag.New(&val, "delete", 'd', "Delete something", flag.Config[bool]{EnvVar: "MYTOOL_DELETE"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("delete")
				test.True(t, exists)
				test.Equal(t, f.String(), format.True)
				test.EqualFunc(t, set.ExtraArgs(), []string{"extra"}, slices.Equal)
			},
			args:    []string{"--", "extra"},
			wantErr: false,
		},
		{
			name: "flag without env var configured is unaffected by environment",
			newSet: func(t *testing.T) *flag.Set {
				// Env var is set in the OS but NOT wired to this flag
				t.Setenv("MYTOOL_DELETE", "true")

				var val bool

				f, err := flag.New(&val, "delete", 'd', "Delete something", flag.Config[bool]{})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("delete")
				test.True(t, exists)
				test.Equal(t, f.String(), "false")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "slice flag set via comma-separated env var",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_ITEMS", "one,two,three")

				var val []string

				f, err := flag.New(&val, "item", 'i', "Add item", flag.Config[[]string]{EnvVar: "MYTOOL_ITEMS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("item")
				test.True(t, exists)
				test.Equal(t, f.String(), `["one", "two", "three"]`)
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "count env var and CLI flags both accumulate",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_VERBOSITY", "2")

				var val publicflag.Count

				f, err := flag.New(&val, "verbosity", 'v', "Increase verbosity", flag.Config[publicflag.Count]{EnvVar: "MYTOOL_VERBOSITY"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("verbosity")
				test.True(t, exists)
				// Env var contributes 2, CLI contributes 1 more — total 3
				test.Equal(t, f.String(), "3")
			},
			args:    []string{"--verbosity"},
			wantErr: false,
		},
		{
			name: "slice env var and CLI flags both accumulate",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_ITEMS", "one,two")

				var val []string

				f, err := flag.New(&val, "item", 'i', "Add item", flag.Config[[]string]{EnvVar: "MYTOOL_ITEMS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("item")
				test.True(t, exists)
				test.Equal(t, f.String(), `["one", "two", "three"]`)
			},
			args:    []string{"--item", "three"},
			wantErr: false,
		},
		{
			name: "net.IP flag set via env var",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_HOST", "192.168.1.1")

				var val net.IP

				f, err := flag.New(&val, "host", 'h', "Host IP address", flag.Config[net.IP]{EnvVar: "MYTOOL_HOST"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("host")
				test.True(t, exists)
				test.Equal(t, f.String(), "192.168.1.1")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "net.IP flag env var is overridden by CLI",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_HOST", "192.168.1.1")

				var val net.IP

				f, err := flag.New(&val, "host", 'h', "Host IP address", flag.Config[net.IP]{EnvVar: "MYTOOL_HOST"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("host")
				test.True(t, exists)
				test.Equal(t, f.String(), "10.0.0.1")
			},
			args:    []string{"--host", "10.0.0.1"},
			wantErr: false,
		},
		{
			name: "invalid net.IP env var returns error",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_HOST", "notanip")

				var val net.IP

				f, err := flag.New(&val, "host", 'h', "Host IP address", flag.Config[net.IP]{EnvVar: "MYTOOL_HOST"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			args:    []string{},
			wantErr: true,
			errMsg:  `could not set flag from env: env var MYTOOL_HOST: parse error: flag "host" received invalid value "notanip" (expected net.IP): invalid IP address`,
		},
		{
			name: "time.Duration flag set via env var",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_TIMEOUT", "30s")

				var val time.Duration

				f, err := flag.New(&val, "timeout", 't', "Request timeout", flag.Config[time.Duration]{EnvVar: "MYTOOL_TIMEOUT"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("timeout")
				test.True(t, exists)
				test.Equal(t, f.String(), "30s")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "time.Duration flag env var is overridden by CLI",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_TIMEOUT", "30s")

				var val time.Duration

				f, err := flag.New(&val, "timeout", 't', "Request timeout", flag.Config[time.Duration]{EnvVar: "MYTOOL_TIMEOUT"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("timeout")
				test.True(t, exists)
				test.Equal(t, f.String(), "1m0s")
			},
			args:    []string{"--timeout", "1m"},
			wantErr: false,
		},
		{
			name: "invalid time.Duration env var returns error",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_TIMEOUT", "notaduration")

				var val time.Duration

				f, err := flag.New(&val, "timeout", 't', "Request timeout", flag.Config[time.Duration]{EnvVar: "MYTOOL_TIMEOUT"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			args:    []string{},
			wantErr: true,
			errMsg:  `could not set flag from env: env var MYTOOL_TIMEOUT: parse error: flag "timeout" received invalid value "notaduration" (expected time.Duration): time: invalid duration "notaduration"`,
		},
		{
			name: "time.Time flag set via env var",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_SINCE", "2024-08-17T10:37:30Z")

				var val time.Time

				f, err := flag.New(
					&val,
					"since",
					publicflag.NoShortHand,
					"Start time (RFC3339)",
					flag.Config[time.Time]{EnvVar: "MYTOOL_SINCE"},
				)
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("since")
				test.True(t, exists)
				test.Equal(t, f.String(), "2024-08-17T10:37:30Z")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "invalid time.Time env var returns error",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_SINCE", "not-a-time")

				var val time.Time

				f, err := flag.New(
					&val,
					"since",
					publicflag.NoShortHand,
					"Start time (RFC3339)",
					flag.Config[time.Time]{EnvVar: "MYTOOL_SINCE"},
				)
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			args:    []string{},
			wantErr: true,
			errMsg:  `could not set flag from env: env var MYTOOL_SINCE: parse error: flag "since" received invalid value "not-a-time" (expected time.Time): parsing time "not-a-time" as "2006-01-02T15:04:05Z07:00": cannot parse "not-a-time" as "2006"`,
		},
		{
			name: "[]int flag set via comma-separated env var",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_PORTS", "8080,8081,8082")

				var val []int

				f, err := flag.New(&val, "port", 'p', "Port numbers", flag.Config[[]int]{EnvVar: "MYTOOL_PORTS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("port")
				test.True(t, exists)
				test.Equal(t, f.String(), "[8080, 8081, 8082]")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "invalid value in []int comma-separated env var returns error",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_PORTS", "8080,notaport,8082")

				var val []int

				f, err := flag.New(&val, "port", 'p', "Port numbers", flag.Config[[]int]{EnvVar: "MYTOOL_PORTS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			args:    []string{},
			wantErr: true,
			errMsg:  `could not set flag from env: env var MYTOOL_PORTS: parse error: flag "port" (type []int) cannot append element "notaport": strconv.ParseInt: parsing "notaport": invalid syntax`,
		},
		{
			name: "bytes flag env var is parsed atomically, not split on comma",
			newSet: func(t *testing.T) *flag.Set {
				// A hex string that contains a comma should be parsed as-is, not split
				t.Setenv("MYTOOL_DATA", "deadbeef")

				var val []byte

				f, err := flag.New(&val, "data", 'd', "Raw bytes", flag.Config[[]byte]{EnvVar: "MYTOOL_DATA"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("data")
				test.True(t, exists)
				test.Equal(t, f.String(), "deadbeef")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "int accepts negative value via long space form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-42")
			},
			args:    []string{"--count", "-42"},
			wantErr: false,
		},
		{
			name: "int accepts negative value via long equals form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-42")
			},
			args:    []string{"--count=-42"},
			wantErr: false,
		},
		{
			name: "int accepts negative value via short space form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-42")
			},
			args:    []string{"-c", "-42"},
			wantErr: false,
		},
		{
			name: "int accepts negative value via short equals form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-42")
			},
			args:    []string{"-c=-42"},
			wantErr: false,
		},
		{
			name: "int accepts negative value via short attached form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int), "count", 'c', "Count", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-42")
			},
			args:    []string{"-c-42"},
			wantErr: false,
		},
		{
			name: "int8 accepts negative value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int8), "count", 'c', "Count", flag.Config[int8]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-128")
			},
			args:    []string{"--count", "-128"},
			wantErr: false,
		},
		{
			name: "int16 accepts negative value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int16), "count", 'c', "Count", flag.Config[int16]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-32768")
			},
			args:    []string{"--count", "-32768"},
			wantErr: false,
		},
		{
			name: "int32 accepts negative value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int32), "count", 'c', "Count", flag.Config[int32]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-2147483648")
			},
			args:    []string{"--count=-2147483648"},
			wantErr: false,
		},
		{
			name: "int64 accepts negative value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(int64), "count", 'c', "Count", flag.Config[int64]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("count")
				test.True(t, exists)
				test.Equal(t, f.String(), "-9223372036854775808")
			},
			args:    []string{"-c", "-9223372036854775808"},
			wantErr: false,
		},
		{
			name: "float accepts negative value",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(float64), "ratio", 'r', "Ratio", flag.Config[float64]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("ratio")
				test.True(t, exists)
				test.Equal(t, f.String(), "-3.14")
			},
			args:    []string{"--ratio", "-3.14"},
			wantErr: false,
		},
		{
			name: "string accepts empty value via long space form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(string), "name", 'n', "Name", flag.Config[string]{DefaultValue: "placeholder"})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("name")
				test.True(t, exists)
				// Default was "placeholder"; confirm it's been overwritten to empty
				test.Equal(t, f.String(), "")
			},
			args:    []string{"--name", ""},
			wantErr: false,
		},
		{
			name: "string accepts empty value via long equals form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(string), "name", 'n', "Name", flag.Config[string]{DefaultValue: "placeholder"})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("name")
				test.True(t, exists)
				test.Equal(t, f.String(), "")
			},
			args:    []string{"--name="},
			wantErr: false,
		},
		{
			name: "string accepts empty value via short space form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(string), "name", 'n', "Name", flag.Config[string]{DefaultValue: "placeholder"})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("name")
				test.True(t, exists)
				test.Equal(t, f.String(), "")
			},
			args:    []string{"-n", ""},
			wantErr: false,
		},
		{
			name: "string accepts empty value via short equals form",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()
				f, err := flag.New(new(string), "name", 'n', "Name", flag.Config[string]{DefaultValue: "placeholder"})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("name")
				test.True(t, exists)
				// -n= should be symmetric with --name= (both set to empty string),
				// matching Go stdlib flag, argparse, clap, urfave, etc.
				test.Equal(t, f.String(), "")
			},
			args:    []string{"-n="},
			wantErr: false,
		},
		{
			name: "slice env var of only commas yields empty slice",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_ITEMS", ",")

				var val []string

				f, err := flag.New(&val, "item", 'i', "Add item", flag.Config[[]string]{EnvVar: "MYTOOL_ITEMS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("item")
				test.True(t, exists)
				test.Equal(t, f.String(), "[]")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "slice env var empty string yields empty slice",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_ITEMS", "")

				var val []string

				f, err := flag.New(&val, "item", 'i', "Add item", flag.Config[[]string]{EnvVar: "MYTOOL_ITEMS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("item")
				test.True(t, exists)
				test.Equal(t, f.String(), "[]")
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "slice env var skips empty and whitespace-only items",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_ITEMS", "a,, ,b")

				var val []string

				f, err := flag.New(&val, "item", 'i', "Add item", flag.Config[[]string]{EnvVar: "MYTOOL_ITEMS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("item")
				test.True(t, exists)
				test.Equal(t, f.String(), `["a", "b"]`)
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "slice env var splits on every comma (no escape mechanism)",
			newSet: func(t *testing.T) *flag.Set {
				// Any comma in an env var value is interpreted as a separator —
				// there is no way to embed a literal comma in a slice item.
				// Users needing commas should pass values via --flag one,two.
				t.Setenv("MYTOOL_ITEMS", "a,b,c")

				var val []string

				f, err := flag.New(&val, "item", 'i', "Add item", flag.Config[[]string]{EnvVar: "MYTOOL_ITEMS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("item")
				test.True(t, exists)
				test.Equal(t, f.String(), `["a", "b", "c"]`)
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "slice env var trims surrounding whitespace from items",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_ITEMS", " a , b , c ")

				var val []string

				f, err := flag.New(&val, "item", 'i', "Add item", flag.Config[[]string]{EnvVar: "MYTOOL_ITEMS"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("item")
				test.True(t, exists)
				test.Equal(t, f.String(), `["a", "b", "c"]`)
			},
			args:    []string{},
			wantErr: false,
		},
		{
			name: "combined short flags with equals value",
			newSet: func(t *testing.T) *flag.Set {
				// -abc=x → a and b are bool and get implied true, c gets "x"
				set := flag.NewSet()

				a, err := flag.New(new(bool), "alpha", 'a', "Alpha", flag.Config[bool]{})
				test.Ok(t, err)
				b, err := flag.New(new(bool), "beta", 'b', "Beta", flag.Config[bool]{})
				test.Ok(t, err)
				c, err := flag.New(new(string), "charlie", 'c', "Charlie", flag.Config[string]{})
				test.Ok(t, err)

				test.Ok(t, flag.AddToSet(set, a))
				test.Ok(t, flag.AddToSet(set, b))
				test.Ok(t, flag.AddToSet(set, c))

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				a, ok := set.Get("alpha")
				test.True(t, ok)
				test.Equal(t, a.String(), format.True)

				b, ok := set.Get("beta")
				test.True(t, ok)
				test.Equal(t, b.String(), format.True)

				c, ok := set.Get("charlie")
				test.True(t, ok)
				test.Equal(t, c.String(), "x")
			},
			args:    []string{"-abc=x"},
			wantErr: false,
		},
		{
			name: "combined short flags with space separated value",
			newSet: func(t *testing.T) *flag.Set {
				// -abc x → a, b bool; c consumes next arg as value
				set := flag.NewSet()

				a, err := flag.New(new(bool), "alpha", 'a', "Alpha", flag.Config[bool]{})
				test.Ok(t, err)
				b, err := flag.New(new(bool), "beta", 'b', "Beta", flag.Config[bool]{})
				test.Ok(t, err)
				c, err := flag.New(new(string), "charlie", 'c', "Charlie", flag.Config[string]{})
				test.Ok(t, err)

				test.Ok(t, flag.AddToSet(set, a))
				test.Ok(t, flag.AddToSet(set, b))
				test.Ok(t, flag.AddToSet(set, c))

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				a, ok := set.Get("alpha")
				test.True(t, ok)
				test.Equal(t, a.String(), format.True)

				b, ok := set.Get("beta")
				test.True(t, ok)
				test.Equal(t, b.String(), format.True)

				c, ok := set.Get("charlie")
				test.True(t, ok)
				test.Equal(t, c.String(), "x")
			},
			args:    []string{"-abc", "x"},
			wantErr: false,
		},
		{
			name: "combined short flags with multiple bools and trailing equals value",
			newSet: func(t *testing.T) *flag.Set {
				// -abcf=value → a, b, c bool; f non-bool gets "value"
				set := flag.NewSet()

				a, err := flag.New(new(bool), "alpha", 'a', "Alpha", flag.Config[bool]{})
				test.Ok(t, err)
				b, err := flag.New(new(bool), "beta", 'b', "Beta", flag.Config[bool]{})
				test.Ok(t, err)
				c, err := flag.New(new(bool), "charlie", 'c', "Charlie", flag.Config[bool]{})
				test.Ok(t, err)
				fl, err := flag.New(new(string), "file", 'f', "File", flag.Config[string]{})
				test.Ok(t, err)

				test.Ok(t, flag.AddToSet(set, a))
				test.Ok(t, flag.AddToSet(set, b))
				test.Ok(t, flag.AddToSet(set, c))
				test.Ok(t, flag.AddToSet(set, fl))

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				a, ok := set.Get("alpha")
				test.True(t, ok)
				test.Equal(t, a.String(), format.True)

				b, ok := set.Get("beta")
				test.True(t, ok)
				test.Equal(t, b.String(), format.True)

				c, ok := set.Get("charlie")
				test.True(t, ok)
				test.Equal(t, c.String(), format.True)

				fl, ok := set.Get("file")
				test.True(t, ok)
				test.Equal(t, fl.String(), "value")
			},
			args:    []string{"-abcf=value"},
			wantErr: false,
		},
		{
			name: "bool explicit equals false beats env var true",
			newSet: func(t *testing.T) *flag.Set {
				// Env var pre-sets to true; CLI --verbose=false must win.
				t.Setenv("MYTOOL_VERBOSE", "true")

				var val bool

				f, err := flag.New(&val, "verbose", 'v', "Verbose", flag.Config[bool]{EnvVar: "MYTOOL_VERBOSE"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("verbose")
				test.True(t, exists)
				test.Equal(t, f.String(), "false")
			},
			args:    []string{"--verbose=false"},
			wantErr: false,
		},
		{
			name: "bool short equals false beats env var true",
			newSet: func(t *testing.T) *flag.Set {
				t.Setenv("MYTOOL_VERBOSE", "true")

				var val bool

				f, err := flag.New(&val, "verbose", 'v', "Verbose", flag.Config[bool]{EnvVar: "MYTOOL_VERBOSE"})
				test.Ok(t, err)

				set := flag.NewSet()
				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("verbose")
				test.True(t, exists)
				test.Equal(t, f.String(), "false")
			},
			args:    []string{"-v=false"},
			wantErr: false,
		},
		{
			name: "combined short flags where early flag needs value captures rest",
			newSet: func(t *testing.T) *flag.Set {
				// -fgh — f is non-bool so consumes the rest of the cluster as its value
				// i.e. f gets "gh", and g/h are not parsed as flags.
				set := flag.NewSet()

				fl, err := flag.New(new(string), "file", 'f', "File", flag.Config[string]{})
				test.Ok(t, err)
				g, err := flag.New(new(bool), "gamma", 'g', "Gamma", flag.Config[bool]{})
				test.Ok(t, err)
				h, err := flag.New(new(bool), "hotel", 'h', "Hotel", flag.Config[bool]{})
				test.Ok(t, err)

				test.Ok(t, flag.AddToSet(set, fl))
				test.Ok(t, flag.AddToSet(set, g))
				test.Ok(t, flag.AddToSet(set, h))

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				fl, ok := set.Get("file")
				test.True(t, ok)
				test.Equal(t, fl.String(), "gh")

				g, ok := set.Get("gamma")
				test.True(t, ok)
				test.Equal(t, g.String(), "false")

				h, ok := set.Get("hotel")
				test.True(t, ok)
				test.Equal(t, h.String(), "false")
			},
			args:    []string{"-fgh"},
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

func TestParseResetsState(t *testing.T) {
	t.Run("positional args do not accumulate", func(t *testing.T) {
		set := flag.NewSet()

		f, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
		test.Ok(t, err)
		test.Ok(t, flag.AddToSet(set, f))

		args := []string{"one", "two", "--delete"}

		test.Ok(t, set.Parse(args))
		test.EqualFunc(t, set.Args(), []string{"one", "two"}, slices.Equal)

		test.Ok(t, set.Parse(args))
		test.EqualFunc(t, set.Args(), []string{"one", "two"}, slices.Equal)
	})

	t.Run("extra args do not accumulate", func(t *testing.T) {
		set := flag.NewSet()

		args := []string{"pos", "--", "extra", "more"}

		test.Ok(t, set.Parse(args))
		test.EqualFunc(t, set.ExtraArgs(), []string{"extra", "more"}, slices.Equal)

		test.Ok(t, set.Parse(args))
		test.EqualFunc(t, set.ExtraArgs(), []string{"extra", "more"}, slices.Equal)
	})

	t.Run("previous extras cleared when next call has none", func(t *testing.T) {
		set := flag.NewSet()

		test.Ok(t, set.Parse([]string{"a", "--", "x"}))
		test.EqualFunc(t, set.ExtraArgs(), []string{"x"}, slices.Equal)

		test.Ok(t, set.Parse([]string{"a", "b"}))
		test.Equal(t, len(set.ExtraArgs()), 0)
	})
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
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("missing")
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "empty short",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.GetShort('d')
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe get",
			newSet: func(t *testing.T) *flag.Set {
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.Get("missing")
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe get short",
			newSet: func(t *testing.T) *flag.Set {
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				f, exists := set.GetShort('m')
				test.False(t, exists)
				test.Equal(t, f, nil)
			},
		},
		{
			name: "nil safe add",
			newSet: func(t *testing.T) *flag.Set {
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				f, err := flag.New(new(bool), "force", 'f', "Force something", flag.Config[bool]{})
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
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
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
				return nil // uh oh
			},
			test: func(t *testing.T, set *flag.Set) {
				test.EqualFunc(t, set.Args(), nil, slices.Equal)
			},
		},
		{
			name: "duplicate flag added",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
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
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				f2, err := flag.New(new(uint), "count", 'c', "Count something 2", flag.Config[uint]{})
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
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
				test.Ok(t, err)

				f2, err := flag.New(new(string), "config", 'c', "Choose a config file", flag.Config[string]{})
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
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				help, ok := set.Help()
				test.False(t, help)
				test.False(t, ok)
			},
		},
		{
			name: "empty version",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				version, ok := set.Version()
				test.False(t, version)
				test.False(t, ok)
			},
		},
		{
			name: "help false",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				f, err := flag.New(new(bool), "help", 'h', "Show help", flag.Config[bool]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				help, ok := set.Help()
				test.True(t, ok)    // It should exist
				test.False(t, help) // But not true
			},
		},
		{
			name: "help non bool",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				f, err := flag.New(new(int), "help", 'h', "Show help", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				help, ok := set.Help()
				test.False(t, ok)   // It should not exist
				test.False(t, help) // And be false because it's not a bool
			},
		},
		{
			name: "help true",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				f, err := flag.New(new(bool), "help", 'h', "Show help", flag.Config[bool]{DefaultValue: true})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				help, ok := set.Help()
				test.True(t, ok)   // It should exist
				test.True(t, help) // And be true
			},
		},
		{
			name: "version false",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				f, err := flag.New(new(bool), "version", 'v', "Show version", flag.Config[bool]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				version, ok := set.Version()
				test.True(t, ok)       // It should exist
				test.False(t, version) // But not true
			},
		},
		{
			name: "version non bool",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				f, err := flag.New(new(int), "version", 'v', "Show version", flag.Config[int]{})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				version, ok := set.Version()
				test.False(t, ok)      // It should not exist
				test.False(t, version) // And be false because it's not a bool
			},
		},
		{
			name: "version true",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				f, err := flag.New(new(bool), "version", 'v', "Show version", flag.Config[bool]{DefaultValue: true})
				test.Ok(t, err)

				err = flag.AddToSet(set, f)
				test.Ok(t, err)

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
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

func TestSorted(t *testing.T) {
	tests := []struct {
		newSet func(t *testing.T) *flag.Set
		test   func(t *testing.T, set *flag.Set)
		name   string
	}{
		{
			name: "empty",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				// Iterator should yield no values
				got := maps.Collect(set.Sorted())
				test.Equal(t, len(got), 0)
			},
		},
		{
			name: "full",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				verbose, err := flag.New(new(bool), "verbose", 'v', "Show verbose info", flag.Config[bool]{})
				test.Ok(t, err)

				debug, err := flag.New(new(bool), "debug", 'd', "Show debug info", flag.Config[bool]{})
				test.Ok(t, err)

				thing, err := flag.New(new(string), "thing", 't', "A thing", flag.Config[string]{})
				test.Ok(t, err)

				number, err := flag.New(new(int), "number", 'n', "Number of times", flag.Config[int]{})
				test.Ok(t, err)

				duration, err := flag.New(new(time.Duration), "duration", 'D', "The time to do something for", flag.Config[time.Duration]{})
				test.Ok(t, err)

				test.Ok(t, flag.AddToSet(set, verbose))
				test.Ok(t, flag.AddToSet(set, debug))
				test.Ok(t, flag.AddToSet(set, thing))
				test.Ok(t, flag.AddToSet(set, number))
				test.Ok(t, flag.AddToSet(set, duration))

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				next, stop := iter.Pull2(set.Sorted())
				defer stop()

				// Should now be in alphabetical order
				name, fl, ok := next()
				test.True(t, ok)
				test.Equal(t, name, "debug")
				test.Equal(t, fl.Name(), "debug")

				name, fl, ok = next()
				test.True(t, ok)
				test.Equal(t, name, "duration")
				test.Equal(t, fl.Name(), "duration")

				name, fl, ok = next()
				test.True(t, ok)
				test.Equal(t, name, "number")
				test.Equal(t, fl.Name(), "number")

				name, fl, ok = next()
				test.True(t, ok)
				test.Equal(t, name, "thing")
				test.Equal(t, fl.Name(), "thing")

				name, fl, ok = next()
				test.True(t, ok)
				test.Equal(t, name, "verbose")
				test.Equal(t, fl.Name(), "verbose")

				// Thats it
				name, fl, ok = next()
				test.False(t, ok)
				test.Equal(t, name, "")
				test.Equal(t, fl, nil)
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

func TestAll(t *testing.T) {
	tests := []struct {
		newSet func(t *testing.T) *flag.Set
		test   func(t *testing.T, set *flag.Set)
		name   string
	}{
		{
			name: "empty",
			newSet: func(t *testing.T) *flag.Set {
				return flag.NewSet()
			},
			test: func(t *testing.T, set *flag.Set) {
				// Iterator should yield no values
				got := maps.Collect(set.All())
				test.Equal(t, len(got), 0)
			},
		},
		{
			name: "full",
			newSet: func(t *testing.T) *flag.Set {
				set := flag.NewSet()

				verbose, err := flag.New(new(bool), "verbose", 'v', "Show verbose info", flag.Config[bool]{})
				test.Ok(t, err)

				debug, err := flag.New(new(bool), "debug", 'd', "Show debug info", flag.Config[bool]{})
				test.Ok(t, err)

				thing, err := flag.New(new(string), "thing", 't', "A thing", flag.Config[string]{})
				test.Ok(t, err)

				number, err := flag.New(new(int), "number", 'n', "Number of times", flag.Config[int]{})
				test.Ok(t, err)

				duration, err := flag.New(new(time.Duration), "duration", 'D', "The time to do something for", flag.Config[time.Duration]{})
				test.Ok(t, err)

				test.Ok(t, flag.AddToSet(set, verbose))
				test.Ok(t, flag.AddToSet(set, debug))
				test.Ok(t, flag.AddToSet(set, thing))
				test.Ok(t, flag.AddToSet(set, number))
				test.Ok(t, flag.AddToSet(set, duration))

				return set
			},
			test: func(t *testing.T, set *flag.Set) {
				// Should get everything back, but order is not deterministic
				all := maps.Collect(set.All())

				want := []string{"verbose", "debug", "thing", "number", "duration"}
				slices.Sort(want)

				got := slices.Collect(maps.Keys(all))
				slices.Sort(got)

				test.EqualFunc(t, got, want, slices.Equal)
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

func BenchmarkParse(b *testing.B) {
	set := flag.NewSet()
	f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
	test.Ok(b, err)

	f2, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
	test.Ok(b, err)

	f3, err := flag.New(new(string), "name", 'n', "Name something", flag.Config[string]{})
	test.Ok(b, err)

	err = flag.AddToSet(set, f)
	test.Ok(b, err)

	err = flag.AddToSet(set, f2)
	test.Ok(b, err)

	err = flag.AddToSet(set, f3)
	test.Ok(b, err)

	args := []string{"some", "args", "here", "--delete", "--count", "10", "--name", "John"}

	for b.Loop() {
		err := set.Parse(args)
		if err != nil {
			b.Fatalf("Parse returned an error: %v", err)
		}
	}
}

// Benchmarks Parse on a flag set invoked exclusively via shorthand flags.
func BenchmarkParseShort(b *testing.B) {
	set := flag.NewSet()
	f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
	test.Ok(b, err)

	f2, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
	test.Ok(b, err)

	f3, err := flag.New(new(string), "name", 'n', "Name something", flag.Config[string]{})
	test.Ok(b, err)

	test.Ok(b, flag.AddToSet(set, f))
	test.Ok(b, flag.AddToSet(set, f2))
	test.Ok(b, flag.AddToSet(set, f3))

	args := []string{"some", "args", "here", "-d", "-c", "10", "-n", "John"}

	for b.Loop() {
		err := set.Parse(args)
		if err != nil {
			b.Fatalf("Parse returned an error: %v", err)
		}
	}
}

// Benchmarks Parse with a mix of long and short flags, representative of how
// users typically invoke a real CLI.
func BenchmarkParseMixed(b *testing.B) {
	set := flag.NewSet()
	f, err := flag.New(new(int), "count", 'c', "Count something", flag.Config[int]{})
	test.Ok(b, err)

	f2, err := flag.New(new(bool), "delete", 'd', "Delete something", flag.Config[bool]{})
	test.Ok(b, err)

	f3, err := flag.New(new(string), "name", 'n', "Name something", flag.Config[string]{})
	test.Ok(b, err)

	f4, err := flag.New(new(bool), "verbose", 'v', "Verbose output", flag.Config[bool]{})
	test.Ok(b, err)

	test.Ok(b, flag.AddToSet(set, f))
	test.Ok(b, flag.AddToSet(set, f2))
	test.Ok(b, flag.AddToSet(set, f3))
	test.Ok(b, flag.AddToSet(set, f4))

	args := []string{"some", "args", "here", "-d", "--count", "10", "-n", "John", "--verbose"}

	for b.Loop() {
		err := set.Parse(args)
		if err != nil {
			b.Fatalf("Parse returned an error: %v", err)
		}
	}
}

// Benchmarks Parse when flags have environment variables attached, so include
// the cost of the Getenv syscall.
func BenchmarkParseWithEnv(b *testing.B) {
	build := func(b *testing.B) (*flag.Set, []string) {
		b.Helper()

		set := flag.NewSet()

		f, err := flag.New(
			new(int),
			"count",
			'c',
			"Count something",
			flag.Config[int]{EnvVar: "BENCH_PARSE_COUNT"},
		)
		test.Ok(b, err)

		f2, err := flag.New(
			new(string),
			"name",
			'n',
			"Name something",
			flag.Config[string]{EnvVar: "BENCH_PARSE_NAME"},
		)
		test.Ok(b, err)

		f3, err := flag.New(
			new(bool),
			"force",
			'f',
			"Force something",
			flag.Config[bool]{EnvVar: "BENCH_PARSE_FORCE"},
		)
		test.Ok(b, err)

		test.Ok(b, flag.AddToSet(set, f))
		test.Ok(b, flag.AddToSet(set, f2))
		test.Ok(b, flag.AddToSet(set, f3))

		return set, []string{"positional"}
	}

	b.Run("unset", func(b *testing.B) {
		set, args := build(b)

		for b.Loop() {
			err := set.Parse(args)
			if err != nil {
				b.Fatalf("Parse returned an error: %v", err)
			}
		}
	})

	b.Run("set", func(b *testing.B) {
		b.Setenv("BENCH_PARSE_COUNT", "10")
		b.Setenv("BENCH_PARSE_NAME", "John")
		b.Setenv("BENCH_PARSE_FORCE", "true")

		set, args := build(b)

		for b.Loop() {
			err := set.Parse(args)
			if err != nil {
				b.Fatalf("Parse returned an error: %v", err)
			}
		}
	})
}

// Benchmarks Parse on a slice flag multiple times. A slice flag is a
// read/write operation.
func BenchmarkParseSliceFlag(b *testing.B) {
	var items []string

	set := flag.NewSet()
	f, err := flag.New(&items, "item", 'i', "An item", flag.Config[[]string]{})
	test.Ok(b, err)

	test.Ok(b, flag.AddToSet(set, f))

	args := []string{"--item", "one", "--item", "two", "--item", "three", "--item", "four"}

	for b.Loop() {
		// Reset between iterations so the slice doesn't grow forever
		items = items[:0]

		err := set.Parse(args)
		if err != nil {
			b.Fatalf("Parse returned an error: %v", err)
		}
	}
}
