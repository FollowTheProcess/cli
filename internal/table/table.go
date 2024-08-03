// Package table implements a thin wrapper around [text/tabwriter] to keep
// formatting consistent across cli.
package table

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// TableWriter config, used for showing subcommands in help.
const (
	minWidth = 2    // Min cell width
	tabWidth = 8    // Tab width in spaces
	padding  = 1    // Padding
	padChar  = '\t' // Char to pad with
)

// Table is a text table.
type Table struct {
	tw *tabwriter.Writer // The underlying writer
}

// New returns a new [Table], writing to w.
func New(w io.Writer) Table {
	tw := tabwriter.NewWriter(w, minWidth, tabWidth, padding, padChar, tabwriter.TabIndent)
	return Table{tw: tw}
}

// Row adds a row to the [Table].
func (t Table) Row(format string, a ...any) {
	fmt.Fprintf(t.tw, format, a...)
}

// Flush flushes the written data to the writer.
func (t Table) Flush() error {
	return t.tw.Flush()
}
