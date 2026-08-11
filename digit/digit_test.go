package digit

import "testing"

func TestToWords(t *testing.T) {
	tests := map[string]struct {
		in   int64
		want string
	}{
		"zero":              {0, "صفر"},
		"one":               {1, "یک"},
		"ten":               {10, "ده"},
		"teen":              {19, "نوزده"},
		"round ten":         {40, "چهل"},
		"tens and ones":     {78, "هفتاد و هشت"},
		"round hundred":     {300, "سیصد"},
		"hundred and rest":  {156, "صد و پنجاه و شش"},
		"bare thousand":     {1000, "هزار"},
		"thousands":         {2000, "دو هزار"},
		"full":              {156789, "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه"},
		"skips empty group": {1000001, "یک میلیون و یک"},
		"million":           {1000000, "یک میلیون"},
		"milliard":          {2500000000, "دو میلیارد و پانصد میلیون"},
		"negative":          {-10, "منفی ده"},
		"negative full":     {-156788, "منفی صد و پنجاه و شش هزار و هفتصد و هشتاد و هشت"},
		"max int64":         {9223372036854775807, "نه کوینتیلیون و دویست و بیست و سه کوادریلیون و سیصد و هفتاد و دو تریلیون و سی و شش میلیارد و هشتصد و پنجاه و چهار میلیون و هفتصد و هفتاد و پنج هزار و هشتصد و هفت"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := ToWords(tt.in); got != tt.want {
				t.Errorf("ToWords(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToWordsMinInt64(t *testing.T) {
	// The magnitude of math.MinInt64 does not fit in an int64, so it is worth
	// its own case: it must not panic or overflow.
	got := ToWords(-9223372036854775808)
	const want = "منفی نه کوینتیلیون و دویست و بیست و سه کوادریلیون و سیصد و هفتاد و دو تریلیون و سی و شش میلیارد و هشتصد و پنجاه و چهار میلیون و هفتصد و هفتاد و پنج هزار و هشتصد و هشت"
	if got != want {
		t.Errorf("ToWords(math.MinInt64) = %q, want %q", got, want)
	}
}

func TestDigitToWord(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"persian digits": {"۱۵۶۷۸۹", "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه"},
		"ascii digits":   {"156789", "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه"},
		"arabic digits":  {"١٥٦٧٨٩", "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه"},
		"negative":       {"-10", "منفی ده"},
		"grouped":        {"156,789", "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه"},
		"not a number":   {"salam", ""},
		"empty":          {"", ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := DigitToWord(tt.in); got != tt.want {
				t.Errorf("DigitToWord(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestAddCommas(t *testing.T) {
	tests := map[string]struct {
		in   int64
		want string
	}{
		"zero":        {0, "0"},
		"below group": {999, "999"},
		"one group":   {1000, "1,000"},
		"large":       {14555478854, "14,555,478,854"},
		"negative":    {-1234567, "-1,234,567"},
		"min int64":   {-9223372036854775808, "-9,223,372,036,854,775,808"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := AddCommas(tt.in); got != tt.want {
				t.Errorf("AddCommas(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestRemoveCommas(t *testing.T) {
	tests := map[string]struct {
		in      string
		want    int64
		wantErr bool
	}{
		"ascii grouped":     {"4,555,522,212", 4555522212, false},
		"ungrouped":         {"455552221212", 455552221212, false},
		"persian separator": {"۱۲۳،۴۵۶", 123456, false},
		"persian digits":    {"۱۴۰۲", 1402, false},
		"negative":          {"-1,234", -1234, false},
		"surrounding space": {" 1,234 ", 1234, false},
		"not a number":      {"not-a-number", 0, true},
		"empty":             {"", 0, true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := RemoveCommas(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RemoveCommas(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("RemoveCommas(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
