package text

import (
	"strings"
	"testing"
	"unicode"
)

func TestFinglish(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"salam":         {"سلام", "salam"},
		"amir":          {"امیر", "amir"},
		"two words":     {"سلام دنیا", "salam dania"},
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

// TestFinglishRules pins the individual rules, so a change to one of them shows
// up here as a named failure rather than as a shift in the accuracy score.
func TestFinglishRules(t *testing.T) {
	t.Parallel()

	tests := []struct{ rule, in, want string }{
		{"long vowel: و between consonants is u", "خوب", "khub"},
		{"long vowel: و between consonants is u", "دوست", "dust"},
		{"consonant: و before a vowel letter is v", "اهواز", "ahvaz"},
		{"consonant: و after a vowel is v", "دیوار", "divar"},
		{"long vowel: ی after a consonant is i", "شیر", "shir"},
		{"long vowel: ی after a consonant is i", "خیابان", "khiaban"},
		{"consonant: ی after a vowel is y", "چای", "chay"},
		{"silent vav in خوا", "خواهر", "khahar"},
		{"silent vav in خوا", "خواندن", "khandan"},
		{"final ه after a consonant is e", "خانه", "khane"},
		{"final ه after a vowel stays h", "ماه", "mah"},
		{"medial ه stays h", "شهر", "shahr"},
		{"initial ای is a long i", "ایران", "iran"},
		{"initial آ is a long a", "آب", "ab"},
		{"a closing pair may end a word", "مرد", "mard"},
		{"a closing pair may end a word", "دست", "dast"},
		{"r cannot close, so the word gains a syllable", "مادر", "madar"},
		{"n cannot close, so the word gains a syllable", "دیدن", "didan"},
		{"a long run splits into syllables", "رفتن", "raftan"},
		{"a long run splits into syllables", "مشهد", "mashhad"},
		{"arabic yeh and kaf are folded", "علي", "ali"},
	}

	for _, tt := range tests {
		t.Run(tt.rule+"/"+tt.in, func(t *testing.T) {
			t.Parallel()
			if got := Finglish(tt.in); got != tt.want {
				t.Errorf("Finglish(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestFinglishStructure checks the properties that must hold for any input,
// whatever the transliteration guesses. These are what separate a readable
// approximation from a pile of consonants.
func TestFinglishStructure(t *testing.T) {
	t.Parallel()

	words := append(finglishReference(), finglishHeldOut()...)
	for _, w := range words {
		got := Finglish(w.word)

		if got == "" {
			t.Errorf("Finglish(%q) is empty", w.word)
			continue
		}
		for _, r := range got {
			if r > unicode.MaxASCII {
				t.Errorf("Finglish(%q) = %q contains non-ASCII %q", w.word, got, r)
			}
		}
		// Every word must carry a vowel: a word of bare consonants is not
		// pronounceable, which is the failure the syllable pass exists to
		// prevent.
		if !strings.ContainsAny(got, "aeiou") {
			t.Errorf("Finglish(%q) = %q has no vowel", w.word, got)
		}
		// No run of four consonants: Persian never stacks that many, so such a
		// run means a syllable lost its nucleus.
		if run := longestConsonantRun(got); run > 3 {
			t.Errorf("Finglish(%q) = %q has a run of %d consonants", w.word, got, run)
		}
	}
}

// TestFinglishWordsAreIndependent checks that a word transliterates the same
// way on its own as it does in a sentence. The previous implementation carried
// its syllable state across spaces, so the same word could come out differently
// depending on what preceded it.
func TestFinglishWordsAreIndependent(t *testing.T) {
	t.Parallel()

	words := []string{"سلام", "کتاب", "مرد", "خانه", "ایران", "دوست"}
	for _, w := range words {
		alone := Finglish(w)
		for _, prefix := range []string{"سلام ", "کتاب ", "من به "} {
			sentence := Finglish(prefix + w)
			if got := sentence[strings.LastIndex(sentence, " ")+1:]; got != alone {
				t.Errorf("Finglish(%q) gives %q for %q, but %q on its own", prefix+w, got, w, alone)
			}
		}
	}
}

func TestFinglishSpacing(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"single space":   {"سلام دنیا", "salam dania"},
		"leading space":  {" سلام", " salam"},
		"trailing space": {"سلام ", "salam "},
		"double space":   {"سلام  دنیا", "salam  dania"},
		"punctuation":    {"سلام، دنیا", "salam dania"},
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

func TestFinglishAlphabetIsComplete(t *testing.T) {
	t.Parallel()

	// The 32 letters of the Persian alphabet, plus the hamza carriers and the
	// Arabic letters that leak into Persian text.
	const alphabet = "اآبپتثجچحخدذرزژسشصضطظعغفقکگلمنوهیءئيك"

	for _, r := range alphabet {
		if !isFinglishLetter(r) {
			t.Errorf("Finglish does not recognize %q", string(r))
		}
		if got := Finglish(string(r)); got == "" && r != 'ء' && r != 'ئ' {
			t.Errorf("Finglish(%q) is empty", string(r))
		}
	}
}

// longestConsonantRun returns the length of the longest run of Latin
// consonants, counting the digraphs kh, sh, ch, zh and gh as one sound each.
func longestConsonantRun(s string) int {
	longest, run := 0, 0
	for i := 0; i < len(s); {
		width := 1
		if i+1 < len(s) {
			switch s[i : i+2] {
			case "kh", "sh", "ch", "zh", "gh":
				width = 2
			}
		}
		if strings.ContainsRune("aeiou", rune(s[i])) || s[i] == ' ' {
			run = 0
		} else {
			run++
			if run > longest {
				longest = run
			}
		}
		i += width
	}
	return longest
}

// TestFinglishLetterVariants covers the letters that appear rarely but must
// still be read correctly: the initial او pair, the hamza carriers and the
// Arabic spellings of alef.
func TestFinglishLetterVariants(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"initial او is a long u": {"او", "u"},
		"او before a consonant":  {"اوج", "uj"},
		"arabic alef with hamza": {"أمير", "amir"},
		"arabic alef with kasra": {"إمير", "amir"},
		"hamza carrier mid-word": {"مسئول", "masaul"},
		"standalone hamza":       {"جزء", "jaza"},
		"ain carries a vowel":    {"سعید", "said"},
		"initial ain":            {"عمر", "amar"},
		"arabic kaf":             {"كتاب", "katab"},
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
