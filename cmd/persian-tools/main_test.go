package main

import (
	"bytes"
	"strings"
	"testing"
)

func exec(t *testing.T, stdin string, argv ...string) (code int, stdout, stderr string) {
	t.Helper()

	var out, errOut bytes.Buffer
	code, err := run(argv, strings.NewReader(stdin), &out, &errOut)
	if err != nil {
		errOut.WriteString(err.Error())
	}
	return code, out.String(), errOut.String()
}

func TestTransformsFromArguments(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		argv []string
		want string
	}{
		"normalize":      {[]string{"normalize", "علي\u200Cرضا"}, "علی رضا\n"},
		"fix arabic":     {[]string{"fix-arabic", "علي كريم"}, "علی کریم\n"},
		"to persian":     {[]string{"to-persian", "1402"}, "۱۴۰۲\n"},
		"to english":     {[]string{"to-english", "۰۹۱۲٣٤٥"}, "0912345\n"},
		"digits default": {[]string{"digits", "a1ب۲"}, "1۲\n"},
		"digits en":      {[]string{"digits", "-set", "en", "a1ب۲"}, "1\n"},
		"digits fa":      {[]string{"digits", "-set", "fa", "a1ب۲"}, "۲\n"},
		"alpha":          {[]string{"alpha", "123شاهینhi"}, "شاهین\n"},
		"currency":       {[]string{"currency", "1234567"}, "۱،۲۳۴،۵۶۷\n"},
		"toman":          {[]string{"currency", "-unit", "toman", "1234567"}, "۱،۲۳۴،۵۶۷ تومان\n"},
		"rial":           {[]string{"currency", "-unit", "rial", "1234567"}, "۱،۲۳۴،۵۶۷ ﷼\n"},
		"key to persian": {[]string{"key-to-persian", "sghl"}, "سلام\n"},
		"key to english": {[]string{"key-to-english", "اثغ"}, "hey\n"},
		"finglish":       {[]string{"finglish", "سلام"}, "salam\n"},
		"reverse":        {[]string{"reverse", "سلام"}, "مالس\n"},
		"joins argv":     {[]string{"key-to-persian", "sghl", "o,fd"}, "سلام خوبی\n"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			code, out, _ := exec(t, "", tt.argv...)
			if code != 0 {
				t.Errorf("exit code = %d, want 0", code)
			}
			if out != tt.want {
				t.Errorf("stdout = %q, want %q", out, tt.want)
			}
		})
	}
}

func TestReadsStdinWhenNoArguments(t *testing.T) {
	t.Parallel()

	// A trailing newline in the input is preserved rather than doubled.
	code, out, _ := exec(t, "علي\u200Cرضا\n", "normalize")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if want := "علی رضا\n"; out != want {
		t.Errorf("stdout = %q, want %q", out, want)
	}
}

func TestStdinWithoutTrailingNewlineGetsOne(t *testing.T) {
	t.Parallel()

	code, out, _ := exec(t, "1402", "to-persian")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if want := "۱۴۰۲\n"; out != want {
		t.Errorf("stdout = %q, want %q", out, want)
	}
}

func TestMultilineStdinKeepsItsLines(t *testing.T) {
	t.Parallel()

	code, out, _ := exec(t, "علي\nكريم\n", "fix-arabic")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if want := "علی\nکریم\n"; out != want {
		t.Errorf("stdout = %q, want %q", out, want)
	}
}

func TestIsEnglishExitCode(t *testing.T) {
	t.Parallel()

	if code, out, _ := exec(t, "", "is-english", "ali"); code != 0 || out != "true\n" {
		t.Errorf("is-english ali = (%d, %q), want (0, \"true\\n\")", code, out)
	}
	if code, out, _ := exec(t, "", "is-english", "علی"); code != 1 || out != "false\n" {
		t.Errorf("is-english علی = (%d, %q), want (1, \"false\\n\")", code, out)
	}
}

func TestHelpAndVersion(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"-h", "--help", "help"} {
		code, out, _ := exec(t, "", arg)
		if code != 0 {
			t.Errorf("%s exit code = %d, want 0", arg, code)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("%s did not print usage", arg)
		}
		for _, c := range commands {
			if !strings.Contains(out, c.name) {
				t.Errorf("%s did not list command %q", arg, c.name)
			}
		}
	}

	for _, arg := range []string{"-v", "--version", "version"} {
		code, out, _ := exec(t, "", arg)
		if code != 0 || strings.TrimSpace(out) == "" {
			t.Errorf("%s = (%d, %q), want a version on exit 0", arg, code, out)
		}
	}
}

func TestUsageErrors(t *testing.T) {
	t.Parallel()

	if code, _, errOut := exec(t, "", "nope"); code != 2 || !strings.Contains(errOut, "unknown command") {
		t.Errorf("unknown command = (%d, %q), want exit 2 and a message", code, errOut)
	}
	if code, _, errOut := exec(t, ""); code != 2 || !strings.Contains(errOut, "Usage:") {
		t.Errorf("no arguments = (%d, %q), want exit 2 and usage", code, errOut)
	}
	if code, _, _ := exec(t, "", "normalize", "-nope"); code != 2 {
		t.Errorf("unknown flag exit code = %d, want 2", code)
	}
}

// TestEveryCommandIsReachable keeps the table and the help output honest: a
// command added to one has to work in the other.
func TestEveryCommandIsReachable(t *testing.T) {
	t.Parallel()

	seen := make(map[string]bool, len(commands))
	for _, c := range commands {
		if seen[c.name] {
			t.Errorf("duplicate command %q", c.name)
		}
		seen[c.name] = true

		if c.summary == "" {
			t.Errorf("command %q has no summary", c.name)
		}
		if code, _, _ := exec(t, "سلام 1402", c.name); code > 1 {
			t.Errorf("command %q exited %d on ordinary input", c.name, code)
		}
	}
}

func TestToolCommands(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		argv []string
		want string
		code int
	}{
		"words":             {[]string{"words", "156789"}, "صد و پنجاه و شش هزار و هفتصد و هشتاد و نه\n", 0},
		"words persian":     {[]string{"words", "۱۴۰۲"}, "هزار و چهارصد و دو\n", 0},
		"commas":            {[]string{"commas", "14555478854"}, "14,555,478,854\n", 0},
		"national id ok":    {[]string{"national-id", "0067749828"}, "true\n", 0},
		"national id bad":   {[]string{"national-id", "0684159415"}, "false\n", 1},
		"national id place": {[]string{"national-id-place", "0499370899"}, "شهرری، تهران\n", 0},
		"card":              {[]string{"card", "6037701689095443"}, "keshavarzi\n", 0},
		"card invalid":      {[]string{"card", "6219861034529008"}, "\n", 0},
		"sheba":             {[]string{"sheba", "IR820540102680020817909002"}, "Parsian Bank — بانک پارسیان\n", 0},
		"phone ok":          {[]string{"phone", "09122221811"}, "true\n", 0},
		"phone bad":         {[]string{"phone", "12903908"}, "false\n", 1},
		"phone operator":    {[]string{"phone-operator", "09123456789"}, "MCI\n", 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			code, out, _ := exec(t, "", tt.argv...)
			if code != tt.code {
				t.Errorf("exit code = %d, want %d", code, tt.code)
			}
			if out != tt.want {
				t.Errorf("stdout = %q, want %q", out, tt.want)
			}
		})
	}
}

// TestEveryCommandIsGrouped keeps the help output complete: a command with no
// group, or an unknown one, would silently vanish from the listing.
func TestEveryCommandIsGrouped(t *testing.T) {
	t.Parallel()

	known := make(map[string]bool, len(groups))
	for _, g := range groups {
		known[g.name] = true
	}
	for _, c := range commands {
		if !known[c.group] {
			t.Errorf("command %q has unknown group %q", c.name, c.group)
		}
	}
}
