package table_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/FollowTheProcess/cli/internal/table"
	"github.com/FollowTheProcess/test"
)

var (
	debug  = flag.Bool("debug", false, "Print debug output during tests")
	update = flag.Bool("update", false, "Update golden files")
)

func TestTable(t *testing.T) {
	buf := &bytes.Buffer{}

	tab := table.New(buf)

	tab.Row("Col1\tCol2\tCol3\n")
	tab.Row("val1\tval2\tval3\n")
	tab.Row("val4\tval5\tval6\n")

	err := tab.Flush()
	test.Ok(t, err)

	file := filepath.Join(test.Data(t), "table.txt")

	if *debug {
		fmt.Printf("DEBUG (%s)\n_____\n\n%s\n", "TestTable", buf.String())
	}

	if *update {
		t.Logf("Updating %s\n", file)
		err := os.WriteFile(file, buf.Bytes(), os.ModePerm)
		test.Ok(t, err)
	}

	test.File(t, buf.String(), file)
}
