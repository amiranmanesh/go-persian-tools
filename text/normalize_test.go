package text

import "testing"

func TestFixArabic(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"arabic yeh":      {"علي\u200Cرضا", "علی\u200Cرضا"},
		"arabic kaf":      {"كتاب", "کتاب"},
		"alef maksura":    {"مصطفى", "مصطفی"},
		"already persian": {"علی", "علی"},
		"keeps zwnj":      {"می\u200Cروم", "می\u200Cروم"},
		"strips harakat":  {"سَلامِ گَرم", "سلام گرم"},
		"strips tatweel":  {"سـلام", "سلام"},
		"keeps madda":     {"آب", "آب"},
		"empty":           {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := FixArabic(tt.in); got != tt.want {
				t.Errorf("FixArabic(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"mixed sentence": {
			"متن نا\u200Cنرمال عربي و فارسی با عددهای ۱ و 1",
			"متن نا نرمال عربی و فارسی با عددهای ۱ و 1",
		},
		"alef variants":  {"آب أمير ٱحمد", "اب امیر احمد"},
		"teh marbuta":    {"فاطمة", "فاطمه"},
		"zwnj to space":  {"می\u200Cروم", "می روم"},
		"strips harakat": {"سَلامِ گَرم", "سلام گرم"},
		"repeated mark":  {"بِِ", "ب"},
		"ezafe hamza":    {"مقالهٔ من", "مقاله من"},
		"empty":          {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := Normalize(tt.in); got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestIsHarakat(t *testing.T) {
	t.Parallel()

	marks := map[string]rune{
		"fathatan":         0x064B,
		"fatha":            0x064E,
		"kasra":            0x0650,
		"shadda":           0x0651,
		"sukun":            0x0652,
		"superscript alef": 0x0670,
		"tatweel":          0x0640,
	}
	for name, r := range marks {
		if !isHarakat(r) {
			t.Errorf("isHarakat(%s, U+%04X) = false, want true", name, r)
		}
	}

	// The marks that change letter identity are not harakat: FixArabic keeps
	// them, and only Normalize drops them.
	for name, r := range map[string]rune{
		"madda above": 0x0653,
		"hamza above": 0x0654,
		"hamza below": 0x0655,
	} {
		if isHarakat(r) {
			t.Errorf("isHarakat(%s, U+%04X) = true, want false", name, r)
		}
		if !isCombiningMadda(r) {
			t.Errorf("isCombiningMadda(%s, U+%04X) = false, want true", name, r)
		}
	}

	// Letters and the marks that change letter identity stay.
	keep := map[string]rune{
		"alef":          'ا',
		"persian yeh":   'ی',
		"madda above":   0x0653,
		"hamza above":   0x0654,
		"hamza below":   0x0655,
		"latin a":       'a',
		"persian digit": '۱',
	}
	for name, r := range keep {
		if isHarakat(r) {
			t.Errorf("isHarakat(%s, U+%04X) = true, want false", name, r)
		}
	}
}

// TestNormalizeIsIdempotent guards the property callers rely on when they
// normalize once at write time and again at query time.
func TestNormalizeIsIdempotent(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"متن نا\u200Cنرمال عربي و فارسی",
		"آب أمير ٱحمد",
		"فاطمة",
		"salam 123",
	}
	for _, in := range inputs {
		once := Normalize(in)
		if twice := Normalize(once); twice != once {
			t.Errorf("Normalize is not idempotent for %q: %q then %q", in, once, twice)
		}
	}
}
