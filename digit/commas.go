package digit

import (
	"strconv"
	"strings"
)

// GroupSeparator is the ASCII thousands separator used by [AddCommas].
// For the Persian separator see [ThousandsSeparator].
const GroupSeparator = ","

// AddCommas formats value in base 10 with an ASCII comma between every group of
// three digits, e.g. 14555478854 becomes "14,555,478,854".
func AddCommas(value int64) string {
	digits := strconv.FormatInt(value, 10)

	sign := ""
	if strings.HasPrefix(digits, "-") {
		sign, digits = "-", digits[1:]
	}

	var b strings.Builder
	b.Grow(len(sign) + len(digits) + (len(digits)-1)/groupSize)
	b.WriteString(sign)
	for i, r := range digits {
		if i > 0 && (len(digits)-i)%groupSize == 0 {
			b.WriteString(GroupSeparator)
		}
		b.WriteRune(r)
	}
	return b.String()
}

// RemoveCommas parses a grouped number string into an int64. ASCII and Persian
// separators are both accepted, as are Persian and Arabic-Indic digits, so
// "۱۲۳،۴۵۶" and "123,456" parse alike.
//
// It returns an error when the remaining text is not a valid integer.
func RemoveCommas(s string) (int64, error) {
	return ParseInt(s)
}

// ParseInt parses a number written with any mix of ASCII, Persian and
// Arabic-Indic digits, ignoring ASCII and Persian thousands separators and
// surrounding space. A leading minus sign is honored.
func ParseInt(s string) (int64, error) {
	s = strings.TrimSpace(ToEnglishDigits(s))
	s = strings.NewReplacer(GroupSeparator, "", ThousandsSeparator, "").Replace(s)
	return strconv.ParseInt(s, 10, 64)
}
