// Package digit works with numbers in Persian text: converting between the
// ASCII, Persian and Arabic-Indic digit sets, extracting digits from mixed
// text, grouping amounts with either separator, spelling numbers out in
// Persian words, and formatting money in Toman or Rial.
//
// Every function is a pure function of its input, holds no state and is safe
// for concurrent use.
package digit

import (
	"strconv"
	"strings"

	"github.com/amiranmanesh/go-persian-tools/internal/runeutil"
)

// The three digit sets this package understands.
const (
	// EnglishDigits are the ASCII digits, U+0030..U+0039.
	EnglishDigits = "0123456789"
	// PersianDigits are the Extended Arabic-Indic digits used in Persian,
	// U+06F0..U+06F9.
	PersianDigits = "۰۱۲۳۴۵۶۷۸۹"
	// ArabicDigits are the Arabic-Indic digits used in Arabic, U+0660..U+0669.
	// They look close enough to Persian digits that they routinely leak into
	// Persian input, so [ToEnglishDigits] accepts them too.
	ArabicDigits = "٠١٢٣٤٥٦٧٨٩"
)

var (
	persianDigitsReplacer = strings.NewReplacer(digitPairs(EnglishDigits, PersianDigits)...)

	englishDigitsReplacer = strings.NewReplacer(append(
		digitPairs(PersianDigits, EnglishDigits),
		digitPairs(ArabicDigits, EnglishDigits)...,
	)...)
)

// digitPairs zips two equally long digit sets into the flat old/new slice that
// [strings.NewReplacer] expects.
func digitPairs(from, to string) []string {
	src, dst := []rune(from), []rune(to)
	pairs := make([]string, 0, len(src)*2)
	for i, r := range src {
		pairs = append(pairs, string(r), string(dst[i]))
	}
	return pairs
}

// ToPersianDigits converts every ASCII digit in text to its Persian
// counterpart. Any other rune is copied through untouched.
func ToPersianDigits(text string) string {
	return persianDigitsReplacer.Replace(text)
}

// ToPersianDigitsFromInt formats value in base 10 using Persian digits.
// A negative value keeps its ASCII minus sign.
func ToPersianDigitsFromInt(value int) string {
	return ToPersianDigits(strconv.Itoa(value))
}

// ToEnglishDigits converts every Persian and Arabic-Indic digit in text to the
// matching ASCII digit. Any other rune is copied through untouched.
func ToEnglishDigits(text string) string {
	return englishDigitsReplacer.Replace(text)
}

// OnlyEnglishNumbers keeps the ASCII digits and decimal points of text and
// discards everything else.
func OnlyEnglishNumbers(text string) string {
	return runeutil.Keep(text, func(r rune) bool {
		return isEnglishDigit(r) || r == '.'
	})
}

// OnlyPersianNumbers keeps the Persian digits and decimal points of text and
// discards everything else.
func OnlyPersianNumbers(text string) string {
	return runeutil.Keep(text, func(r rune) bool {
		return isPersianDigit(r) || r == '.'
	})
}

// OnlyNumbers keeps the ASCII digits, Persian digits and decimal points of text
// and discards everything else.
func OnlyNumbers(text string) string {
	return runeutil.Keep(text, func(r rune) bool {
		return isEnglishDigit(r) || isPersianDigit(r) || r == '.'
	})
}

func isEnglishDigit(r rune) bool { return r >= '0' && r <= '9' }

func isPersianDigit(r rune) bool { return r >= '۰' && r <= '۹' }
