package table_test

import (
	"bytes"
	"flag"
	"fmt"
	"testing"

	"followtheprocess.codes/cli/internal/table"
	"github.com/FollowTheProcess/snapshot"
	"github.com/FollowTheProcess/test"
)

var (
	debug  = flag.Bool("debug", false, "Print debug output during tests")
	update = flag.Bool("update", false, "Update golden files")
)

func TestTable(t *testing.T) {
	snap := snapshot.New(t, snapshot.Update(*update))
	buf := &bytes.Buffer{}

	tab := table.New(buf)

	tab.Row("Col1\tCol2\tCol3\n")
	tab.Row("val1\tval2\tval3\n")
	tab.Row("val4\tval5\tval6\n")

	err := tab.Flush()
	test.Ok(t, err)

	if *debug {
		fmt.Printf("DEBUG (%s)\n_____\n\n%s\n", "TestTable", buf.String())
	}

	snap.Snap(buf.String())
}
