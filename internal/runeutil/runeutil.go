// Package runeutil holds the small rune helpers shared by the public
// packages. It is internal: the helpers are implementation details, not API.
package runeutil

import "strings"

// Keep returns the runes of text for which allow reports true, in order.
func Keep(text string, allow func(rune) bool) string {
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		if allow(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Replacer builds a [strings.Replacer] from an ordered table of pairs. When
// reverse is true the pairs are flipped, which is how the two directions of a
// keyboard layout are derived from a single table.
func Replacer(pairs [][2]string, reverse bool) *strings.Replacer {
	flat := make([]string, 0, len(pairs)*2)
	for _, p := range pairs {
		if reverse {
			flat = append(flat, p[1], p[0])
		} else {
			flat = append(flat, p[0], p[1])
		}
	}
	return strings.NewReplacer(flat...)
}

// ConcatPairs joins several replacement tables into one, preserving order.
func ConcatPairs(tables ...[][2]string) [][2]string {
	var all [][2]string
	for _, t := range tables {
		all = append(all, t...)
	}
	return all
}
