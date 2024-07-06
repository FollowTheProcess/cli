package cli_test

import (
	"testing"

	"github.com/FollowTheProcess/cli"
)

func TestHello(t *testing.T) {
	got := cli.Hello()
	want := "Hello cli"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
