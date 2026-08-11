package text

import "strings"

// Finglish transliteration works in three passes, because Persian spelling and
// Latin spelling disagree about what a word is made of:
//
//  1. The text is split into words. Every word is transliterated on its own, so
//     nothing leaks across a space.
//  2. Each word's letters become a list of phonemes — consonants and vowels.
//     This is where the letters that can be either (ا و ی ه ع) are resolved
//     from their position, and where the silent letters are dropped.
//  3. The phonemes are grouped into syllables. Persian does not write its short
//     vowels, so most words arrive with consonants that have no vowel between
//     them; this pass supplies one wherever a syllable needs a nucleus.
//
// Splitting the work this way is what makes the vowel rules expressible at all:
// "و is a vowel between two consonants" is a statement about phonemes, not
// about letters, and "every syllable needs a nucleus" is a statement about
// syllables, not about either.

// phoneme is one sound: a Latin spelling plus whether it can be a syllable
// nucleus.
type phoneme struct {
	text  string
	vowel bool
}

// finglishConsonants maps each Persian letter to its Latin consonant. Letters
// that share a sound in modern Persian (the four z's, the three s's) map to the
// same Latin letter on purpose.
var finglishConsonants = map[rune]string{
	'ب': "b",
	'پ': "p",
	'ت': "t",
	'ث': "s",
	'ج': "j",
	'چ': "ch",
	'ح': "h",
	'خ': "kh",
	'د': "d",
	'ذ': "z",
	'ر': "r",
	'ز': "z",
	'ژ': "zh",
	'س': "s",
	'ش': "sh",
	'ص': "s",
	'ض': "z",
	'ط': "t",
	'ظ': "z",
	'غ': "gh",
	'ف': "f",
	'ق': "gh",
	'ک': "k",
	'گ': "g",
	'ل': "l",
	'م': "m",
	'ن': "n",
	'ه': "h",
	'و': "v",
	'ی': "y",
}

// The default short vowel. Persian leaves its short vowels unwritten, so a word
// like "مرد" carries no clue that it is "mard" rather than "merd" or "mord".
// Fatha is the most common of the three, so it is what an unwritten vowel
// becomes.
const finglishDefaultVowel = "a"

// isFinglishLetter reports whether r is a letter this transliterator handles.
func isFinglishLetter(r rune) bool {
	switch r {
	case 'ا', 'آ', 'أ', 'إ', 'ئ', 'ء', 'ؤ', 'ع', 'ی', 'ي', 'ى', 'و', 'ک', 'ك':
		return true
	}
	_, ok := finglishConsonants[r]
	return ok
}

// Finglish transliterates Persian text into the Latin script, the informal
// romanization Iranians call "Finglish": "سلام" becomes "salam".
//
// Spaces are kept as word boundaries and every word is transliterated
// independently. Characters outside the Persian alphabet are dropped.
//
// The mapping is a heuristic, and it is worth knowing where it guesses.
// Persian does not write its short vowels, so nothing on the page distinguishes
// "gol" from "gal"; this function supplies the most common one, an "a". It
// therefore writes "mard", "salam" and "khub" correctly but renders "گل" as
// "gal" rather than "gol" and "پدر" as "padar" rather than "pedar", and
// homographs such as "کرم" (kerm, karam, kerem) have no single right answer at
// all.
//
// What it does get right is the structure: long vowels, silent letters, the
// consonant-or-vowel letters resolved from context, and a vowel in every
// syllable. Expect a readable approximation, not a reversible encoding — and
// see [TestFinglishAccuracy] in the tests for the measured hit rate against a
// reference word list.
func Finglish(text string) string {
	var b strings.Builder
	b.Grow(len(text))

	runes := []rune(text)
	for i := 0; i < len(runes); {
		if !isFinglishLetter(runes[i]) {
			// Spaces separate words; everything else is dropped.
			if runes[i] == ' ' {
				b.WriteByte(' ')
			}
			i++
			continue
		}

		// Take the whole word, so its letters can be read in context.
		start := i
		for i < len(runes) && isFinglishLetter(runes[i]) {
			i++
		}
		b.WriteString(finglishWord(runes[start:i]))
	}

	return b.String()
}

// finglishWord transliterates a single word.
func finglishWord(word []rune) string {
	return finglishSyllables(finglishPhonemes(word))
}

// finglishPhonemes reads a word's letters into phonemes, resolving the letters
// whose sound depends on where they sit.
func finglishPhonemes(word []rune) []phoneme {
	out := make([]phoneme, 0, len(word)+2)

	// lastWasVowel tracks the previous phoneme, which is what decides whether
	// و and ی are consonants or long vowels.
	lastWasVowel := false
	emit := func(text string, vowel bool) {
		out = append(out, phoneme{text: text, vowel: vowel})
		lastWasVowel = vowel
	}

	for i := 0; i < len(word); i++ {
		r := normalizeFinglishLetter(word[i])
		atStart := len(out) == 0
		atEnd := i == len(word)-1

		// "خوا" is written with a vav that is not pronounced: خواهر is
		// "khahar", not "khavahar". The alef after it carries the vowel.
		if r == 'خ' && i+2 < len(word) &&
			normalizeFinglishLetter(word[i+1]) == 'و' &&
			normalizeFinglishLetter(word[i+2]) == 'ا' {
			emit("kh", false)
			emit("a", true)
			i += 2
			continue
		}

		switch r {
		case 'آ':
			emit("a", true)

		case 'ا':
			switch {
			case !atStart:
				// Mid-word alef is always the long vowel.
				emit("a", true)
			case i+1 < len(word) && normalizeFinglishLetter(word[i+1]) == 'ی':
				// ایران: the pair opens the word as a long "i".
				emit("i", true)
				i++
			case i+1 < len(word) && normalizeFinglishLetter(word[i+1]) == 'و':
				// او: the pair opens the word as a long "u".
				emit("u", true)
				i++
			default:
				emit("a", true)
			}

		case 'و':
			// A vav is the consonant "v" when it opens a syllable — at the
			// start of a word, straight after a vowel, or straight before one,
			// as in قزوین and اهواز — and the long vowel "u" when it sits
			// between two consonants, as in خوب and دوست.
			if atStart || lastWasVowel || finglishVowelFollows(word, i) {
				emit("v", false)
			} else {
				emit("u", true)
			}

		case 'ی':
			// A yeh is the consonant "y" when it opens a syllable — at the
			// start of a word, or straight after a vowel, as in "چای" — and
			// the long vowel "i" after a consonant, as in "شیر" and "خیابان".
			if atStart || lastWasVowel {
				emit("y", false)
			} else {
				emit("i", true)
			}

		case 'ه':
			// A final heh after a consonant is the vowel of the last syllable:
			// خانه is "khane". Anywhere else it is the consonant "h".
			if atEnd && !atStart && !lastWasVowel {
				emit("e", true)
			} else {
				emit("h", false)
			}

		case 'ع', 'ء', 'ئ', 'ؤ':
			// The glottal letters have no Latin consonant of their own: they
			// carry a vowel and are otherwise silent.
			if !lastWasVowel {
				emit(finglishDefaultVowel, true)
			}
			// A carrier does not stand between a following و or ی and the
			// consonant before it, so it must not make them look like onsets:
			// سعید is "said", not "sayd".
			lastWasVowel = false

		default:
			if c, ok := finglishConsonants[r]; ok {
				emit(c, false)
			}
		}
	}

	return out
}

// normalizeFinglishLetter folds the Arabic spellings onto their Persian
// counterparts so the rules only have to handle one shape of each letter.
func normalizeFinglishLetter(r rune) rune {
	switch r {
	case 'ي', 'ى':
		return 'ی'
	case 'ك':
		return 'ک'
	case 'أ', 'إ':
		return 'ا'
	}
	return r
}

// finglishVowelFollows reports whether the letter after position i is one that
// is always read as a vowel, which makes a vav or yeh at i an onset consonant
// rather than a long vowel of its own.
func finglishVowelFollows(word []rune, i int) bool {
	if i+1 >= len(word) {
		return false
	}
	switch normalizeFinglishLetter(word[i+1]) {
	case 'ا', 'آ', 'ی':
		return true
	}
	return false
}

// finglishConsonantRun counts the consonants starting at i, stopping at the
// next vowel or at the end of the word.
func finglishConsonantRun(phonemes []phoneme, i int) int {
	n := 0
	for ; i < len(phonemes) && !phonemes[i].vowel; i++ {
		n++
	}
	return n
}

// finglishClosesSyllable reports whether the pair first+second may close a
// syllable together.
//
// Persian is content to end a syllable on two consonants when the second is an
// obstruent or an "m" — mard, dast, sakht, garm, cheshm, sobh — but an "r" or
// an "n" there pulls the consonant before it into a new syllable, which is why
// پدر is "pa-dar" and دیدن is "di-dan" rather than "padr" and "didn".
//
// An "h" after a short vowel is the exception: it is weak enough to close
// against anything, so شهر stays "shahr" rather than becoming "shahar". After a
// written vowel even "h" gives way, which is what separates خواهر ("kha-har")
// from شهر ("shahr") — the two carry the same pair of consonants and differ
// only in the vowel before it.
func finglishClosesSyllable(first, second string, longVowel bool) bool {
	switch second {
	case "r", "n", "y":
		return first == "h" && !longVowel
	}
	return true
}

// finglishSyllables groups phonemes into syllables and writes them out.
//
// A Persian syllable is an onset consonant, a vowel, and up to two closing
// consonants. Since the short vowels are not written, most words arrive as a
// run of consonants with no nucleus at all; this is where the missing vowels
// are supplied, one per syllable rather than one per gap.
func finglishSyllables(phonemes []phoneme) string {
	var b strings.Builder

	for i := 0; i < len(phonemes); {
		// Onset: one consonant, if there is one. A word may also open on a
		// vowel, in which case the syllable has no onset.
		if !phonemes[i].vowel {
			b.WriteString(phonemes[i].text)
			i++
		}

		// Nucleus: the written vowel, or the default short vowel when the
		// spelling leaves it out. Which of the two it was decides how much the
		// syllable may keep as a coda below.
		long := false
		if i < len(phonemes) && phonemes[i].vowel {
			b.WriteString(phonemes[i].text)
			long = true
			i++
		} else {
			b.WriteString(finglishDefaultVowel)
		}

		// Coda: the closing consonants. How many this syllable may keep depends
		// on what is left of the word, because every later syllable needs an
		// onset of its own:
		//
		//   run of 1  the word ends here, so the consonant closes it: "gal"
		//   run of 2  both close the last syllable: "mard"
		//   run of 3+ one closes this syllable and the rest start the next
		//             one, which is what makes "رفتن" read "raf-tan"
		//
		// A consonant with a vowel after it always belongs to the next
		// syllable, so the run stops there.
		run := finglishConsonantRun(phonemes, i)
		endsWord := i+run == len(phonemes)

		var keep int
		switch {
		case !endsWord:
			// A vowel follows the run, so its last consonant is the onset of
			// the next syllable and cannot be taken here.
			if run >= 2 {
				keep = 1
			}
		case run <= 1:
			keep = run
		case run == 2:
			keep = 2
			if !finglishClosesSyllable(phonemes[i].text, phonemes[i+1].text, long) {
				// The pair cannot close the word, so the second consonant opens
				// another syllable instead: "مادر" is "ma-dar", not "madr".
				keep = 0
			}
		default:
			// A longer run splits into further syllables, one consonant
			// closing this one and the rest opening the next: "رفتن" is
			// "raf-tan".
			keep = 1
		}
		for ; keep > 0; keep-- {
			b.WriteString(phonemes[i].text)
			i++
		}
	}

	return b.String()
}
