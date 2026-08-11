package bank

import (
	"errors"
	"regexp"
	"strings"
	"testing"
)

func TestCardInfoErrors(t *testing.T) {
	tests := map[string]struct {
		card string
		want string
		err  error
	}{
		"keshavarzi":       {"6037701689095443", "keshavarzi", nil},
		"saman":            {"6219861034529007", "saman", nil},
		"too short":        {"62198610", "", ErrInvalidCard},
		"letters":          {"603770asdfgbvcfg", "", ErrInvalidCard},
		"bad luhn":         {"6219861034529008", "", ErrInvalidCard},
		"empty":            {"", "", ErrInvalidCard},
		"too long":         {"62198610345290071", "", ErrInvalidCard},
		"persian digits":   {"۶۰۳۷۷۰۱۶۸۹۰۹۵۴۴۳", "", ErrInvalidCard},
		"unknown prefix":   {"9999999999999995", "", ErrBankNotFound},
		"tourism new bin":  {"5054260000000004", "tourism", nil},
		"central bank bin": {"6367970000000008", "central-bank", nil},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := CardInfo(tt.card)
			if !errors.Is(err, tt.err) {
				t.Fatalf("CardInfo(%q) error = %v, want %v", tt.card, err, tt.err)
			}
			if got != tt.want {
				t.Errorf("CardInfo(%q) = %q, want %q", tt.card, got, tt.want)
			}
		})
	}
}

// TestBankCodeTableIsWellFormed guards the card prefix table.
func TestBankCodeTableIsWellFormed(t *testing.T) {
	prefix := regexp.MustCompile(`^\d{6}$`)
	for bin, name := range bankCode {
		if !prefix.MatchString(bin) {
			t.Errorf("card prefix %q is not six digits", bin)
		}
		if strings.TrimSpace(name) != name || name == "" {
			t.Errorf("prefix %q has an untrimmed or empty bank name %q", bin, name)
		}
	}
}

func TestShebaValidation(t *testing.T) {
	tests := map[string]struct {
		code  string
		valid bool
		bank  string
	}{
		"parsian":            {"IR820540102680020817909002", true, "Parsian Bank"},
		"bad checksum":       {"IR820540102680020817909003", false, ""},
		"too short":          {"IR8205401026800208179090", false, ""},
		"too long":           {"IR8205401026800208179090021", false, ""},
		"empty":              {"", false, ""},
		"missing IR":         {"XX820540102680020817909002", false, ""},
		"letters in digits":  {"IR8205401026800208179090AB", false, ""},
		"unknown bank code":  {"IR062960000000100324200001", true, ""},
		"lowercase rejected": {"ir820540102680020817909002", false, ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sheba := ShebaCode{Code: tt.code}
			if got := sheba.IsValid(); got != tt.valid {
				t.Errorf("IsValid(%q) = %v, want %v", tt.code, got, tt.valid)
			}
			if got := sheba.IsSheba(); got.Name != tt.bank {
				t.Errorf("IsSheba(%q).Name = %q, want %q", tt.code, got.Name, tt.bank)
			}
		})
	}
}

// TestShebaTableIsWellFormed guards the bank table and the merger annotations.
func TestShebaTableIsWellFormed(t *testing.T) {
	code := regexp.MustCompile(`^\d{3}$`)
	seen := map[string]bool{}

	for _, want := range []string{"010", "012", "017", "054", "057", "061"} {
		if got := shebaHashTable(want); got.Code != want {
			t.Errorf("shebaHashTable(%q) returned code %q", want, got.Code)
		}
	}
	if got := shebaHashTable("000"); got.Name != "" {
		t.Errorf("shebaHashTable of an unknown code returned %v", got)
	}

	// Walk every code the table can answer for.
	for i := 0; i < 1000; i++ {
		key := string(rune('0'+i/100)) + string(rune('0'+i/10%10)) + string(rune('0'+i%10))
		bank := shebaHashTable(key)
		if bank.Name == "" {
			continue
		}
		if !code.MatchString(bank.Code) {
			t.Errorf("bank %q has a malformed code %q", bank.Name, bank.Code)
		}
		if seen[bank.Code] {
			t.Errorf("duplicate bank code %q", bank.Code)
		}
		seen[bank.Code] = true
		if bank.NickName == "" || bank.PersianName == "" {
			t.Errorf("bank %q is missing a nickname or Persian name", bank.Code)
		}
		if bank.AccountNumberAvailable && bank.Process == nil {
			t.Errorf("bank %q claims account numbers but has no Process", bank.Code)
		}
	}

	if len(seen) < 30 {
		t.Errorf("only %d banks in the table, expected the full list", len(seen))
	}
}

// TestMergedBanksStillResolve checks that the institutions absorbed by Sepah
// keep resolving, since IBANs issued before the merger are still in the wild.
func TestMergedBanksStillResolve(t *testing.T) {
	merged := map[string]string{
		"052": "ghavamin",
		"063": "ansar",
		"065": "hekmat-iranian",
		"073": "kosar",
		"079": "mehr-eqtesad",
	}
	for code, nickname := range merged {
		bank := shebaHashTable(code)
		if bank.NickName != nickname {
			t.Errorf("code %q resolved to %q, want %q", code, bank.NickName, nickname)
		}
		if bank.MergedInto != "sepah" {
			t.Errorf("bank %q has MergedInto %q, want \"sepah\"", nickname, bank.MergedInto)
		}
	}
	if got := shebaHashTable("015"); got.MergedInto != "" {
		t.Errorf("Sepah itself is marked as merged into %q", got.MergedInto)
	}
}

// TestShebaAccountNumberProcessors exercises the per-bank account formatters.
func TestShebaAccountNumberProcessors(t *testing.T) {
	tests := map[string]struct {
		fn    func(string) ShebaProcess
		code  string
		empty bool
	}{
		"parsian":   {parsian, "IR820540102680020817909002", false},
		"pasargad":  {pasargad, "IR870570023880010930909001", false},
		"bankshahr": {bankshahr, "IR230610000000700796959601", false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := tt.fn(tt.code)
			if got.normal == "" {
				t.Errorf("%s(%q) produced an empty account number", name, tt.code)
			}
			if got.formatted == "" {
				t.Errorf("%s(%q) produced an empty formatted account number", name, tt.code)
			}
		})
	}
}
