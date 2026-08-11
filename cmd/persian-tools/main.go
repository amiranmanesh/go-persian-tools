// Command persian-tools applies the transformations and validators of the
// go-persian-tools packages to text on the command line.
//
// Text comes from the arguments, or from standard input when there are none,
// so the tool works both as a one-off and as a filter in a pipeline:
//
//	persian-tools normalize "علي كريم"
//	cat names.txt | persian-tools normalize > keys.txt
//
// Every transformation is applied to the whole input at once. Line structure
// survives the transformations that work rune by rune; the filters
// (digits, alpha) drop newlines along with everything else they do not keep.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"

	"github.com/amiranmanesh/go-persian-tools/bank"
	"github.com/amiranmanesh/go-persian-tools/digit"
	"github.com/amiranmanesh/go-persian-tools/nationalid"
	"github.com/amiranmanesh/go-persian-tools/phonenumbers"
	"github.com/amiranmanesh/go-persian-tools/text"
)

// version is stamped by the release build with -ldflags -X main.version=vX.Y.Z.
// Outside that build it comes from the module the binary was installed from.
var version = ""

type command struct {
	name    string
	args    string
	group   string
	summary string

	// setup registers the command's flags and returns the transformation bound
	// to them. Commands without flags ignore the FlagSet.
	setup func(fs *flag.FlagSet) func(text string) string

	// predicate marks a command that answers yes or no: it prints "true" or
	// "false" and exits 0 or 1, so a shell can branch on it.
	predicate bool
}

func plain(fn func(string) string) func(*flag.FlagSet) func(string) string {
	return func(*flag.FlagSet) func(string) string { return fn }
}

// boolean adapts a predicate into the string-to-string shape of a command.
func boolean(fn func(string) bool) func(*flag.FlagSet) func(string) string {
	return plain(func(s string) string { return fmt.Sprint(fn(strings.TrimSpace(s))) })
}

var commands = []command{
	// text
	{
		name: "normalize", args: "[text]", group: "text",
		summary: "fold text for comparison, sorting and indexing",
		setup:   plain(text.Normalize),
	},
	{
		name: "fix-arabic", args: "[text]", group: "text",
		summary: "replace Arabic characters with their Persian equivalents",
		setup:   plain(text.FixArabic),
	},
	{
		name: "alpha", args: "[text]", group: "text",
		summary: "keep only Persian script",
		setup:   plain(text.OnlyPersianAlpha),
	},
	{
		name: "key-to-persian", args: "[text]", group: "text",
		summary: "re-read text typed on an English keyboard layout",
		setup:   plain(text.SwitchToPersianKey),
	},
	{
		name: "key-to-english", args: "[text]", group: "text",
		summary: "re-read text typed on a Persian keyboard layout",
		setup:   plain(text.SwitchToEnglishKey),
	},
	{
		name: "finglish", args: "[text]", group: "text",
		summary: "romanize Persian text",
		setup:   plain(text.Finglish),
	},
	{
		name: "reverse", args: "[text]", group: "text",
		summary: "reverse the text by runes",
		setup:   plain(text.Reverse),
	},
	{
		name: "is-english", args: "[text]", group: "text",
		summary:   "report whether the text is only ASCII letters",
		setup:     boolean(text.CheckIsEnglish),
		predicate: true,
	},

	// digits and money
	{
		name: "to-persian", args: "[text]", group: "digits",
		summary: "convert ASCII digits to Persian digits",
		setup:   plain(digit.ToPersianDigits),
	},
	{
		name: "to-english", args: "[text]", group: "digits",
		summary: "convert Persian and Arabic-Indic digits to ASCII",
		setup:   plain(digit.ToEnglishDigits),
	},
	{
		name: "digits", args: "[text]", group: "digits",
		summary: "keep only digits (-set en|fa|all)",
		setup: func(fs *flag.FlagSet) func(string) string {
			set := fs.String("set", "all", "which digits to keep: en, fa or all")
			return func(s string) string {
				switch *set {
				case "en":
					return digit.OnlyEnglishNumbers(s)
				case "fa":
					return digit.OnlyPersianNumbers(s)
				default:
					return digit.OnlyNumbers(s)
				}
			}
		},
	},
	{
		name: "currency", args: "[amount]", group: "digits",
		summary: "group an amount in threes (-unit toman|rial|none)",
		setup: func(fs *flag.FlagSet) func(string) string {
			unit := fs.String("unit", "none", "unit to append: toman, rial or none")
			return func(s string) string {
				switch *unit {
				case "toman":
					return digit.Toman(s)
				case "rial":
					return digit.Rial(s)
				default:
					return digit.Currency(s)
				}
			}
		},
	},
	{
		name: "commas", args: "[amount]", group: "digits",
		summary: "group an amount with ASCII commas",
		setup: plain(func(s string) string {
			n, err := digit.ParseInt(s)
			if err != nil {
				return ""
			}
			return digit.AddCommas(n)
		}),
	},
	{
		name: "words", args: "[amount]", group: "digits",
		summary: "spell an amount out in Persian words",
		setup:   plain(digit.DigitToWord),
	},

	// validators
	{
		name: "national-id", args: "[code]", group: "validate",
		summary:   "validate an Iranian national number (code-e Melli)",
		setup:     boolean(nationalid.Validate),
		predicate: true,
	},
	{
		name: "national-id-place", args: "[code]", group: "validate",
		summary: "resolve the city and province of a national number",
		setup: plain(func(s string) string {
			place := nationalid.GetPlaceByIranNationalId(strings.TrimSpace(s))
			if place.City == "" {
				return ""
			}
			return place.City + "، " + place.Province
		}),
	},
	{
		name: "card", args: "[number]", group: "validate",
		summary: "resolve the bank that issued a card number",
		setup: plain(func(s string) string {
			name, err := bank.CardInfo(strings.TrimSpace(s))
			if err != nil {
				return ""
			}
			return name
		}),
	},
	{
		name: "sheba", args: "[code]", group: "validate",
		summary: "resolve the bank of a Sheba (IBAN) code",
		setup: plain(func(s string) string {
			result := bank.ShebaCode{Code: strings.TrimSpace(s)}.IsSheba()
			if result.Name == "" {
				return ""
			}
			return result.Name + " — " + result.PersianName
		}),
	},
	{
		name: "phone", args: "[number]", group: "validate",
		summary:   "report whether an Iranian mobile number is valid",
		setup:     boolean(phonenumbers.IsPhoneValid),
		predicate: true,
	},
	{
		name: "phone-operator", args: "[number]", group: "validate",
		summary: "resolve the operator of an Iranian mobile number",
		setup: plain(func(s string) string {
			details, err := phonenumbers.GetPhoneDetails(strings.TrimSpace(s))
			if err != nil {
				return ""
			}
			return string(details.GetOperator())
		}),
	},
}

// groups fixes the order the command list is printed in.
var groups = []struct{ name, title string }{
	{"text", "Text"},
	{"digits", "Digits and money"},
	{"validate", "Validators"},
}

func main() {
	code, err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "persian-tools: %v\n", err)
	}
	os.Exit(code)
}

func run(argv []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	if len(argv) == 0 {
		usage(stderr)
		return 2, nil
	}

	switch argv[0] {
	case "-h", "--help", "help":
		usage(stdout)
		return 0, nil
	case "-v", "--version", "version":
		fmt.Fprintln(stdout, buildVersion())
		return 0, nil
	}

	cmd, ok := lookup(argv[0])
	if !ok {
		usage(stderr)
		return 2, fmt.Errorf("unknown command %q", argv[0])
	}

	fs := flag.NewFlagSet("persian-tools "+cmd.name, flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: persian-tools %s %s\n\n%s\n", cmd.name, cmd.args, cmd.summary)
		fs.PrintDefaults()
	}
	transform := cmd.setup(fs)
	if err := fs.Parse(argv[1:]); err != nil {
		return 2, nil // flag has already reported the problem
	}

	input, fromStdin, err := readInput(fs.Args(), stdin)
	if err != nil {
		return 1, err
	}

	out := transform(input)
	if _, err := io.WriteString(stdout, out); err != nil {
		return 1, err
	}
	// Keep a filter's output byte-faithful when the input already ended in a
	// newline, and still leave the terminal on its own line when it did not.
	if !fromStdin || (out != "" && !strings.HasSuffix(out, "\n")) {
		if _, err := io.WriteString(stdout, "\n"); err != nil {
			return 1, err
		}
	}

	if cmd.predicate && out == "false" {
		return 1, nil
	}
	return 0, nil
}

// readInput returns the text to transform: the arguments joined by a space, or
// all of standard input when there are no arguments.
func readInput(args []string, stdin io.Reader) (text string, fromStdin bool, err error) {
	if len(args) > 0 {
		return strings.Join(args, " "), false, nil
	}
	b, err := io.ReadAll(stdin)
	if err != nil {
		return "", true, fmt.Errorf("reading standard input: %w", err)
	}
	return string(b), true, nil
}

func lookup(name string) (command, bool) {
	for _, c := range commands {
		if c.name == name {
			return c, true
		}
	}
	return command{}, false
}

func buildVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return "devel"
}

func usage(w io.Writer) {
	fmt.Fprintf(w, `persian-tools %s — tools for Persian (Iranian) data

Usage:
  persian-tools <command> [flags] [text]
  <input> | persian-tools <command> [flags]

Text comes from the arguments, or from standard input when there are none.
`, buildVersion())

	width := 0
	for _, c := range commands {
		if n := len(c.name); n > width {
			width = n
		}
	}
	for _, g := range groups {
		fmt.Fprintf(w, "\n%s:\n", g.title)
		for _, c := range commands {
			if c.group == g.name {
				fmt.Fprintf(w, "  %-*s  %s\n", width, c.name, c.summary)
			}
		}
	}

	fmt.Fprintf(w, `
Run "persian-tools <command> -h" for a command's flags.

Examples:
  persian-tools normalize "علي كريم"
  persian-tools currency -unit toman 1234567
  persian-tools words 156789
  persian-tools key-to-persian sghl
  persian-tools national-id 0067749828
  cat names.txt | persian-tools normalize > keys.txt

Docs: https://pkg.go.dev/github.com/amiranmanesh/go-persian-tools
`)
}
