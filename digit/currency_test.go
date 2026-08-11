package digit

import "testing"

func TestCurrency(t *testing.T) {
	t.Parallel()

	tests := map[string]struct{ in, want string }{
		"english digits":  {"1234567", "۱،۲۳۴،۵۶۷"},
		"mixed digits":    {"123۴۵۶7", "۱،۲۳۴،۵۶۷"},
		"already grouped": {"1,234,567", "۱،۲۳۴،۵۶۷"},
		"exactly three":   {"123", "۱۲۳"},
		"one group plus":  {"1234", "۱،۲۳۴"},
		"strips letters":  {"مبلغ 5000 ریال", "۵،۰۰۰"},
		"no digits":       {"سلام", ""},
		"empty":           {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := Currency(tt.in); got != tt.want {
				t.Errorf("Currency(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestTomanAndRial(t *testing.T) {
	t.Parallel()

	const amount = "123۴۵۶7"

	if got, want := Toman(amount), "۱،۲۳۴،۵۶۷ تومان"; got != want {
		t.Errorf("Toman(%q) = %q, want %q", amount, got, want)
	}
	if got, want := Rial(amount), "۱،۲۳۴،۵۶۷ ﷼"; got != want {
		t.Errorf("Rial(%q) = %q, want %q", amount, got, want)
	}
}

// TestCurrencyGroupsFromTheRight checks the separator placement for every
// length up to four full groups.
func TestCurrencyGroupsFromTheRight(t *testing.T) {
	t.Parallel()

	tests := []struct{ in, want string }{
		{"1", "۱"},
		{"12", "۱۲"},
		{"123", "۱۲۳"},
		{"1234", "۱،۲۳۴"},
		{"12345", "۱۲،۳۴۵"},
		{"123456", "۱۲۳،۴۵۶"},
		{"1234567", "۱،۲۳۴،۵۶۷"},
		{"123456789012", "۱۲۳،۴۵۶،۷۸۹،۰۱۲"},
	}
	for _, tt := range tests {
		if got := Currency(tt.in); got != tt.want {
			t.Errorf("Currency(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
