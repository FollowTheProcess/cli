// Package style simply provides a uniform terminal printing style via [hue] for use across
// the library.
//
// [hue]: https://github.com/FollowTheProcess/hue
package style

import "go.followtheprocess.codes/hue"

const (
	// Title is the style for titles of help text sections like arguments or commands.
	Title = hue.Bold | hue.White | hue.Underline

	// Bold is simply plain bold text.
	Bold = hue.Bold

	// MinWidth is the minimum cell width for hue's colour-enabled tabwriter.
	MinWidth = 1

	// TabWidth is the width of tabs in spaces for tabwriter.
	TabWidth = 8

	// Padding is the number of PadChars to pad table cells with.
	Padding = 2

	// PadChar is the character with which to pad table cells.
	PadChar = ' '

	// Flags is the tabwriter config flags.
	Flags = 0
)
