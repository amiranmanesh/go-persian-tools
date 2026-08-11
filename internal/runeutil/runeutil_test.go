package runeutil

import "testing"

func TestKeep(t *testing.T) {
	isDigit := func(r rune) bool { return r >= '0' && r <= '9' }
	tests := map[string]struct{ in, want string }{
		"mixed":       {"a1b2c3", "123"},
		"none kept":   {"abc", ""},
		"all kept":    {"123", "123"},
		"empty":       {"", ""},
		"multi-byte":  {"س1ل2ا3م", "123"},
		"invalid utf": {"\xff1\xfe2", "12"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := Keep(tt.in, isDigit); got != tt.want {
				t.Errorf("Keep(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestReplacer(t *testing.T) {
	pairs := [][2]string{{"a", "ا"}, {"b", "ب"}}

	forward := Replacer(pairs, false)
	if got := forward.Replace("ab c"); got != "اب c" {
		t.Errorf("forward = %q, want %q", got, "اب c")
	}

	reverse := Replacer(pairs, true)
	if got := reverse.Replace("اب c"); got != "ab c" {
		t.Errorf("reverse = %q, want %q", got, "ab c")
	}

	if got := Replacer(nil, false).Replace("untouched"); got != "untouched" {
		t.Errorf("empty replacer changed the input: %q", got)
	}
}

func TestConcatPairs(t *testing.T) {
	a := [][2]string{{"1", "one"}}
	b := [][2]string{{"2", "two"}}

	got := ConcatPairs(a, b)
	if len(got) != 2 || got[0] != a[0] || got[1] != b[0] {
		t.Errorf("ConcatPairs = %v, want the tables in order", got)
	}
	if got := ConcatPairs(); got != nil {
		t.Errorf("ConcatPairs() = %v, want nil", got)
	}
	// Order must be preserved: the first table wins for a duplicate key.
	dup := ConcatPairs([][2]string{{"x", "first"}}, [][2]string{{"x", "second"}})
	if Replacer(dup, false).Replace("x") != "first" {
		t.Error("ConcatPairs did not preserve table order")
	}
}
