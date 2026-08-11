package digit

import (
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

var fuzzSeeds = []string{
	"",
	"سلام",
	"123salam456",
	"۱۲۳",
	"٠١٢٣٤٥٦٧٨٩",
	"مبلغ ۱۲۳۴۵۶۷ ریال",
	"1,234,567",
	"۱،۲۳۴،۵۶۷",
	"-42",
	"\x00\xff",
}

func addFuzzSeeds(f *testing.F) {
	f.Helper()
	for _, seed := range fuzzSeeds {
		f.Add(seed)
	}
}

// FuzzToEnglishDigits checks that no Persian or Arabic-Indic digit survives the
// conversion, and that nothing else about the string changes length in runes.
func FuzzToEnglishDigits(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		out := ToEnglishDigits(in)

		if strings.ContainsAny(out, PersianDigits+ArabicDigits) {
			t.Fatalf("ToEnglishDigits(%q) = %q still contains non-ASCII digits", in, out)
		}
		if got, want := utf8.RuneCountInString(out), utf8.RuneCountInString(in); got != want {
			t.Fatalf("ToEnglishDigits(%q) changed the rune count: %d, want %d", in, got, want)
		}
	})
}

// FuzzDigitRoundTrip checks that converting to Persian digits and back is
// lossless for any input.
func FuzzDigitRoundTrip(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		// Start from a string with no Persian or Arabic digits, so the round
		// trip is expected to be exact.
		base := ToEnglishDigits(in)
		if got := ToEnglishDigits(ToPersianDigits(base)); got != base {
			t.Fatalf("round trip of %q produced %q", base, got)
		}
	})
}

// FuzzOnlyNumbers checks that the filters only ever remove runes.
func FuzzOnlyNumbers(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		for name, fn := range map[string]func(string) string{
			"OnlyEnglishNumbers": OnlyEnglishNumbers,
			"OnlyPersianNumbers": OnlyPersianNumbers,
			"OnlyNumbers":        OnlyNumbers,
		} {
			out := fn(in)
			if !isSubsequence(out, in) {
				t.Fatalf("%s(%q) = %q is not a subsequence of the input", name, in, out)
			}
		}
	})
}

// FuzzCurrency checks that formatting an amount never invents or loses digits.
func FuzzCurrency(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		out := Currency(in)

		want := ToPersianDigits(OnlyNumbers(in))
		got := strings.ReplaceAll(out, ThousandsSeparator, "")
		if got != want {
			t.Fatalf("Currency(%q) = %q; digits %q, want %q", in, out, got, want)
		}
	})
}

// FuzzCommasRoundTrip checks that grouping a number and parsing it back is
// lossless, for any int64.
func FuzzCommasRoundTrip(f *testing.F) {
	for _, seed := range []int64{0, 1, -1, 999, 1000, -1234567, 1 << 62} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, n int64) {
		grouped := AddCommas(n)

		if stripped := strings.ReplaceAll(grouped, GroupSeparator, ""); stripped != strconv.FormatInt(n, 10) {
			t.Fatalf("AddCommas(%d) = %q, whose digits are %q", n, grouped, stripped)
		}
		back, err := RemoveCommas(grouped)
		if err != nil {
			t.Fatalf("RemoveCommas(%q): %v", grouped, err)
		}
		if back != n {
			t.Fatalf("round trip of %d produced %d", n, back)
		}
	})
}

// FuzzToWords checks that every int64 gets a non-empty spelling that never
// contains a digit, and that DigitToWord agrees with it.
func FuzzToWords(f *testing.F) {
	for _, seed := range []int64{0, 1, -1, 10, 19, 100, 1000, -156788, 1 << 62} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, n int64) {
		words := ToWords(n)
		if words == "" {
			t.Fatalf("ToWords(%d) is empty", n)
		}
		if strings.ContainsAny(words, EnglishDigits+PersianDigits) {
			t.Fatalf("ToWords(%d) = %q contains a digit", n, words)
		}
		if got := DigitToWord(strconv.FormatInt(n, 10)); got != words {
			t.Fatalf("DigitToWord disagrees with ToWords for %d: %q vs %q", n, got, words)
		}
	})
}

// isSubsequence reports whether the runes of sub appear in s in the same order.
func isSubsequence(sub, s string) bool {
	rest := s
	for _, r := range sub {
		i := strings.IndexRune(rest, r)
		if i < 0 {
			return false
		}
		rest = rest[i+len(string(r)):]
	}
	return true
}
