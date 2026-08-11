package text

import "testing"

func TestSwitchToPersianKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"sentence":         {"sghl o,fd ? o,fl llk,k", "سلام خوبی ؟ خوبم ممنون"},
		"brackets":         {"[]", "جچ"},
		"unmapped is kept": {"A1", "A1"},
		"empty":            {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := SwitchToPersianKey(tt.in); got != tt.want {
				t.Errorf("SwitchToPersianKey(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSwitchToEnglishKey(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"sentence": {"اثغ صاشفس عح ؟", "hey whats up ?"},
		"brackets": {"جچ", "[]"},
		"empty":    {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := SwitchToEnglishKey(tt.in); got != tt.want {
				t.Errorf("SwitchToEnglishKey(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestKeyboardRoundTrip checks that the layout table is a bijection: switching
// to Persian and back has to return the original text.
func TestKeyboardRoundTrip(t *testing.T) {
	t.Parallel()

	for _, pair := range persianKeyboardLayout {
		english, persian := pair[0], pair[1]

		if got := SwitchToPersianKey(english); got != persian {
			t.Errorf("SwitchToPersianKey(%q) = %q, want %q", english, got, persian)
		}
		if got := SwitchToEnglishKey(persian); got != english {
			t.Errorf("SwitchToEnglishKey(%q) = %q, want %q", persian, got, english)
		}
	}
}

func TestKeyboardLayoutHasNoDuplicates(t *testing.T) {
	t.Parallel()

	seenEnglish := make(map[string]bool, len(persianKeyboardLayout))
	seenPersian := make(map[string]bool, len(persianKeyboardLayout))
	for _, pair := range persianKeyboardLayout {
		if seenEnglish[pair[0]] {
			t.Errorf("duplicate English key %q", pair[0])
		}
		if seenPersian[pair[1]] {
			t.Errorf("duplicate Persian key %q", pair[1])
		}
		seenEnglish[pair[0]] = true
		seenPersian[pair[1]] = true
	}
}
