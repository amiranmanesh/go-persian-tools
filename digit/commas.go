// Package digit provides helpers to format numbers with thousands separators
// and to convert digits into their Persian word representation.
package digit

import (
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

// AddCommas formats the given integer with thousands separators,
// e.g. 14555478854 becomes "14,555,478,854".
func AddCommas(digit int64) string {
	return humanize.Comma(digit)
}

// RemoveCommas parses a comma-separated number string into an int64.
// It returns an error when the input (with commas removed) is not a valid
// integer.
func RemoveCommas(s string) (int64, error) {
	return strconv.ParseInt(strings.ReplaceAll(s, ",", ""), 10, 64)
}
