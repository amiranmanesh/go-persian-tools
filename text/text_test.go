package text

import "testing"

func TestReverse(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"ascii":       {"hello", "olleh"},
		"persian":     {"سلام", "مالس"},
		"single rune": {"ا", "ا"},
		"empty":       {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := Reverse(tt.in); got != tt.want {
				t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if got := Reverse(Reverse(tt.in)); got != tt.in {
				t.Errorf("Reverse is not its own inverse for %q: got %q", tt.in, got)
			}
		})
	}
}

func TestCheckIsEnglish(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want bool
	}{
		"lowercase":  {"ali", true},
		"mixed case": {"Ali", true},
		"persian":    {"علي", false},
		"with digit": {"ali1", false},
		"with space": {"ali reza", false},
		"empty":      {"", true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := CheckIsEnglish(tt.in); got != tt.want {
				t.Errorf("CheckIsEnglish(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestZeroWidthConstants(t *testing.T) {
	t.Parallel()

	if ZWNJ != "\u200C" {
		t.Errorf("ZWNJ = %q, want U+200C", ZWNJ)
	}
	if ZWJ != "\u200D" {
		t.Errorf("ZWJ = %q, want U+200D", ZWJ)
	}
}

func TestOnlyPersianAlpha(t *testing.T) {
	tests := map[string]struct{ in, want string }{
		"drops latin and digits": {"123456شاهینshaahin", "شاهین"},
		"keeps the dot":          {"علی.رضا", "علی.رضا"},
		"empty":                  {"", ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := OnlyPersianAlpha(tt.in); got != tt.want {
				t.Errorf("OnlyPersianAlpha(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
