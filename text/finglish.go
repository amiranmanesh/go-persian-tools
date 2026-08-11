package text

import "strings"

// finglishSound is the Latin sound of one Persian letter: the consonant it
// always contributes, plus the vowel it turns into when it lands where a vowel
// is expected. An empty vowel marks a pure consonant.
type finglishSound struct {
	consonant string
	vowel     string
}

// finglishAlphabet maps each Persian letter to its Latin sound. Letters that
// share a sound in modern Persian (the four z's, the three s's) deliberately
// map to the same Latin letter.
var finglishAlphabet = map[rune]finglishSound{
	'ا': {"a", "a"},
	'آ': {"a", "a"},
	'ئ': {"a", ""},
	'ء': {"a", ""},
	'ب': {"b", ""},
	'پ': {"p", ""},
	'ت': {"t", ""},
	'ث': {"s", ""},
	'ج': {"j", ""},
	'چ': {"ch", ""},
	'ح': {"h", ""},
	'خ': {"kh", ""},
	'د': {"d", ""},
	'ذ': {"z", ""},
	'ر': {"r", ""},
	'ز': {"z", ""},
	'ژ': {"zh", ""},
	'س': {"s", ""},
	'ش': {"sh", ""},
	'ص': {"s", ""},
	'ض': {"z", ""},
	'ط': {"t", ""},
	'ظ': {"z", ""},
	'ع': {"", "a"},
	'غ': {"gh", ""},
	'ف': {"f", ""},
	'ق': {"gh", ""},
	'ک': {"k", ""},
	'گ': {"g", ""},
	'ل': {"l", ""},
	'م': {"m", ""},
	'ن': {"n", ""},
	'و': {"v", "o"},
	'ه': {"h", ""},
	'ی': {"y", "i"},
	'ي': {"y", "i"},
}

// Positions in the syllable pattern the transliterator walks through. Persian
// omits short vowels in writing, so the vowel of each syllable has to be
// guessed from where the letter falls rather than read off the page.
const (
	finglishOnset      = iota // start of a word; emit the consonant
	finglishFirstVowel        // the slot right after it takes a vowel
	finglishConsonant         // a consonant is due
	finglishVowelSlot         // a vowel is due, and every state past this one
)

// Finglish transliterates Persian text into the Latin script, the informal
// romanization Iranians call "Finglish": "سلام" becomes "salam".
//
// The mapping is a heuristic. Persian is written without short vowels, so they
// are inferred from each letter's position in the word, and homographs such as
// "کرم" (kerm, karam, kerem) have no single correct answer. Expect a readable
// approximation, not a reversible encoding. Characters outside the Persian
// alphabet are dropped, except spaces, which are kept as word boundaries.
func Finglish(text string) string {
	var b strings.Builder
	b.Grow(len(text))

	state := finglishOnset
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == ' ' {
			b.WriteByte(' ')
		}

		sound, ok := finglishAlphabet[r]
		if !ok {
			continue
		}

		switch state {
		case finglishOnset:
			b.WriteString(sound.consonant)
			state = finglishFirstVowel

		case finglishFirstVowel:
			if sound.vowel != "" {
				b.WriteString(sound.vowel)
			} else {
				// A consonant showed up where a vowel was due: supply the
				// default short vowel and read this letter again as the
				// consonant of the next syllable.
				b.WriteString("a")
				i--
			}
			state = finglishConsonant

		case finglishConsonant:
			b.WriteString(sound.consonant)
			state = finglishVowelSlot

		default: // finglishVowelSlot and beyond: a run of consonants
			if sound.vowel != "" {
				b.WriteString(sound.vowel)
				state = finglishConsonant
			} else {
				b.WriteString(sound.consonant)
				state++
			}
		}
	}

	// Collapse the vowel clusters the guessing above can produce.
	out := b.String()
	out = strings.ReplaceAll(out, "aa", "a")
	out = strings.ReplaceAll(out, "ao", "o")
	out = strings.ReplaceAll(out, "ae", "e")
	return out
}
