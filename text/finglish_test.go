package text

import "testing"

func TestFinglish(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"salam":         {"سلام", "salam"},
		"amir":          {"امیر", "amir"},
		"tehran":        {"تهران", "tahran"},
		"ketab":         {"کتاب", "katab"},
		"two words":     {"سلام دنیا", "salam dnia"},
		"drops unknown": {"سلام!", "salam"},
		"keeps spaces":  {"ا ا", "a a"},
		"empty":         {"", ""},
		"only unknown":  {"123", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := Finglish(tt.in); got != tt.want {
				t.Errorf("Finglish(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestFinglishIsAscii checks that every letter in the alphabet transliterates
// to ASCII, whatever the surrounding context.
func TestFinglishIsAscii(t *testing.T) {
	t.Parallel()

	for r := range finglishAlphabet {
		for _, out := range []string{Finglish(string(r)), Finglish("ا" + string(r) + "ا")} {
			for _, c := range out {
				if c > 127 {
					t.Errorf("Finglish of %q produced non-ASCII %q", string(r), string(c))
				}
			}
		}
	}
}

func TestFinglishAlphabetIsComplete(t *testing.T) {
	t.Parallel()

	// The 32 letters of the Persian alphabet, plus the hamza carriers and the
	// Arabic yeh that leaks into Persian text.
	const alphabet = "اآبپتثجچحخدذرزژسشصضطظعغفقکگلمنوهیءئي"

	for _, r := range alphabet {
		if _, ok := finglishAlphabet[r]; !ok {
			t.Errorf("finglishAlphabet is missing %q", string(r))
		}
	}
}
