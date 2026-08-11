// Package text provides utilities for working with Persian (Farsi) text.
//
// It covers the text-shaped problems that show up in almost every
// Persian-facing application:
//
//   - Normalization: folding Arabic look-alike letters onto their Persian
//     counterparts so text can be compared, sorted and indexed reliably.
//     See [FixArabic] and [Normalize].
//   - Keyboard: recovering text typed with the wrong keyboard layout.
//     See [SwitchToPersianKey] and [SwitchToEnglishKey].
//   - Transliteration: romanizing Persian words. See [Finglish].
//   - Inspection: [Reverse], [CheckIsEnglish] and [OnlyPersianAlpha].
//
// For digits, money and number words see the digit package.
//
// Every function is a pure function of its input, holds no state and is safe
// for concurrent use.
//
// # A note on runes
//
// All functions operate on runes, never on bytes, so multi-byte Persian
// characters are handled correctly. Input that is not valid UTF-8 is not
// rejected: invalid bytes simply do not match any rule and are dropped by the
// filtering functions.
package text

import "github.com/amiranmanesh/go-persian-tools/internal/runeutil"

// Zero-width control characters that shape Persian words.
const (
	// ZWNJ is the zero-width non-joiner, U+200C. Persian uses it to separate
	// the parts of a compound word without letting them join visually, as in
	// the word "mi-ravam", written as می + ZWNJ + روم.
	ZWNJ = "\u200C"
	// ZWJ is the zero-width joiner, U+200D. It forces a joining form where the
	// letters would otherwise stand apart.
	ZWJ = "\u200D"
)

// Reverse returns s with its runes in reverse order. Combining marks are moved
// along with the rest, so reversing decomposed text is not visually lossless.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// CheckIsEnglish reports whether text is made up exclusively of ASCII letters.
// Digits, spaces and punctuation all make it report false; the empty string
// reports true.
func CheckIsEnglish(text string) bool {
	for _, r := range text {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

// OnlyPersianAlpha keeps the characters of text that live in the Arabic Unicode
// block (U+0600..U+06FF), plus decimal points, and discards everything else.
// The block covers Persian letters, Persian digits and Persian punctuation.
func OnlyPersianAlpha(text string) string {
	return runeutil.Keep(text, func(r rune) bool {
		return isArabicBlock(r) || r == '.'
	})
}

func isArabicBlock(r rune) bool { return r >= 0x0600 && r <= 0x06FF }
