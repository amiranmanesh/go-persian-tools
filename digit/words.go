package digit

import "strings"

// The building blocks of a Persian number word.
var (
	onesWords = [...]string{
		"صفر", "یک", "دو", "سه", "چهار", "پنج", "شش", "هفت", "هشت", "نه",
	}
	teensWords = [...]string{
		"ده", "یازده", "دوازده", "سیزده", "چهارده",
		"پانزده", "شانزده", "هفده", "هجده", "نوزده",
	}
	tensWords = [...]string{
		"", "", "بیست", "سی", "چهل", "پنجاه", "شصت", "هفتاد", "هشتاد", "نود",
	}
	hundredsWords = [...]string{
		"", "صد", "دویست", "سیصد", "چهارصد",
		"پانصد", "ششصد", "هفتصد", "هشتصد", "نهصد",
	}
	// scaleWords[i] names the 10^(3i) group. int64 reaches just past 9
	// quintillion, so the table stops there.
	scaleWords = [...]string{
		"", "هزار", "میلیون", "میلیارد", "تریلیون", "کوادریلیون", "کوینتیلیون",
	}
)

// Persian joins the parts of a number word with this conjunction.
const wordSeparator = " و "

// Zero and the negative prefix.
const (
	zeroWord     = "صفر"
	negativeWord = "منفی"
)

// ToWords writes n out in Persian words, e.g. 156789 becomes
// "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه". Negative numbers are prefixed
// with "منفی".
func ToWords(n int64) string {
	if n == 0 {
		return zeroWord
	}

	// Take the magnitude in uint64 so that math.MinInt64, whose negation does
	// not fit in an int64, is handled like any other value.
	negative := n < 0
	magnitude := uint64(n)
	if negative {
		magnitude = -uint64(n)
	}

	words := unsignedToWords(magnitude)
	if negative {
		return negativeWord + " " + words
	}
	return words
}

// unsignedToWords writes a non-zero magnitude out in Persian words.
func unsignedToWords(n uint64) string {
	// Split into three-digit groups, least significant first.
	var groups []uint64
	for n > 0 {
		groups = append(groups, n%1000)
		n /= 1000
	}

	parts := make([]string, 0, len(groups))
	// Walk back down so the most significant group is written first.
	for i := len(groups) - 1; i >= 0; i-- {
		group := groups[i]
		if group == 0 {
			continue
		}

		scale := scaleWords[i]
		// Persian says "هزار", not "یک هزار", for a bare thousand. Every
		// larger scale keeps its "یک".
		if group == 1 && i == 1 {
			parts = append(parts, scale)
			continue
		}

		words := groupToWords(group)
		if scale != "" {
			words += " " + scale
		}
		parts = append(parts, words)
	}

	return strings.Join(parts, wordSeparator)
}

// groupToWords writes a number in 1..999 out in Persian words.
func groupToWords(n uint64) string {
	parts := make([]string, 0, 3)

	if h := n / 100; h > 0 {
		parts = append(parts, hundredsWords[h])
	}

	switch rest := n % 100; {
	case rest == 0:
		// Nothing below the hundreds.
	case rest < 10:
		parts = append(parts, onesWords[rest])
	case rest < 20:
		parts = append(parts, teensWords[rest-10])
	default:
		parts = append(parts, tensWords[rest/10])
		if ones := rest % 10; ones > 0 {
			parts = append(parts, onesWords[ones])
		}
	}

	return strings.Join(parts, wordSeparator)
}

// ToWord converts a numeric string into its Persian word representation,
// e.g. "156789" becomes "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه".
//
// Persian, Arabic-Indic and ASCII digits are all accepted, and a leading minus
// sign is honored. It returns an empty string when the input is not a valid
// integer or does not fit in an int64.
func ToWord(text string) string {
	n, err := ParseInt(text)
	if err != nil {
		return ""
	}
	return ToWords(n)
}

// DigitToWord is the former name of [ToWord].
//
// Deprecated: use [ToWord].
func DigitToWord(text string) string { return ToWord(text) } //nolint:revive // deprecated alias kept for compatibility
