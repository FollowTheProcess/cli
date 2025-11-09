// Package style simply provides a uniform terminal printing style via [hue] for use across
// the library.
//
// [hue]: https://github.com/FollowTheProcess/hue
package style

import (
	"io"

	"go.followtheprocess.codes/hue"
	"go.followtheprocess.codes/hue/tabwriter"
)

const (
	// Title is the style for titles of help text sections like arguments or commands.
	Title = hue.Bold | hue.Underline

	// Bold is simply plain bold text.
	Bold = hue.Bold

	// minWidth is the minimum cell width for hue's colour-enabled tabwriter.
	minWidth = 1

	// tabWidth is the width of tabs in spaces for tabwriter.
	tabWidth = 8

	// padding is the number of PadChars to pad table cells with.
	padding = 2

	// padChar is the character with which to pad table cells.
	padChar = ' '

	// flags is the tabwriter config flags.
	flags = 0
)

// Tabwriter returns a [hue.Tabwriter] configured with cli house style.
func Tabwriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, minWidth, tabWidth, padding, padChar, flags)
}
