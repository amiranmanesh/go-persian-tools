package digit

import "testing"

func TestToPersianDigits(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"mixed text":      {"123salam456", "۱۲۳salam۴۵۶"},
		"digits only":     {"0123456789", "۰۱۲۳۴۵۶۷۸۹"},
		"no digits":       {"سلام", "سلام"},
		"already Persian": {"۱۲۳", "۱۲۳"},
		"empty":           {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := ToPersianDigits(tt.in); got != tt.want {
				t.Errorf("ToPersianDigits(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToPersianDigitsFromInt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   int
		want string
	}{
		"positive": {123, "۱۲۳"},
		"zero":     {0, "۰"},
		"negative": {-42, "-۴۲"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := ToPersianDigitsFromInt(tt.in); got != tt.want {
				t.Errorf("ToPersianDigitsFromInt(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToEnglishDigits(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"mixed text":    {"۱۲۳salam۴۵۶", "123salam456"},
		"persian only":  {"۰۱۲۳۴۵۶۷۸۹", "0123456789"},
		"arabic indic":  {"٠١٢٣٤٥٦٧٨٩", "0123456789"},
		"both variants": {"۱٢۳", "123"},
		"empty":         {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := ToEnglishDigits(tt.in); got != tt.want {
				t.Errorf("ToEnglishDigits(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestOnly(t *testing.T) {
	t.Parallel()

	const mixed = "salam123hello456۱۲۳سلام۴۵۶"

	tests := map[string]struct {
		fn       func(string) string
		in, want string
	}{
		"english numbers":          {OnlyEnglishNumbers, mixed, "123456"},
		"english numbers with dot": {OnlyEnglishNumbers, "price: 12.50$", "12.50"},
		"persian numbers":          {OnlyPersianNumbers, mixed, "۱۲۳۴۵۶"},
		"all numbers":              {OnlyNumbers, mixed, "123456۱۲۳۴۵۶"},
		"nothing to keep":          {OnlyPersianNumbers, "abc", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := tt.fn(tt.in); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDigitSetsAreAligned(t *testing.T) {
	t.Parallel()

	sets := map[string]string{
		"EnglishDigits": EnglishDigits,
		"PersianDigits": PersianDigits,
		"ArabicDigits":  ArabicDigits,
	}
	for name, set := range sets {
		if n := len([]rune(set)); n != 10 {
			t.Errorf("%s has %d runes, want 10", name, n)
		}
	}
}
