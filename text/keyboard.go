package text

import "github.com/amiranmanesh/go-persian-tools/internal/runeutil"

// persianKeyboardLayout pairs each key of a US QWERTY keyboard with the
// character the same key produces on the standard Iranian (ISIRI) layout.
// Both conversion directions are derived from this single table.
var persianKeyboardLayout = [][2]string{
	// Top letter row.
	{"q", "ض"},
	{"w", "ص"},
	{"e", "ث"},
	{"r", "ق"},
	{"t", "ف"},
	{"y", "غ"},
	{"u", "ع"},
	{"i", "ه"},
	{"o", "خ"},
	{"p", "ح"},
	{"[", "ج"},
	{"]", "چ"},
	// Home row.
	{"a", "ش"},
	{"s", "س"},
	{"d", "ی"},
	{"f", "ب"},
	{"g", "ل"},
	{"h", "ا"},
	{"j", "ت"},
	{"k", "ن"},
	{"l", "م"},
	{";", "ک"},
	{"'", "گ"},
	// Bottom row.
	{"z", "ظ"},
	{"x", "ط"},
	{"c", "ز"},
	{"v", "ر"},
	{"b", "ذ"},
	{"n", "د"},
	{"m", "پ"},
	{",", "و"},
	// Punctuation.
	{"?", "؟"},
}

var (
	persianKeyReplacer = runeutil.Replacer(persianKeyboardLayout, false)
	englishKeyReplacer = runeutil.Replacer(persianKeyboardLayout, true)
)

// SwitchToPersianKey rewrites text as if it had been typed on a Persian
// keyboard: every character is replaced by the one sharing its physical key.
// It is the usual fix for "sghl" when the writer meant "سلام".
//
// Characters outside the layout, including uppercase letters and digits, are
// left as they are.
func SwitchToPersianKey(text string) string {
	return persianKeyReplacer.Replace(text)
}

// SwitchToEnglishKey is the inverse of [SwitchToPersianKey]: it rewrites text
// as if it had been typed on an English keyboard, turning "اثغ" back into
// "hey".
func SwitchToEnglishKey(text string) string {
	return englishKeyReplacer.Replace(text)
}
