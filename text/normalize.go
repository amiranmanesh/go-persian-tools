package text

import "github.com/amiranmanesh/go-persian-tools/internal/runeutil"

// arabicToPersian folds Arabic letters onto the Persian letter they are
// habitually mistaken for. The Arabic yeh and kaf are by far the most common:
// they are visually identical to the Persian ones in most fonts, so text
// copied from Arabic keyboards and older systems mixes them freely, and a
// naive string comparison then says "علي" != "علی".
var arabicToPersian = [][2]string{
	{"ي", "ی"}, // Arabic yeh          U+064A -> Persian yeh    U+06CC
	{"ى", "ی"}, // Arabic alef maksura U+0649 -> Persian yeh    U+06CC
	{"ك", "ک"}, // Arabic kaf          U+0643 -> Persian keheh  U+06A9
}

// normalizeExtra is applied on top of [arabicToPersian] by [Normalize]. These
// rules lose information on purpose: they collapse characters that a reader
// treats as the same letter, so that sorting, searching and equality checks
// behave the way a Persian speaker expects.
var normalizeExtra = [][2]string{
	{ZWJ, " "},
	{ZWNJ, " "},
	{"ٱ", "ا"}, // alef wasla
	{"آ", "ا"}, // alef with madda
	{"أ", "ا"}, // alef with hamza above
	{"ء", "ا"}, // hamza
	{"ئ", "ی"}, // yeh with hamza above
	{"ة", "ه"}, // teh marbuta
}

var (
	fixArabicReplacer = runeutil.Replacer(arabicToPersian, false)
	normalizeReplacer = runeutil.Replacer(runeutil.ConcatPairs(arabicToPersian, normalizeExtra), false)
)

// isHarakat reports whether r is an Arabic vocalization mark. These marks sit
// on top of a letter to spell out a short vowel; Persian writes them only in
// dictionaries, poetry and the Quran, so identical words differ by nothing but
// their presence. They carry no letter identity and are dropped.
//
// The hamza and madda marks (U+0653..U+0655) are deliberately not in the set:
// they do change which letter is being written.
func isHarakat(r rune) bool {
	switch {
	case r >= 0x064B && r <= 0x0652:
		// Tanwin, fatha, damma, kasra, shadda and sukun.
		return true
	case r == 0x0670:
		// Superscript alef.
		return true
	case r == 0x0640:
		// Tatweel, the decorative elongation of a joining stroke.
		return true
	default:
		return false
	}
}

// isCombiningMadda reports whether r is a standalone madda or hamza mark,
// U+0653..U+0655. Unlike the marks in [isHarakat] these do change the letter
// they sit on — an alef carrying U+0654 is أ — so [FixArabic] leaves them
// alone.
//
// [Normalize] drops them, because it already folds every precomposed form of
// those letters (آ أ ٱ ء ئ) onto a plain one. Without this, the ezafe in
// "مقالهٔ" would keep it from matching "مقاله", which is the exact class of
// near-miss Normalize exists to remove.
func isCombiningMadda(r rune) bool { return r >= 0x0653 && r <= 0x0655 }

// removeHarakat strips every Arabic vocalization mark from text.
func removeHarakat(text string) string {
	return runeutil.Keep(text, func(r rune) bool { return !isHarakat(r) })
}

// removeMarks strips the vocalization marks and the standalone madda and hamza.
func removeMarks(text string) string {
	return runeutil.Keep(text, func(r rune) bool {
		return !isHarakat(r) && !isCombiningMadda(r)
	})
}

// FixArabic replaces the Arabic characters in text with their Persian
// equivalents and strips Arabic vocalization marks. It is a conservative,
// meaning-preserving cleanup: use it on anything you are about to store.
func FixArabic(text string) string {
	return fixArabicReplacer.Replace(removeHarakat(text))
}

// Normalize prepares text for comparison, sorting and indexing. On top of
// everything [FixArabic] does, it folds the alef and hamza variants onto plain
// alef, teh marbuta onto heh, drops the standalone madda and hamza marks, and
// turns zero-width joiners into spaces.
//
// The result is meant to be compared, not displayed: "می‌روم" normalizes to
// "می روم" and "مقالهٔ من" to "مقاله من". Keep the original around for showing
// back to the user.
//
// Normalize is idempotent, so callers can safely normalize once when writing a
// record and again when querying it.
func Normalize(text string) string {
	return normalizeReplacer.Replace(removeMarks(text))
}
