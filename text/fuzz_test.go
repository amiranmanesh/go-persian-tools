package text

import (
	"strings"
	"testing"
)

var fuzzSeeds = []string{
	"",
	"سلام",
	"123salam456",
	"علي‌رضا",
	"sghl o,fd ?",
	"مقالهٔ من",
	"\x00\xff",
}

func addFuzzSeeds(f *testing.F) {
	f.Helper()
	for _, seed := range fuzzSeeds {
		f.Add(seed)
	}
}

// FuzzNormalize checks that normalizing twice is the same as normalizing once,
// which is what lets callers normalize at write time and at query time.
func FuzzNormalize(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		once := Normalize(in)
		if twice := Normalize(once); twice != once {
			t.Fatalf("Normalize(%q) is not idempotent: %q then %q", in, once, twice)
		}
	})
}

// FuzzFixArabic checks that FixArabic is idempotent too.
func FuzzFixArabic(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		once := FixArabic(in)
		if twice := FixArabic(once); twice != once {
			t.Fatalf("FixArabic(%q) is not idempotent: %q then %q", in, once, twice)
		}
	})
}

// FuzzOnlyPersianAlpha checks that the filter only ever removes runes.
func FuzzOnlyPersianAlpha(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		out := OnlyPersianAlpha(in)
		if !isSubsequence(out, in) {
			t.Fatalf("OnlyPersianAlpha(%q) = %q is not a subsequence of the input", in, out)
		}
	})
}

// FuzzReverse checks that reversing twice returns the original runes.
func FuzzReverse(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		if got := Reverse(Reverse(in)); got != string([]rune(in)) {
			t.Fatalf("Reverse is not an involution for %q: got %q", in, got)
		}
	})
}

// FuzzKeyboardRoundTrip checks that a layout switch and its inverse cancel out
// for text made only of characters the layout covers.
func FuzzKeyboardRoundTrip(f *testing.F) {
	addFuzzSeeds(f)
	f.Fuzz(func(t *testing.T, in string) {
		// Restrict the input to keys the layout actually maps, so the round
		// trip is expected to be exact.
		var b strings.Builder
		for _, pair := range persianKeyboardLayout {
			if strings.Contains(in, pair[0]) {
				b.WriteString(pair[0])
			}
		}
		base := b.String()
		if got := SwitchToEnglishKey(SwitchToPersianKey(base)); got != base {
			t.Fatalf("round trip of %q produced %q", base, got)
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
