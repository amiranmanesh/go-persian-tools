package digit

import "strings"

// Persian money formatting.
const (
	// ThousandsSeparator is the Arabic comma, U+060C, used to group digits.
	ThousandsSeparator = "،"
	// TomanSuffix is appended by [Toman].
	TomanSuffix = " تومان"
	// RialSuffix is appended by [Rial]. It is the Rial sign, U+FDFC.
	RialSuffix = " ﷼"
)

// groupSize is how many digits sit between two separators.
const groupSize = 3

// Currency formats amount in the Persian money style: digits grouped in threes
// from the right, separated by [ThousandsSeparator], written with Persian
// digits.
//
// Anything in amount that is not a digit is discarded first, so both
// "1234567" and "1,234,567" work, and either digit set is accepted. The
// function is meant for whole amounts; a decimal point is kept but is not
// treated as a decimal separator when grouping.
func Currency(amount string) string {
	digits := []rune(OnlyNumbers(amount))
	if len(digits) == 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(len(digits)*3 + len(digits)/groupSize*len(ThousandsSeparator))
	for i, r := range digits {
		if i > 0 && (len(digits)-i)%groupSize == 0 {
			b.WriteString(ThousandsSeparator)
		}
		b.WriteRune(r)
	}
	return ToPersianDigits(b.String())
}

// Toman formats amount as [Currency] does and appends the Toman unit.
func Toman(amount string) string {
	return Currency(amount) + TomanSuffix
}

// Rial formats amount as [Currency] does and appends the Rial sign.
func Rial(amount string) string {
	return Currency(amount) + RialSuffix
}
