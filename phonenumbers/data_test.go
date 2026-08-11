package phonenumbers

import (
	"errors"
	"regexp"
	"testing"
)

// allTables lists every operator prefix table the package resolves against.
var allTables = map[Operator]map[string]OperatorDetails{
	MCI:          MCIMap,
	Taliya:       TALIYA,
	RightTel:     RIGHTTEL,
	Irancell:     IRANCELL,
	ShatelMobile: SHATELMOBILE,
	ApTel:        APTEL,
	TeleKish:     TELEKISH,
	Espadan:      ESPADAN,
}

func TestGetPhoneDetails(t *testing.T) {
	tests := map[string]struct {
		number   string
		operator Operator
		err      error
	}{
		"mci":              {"09123456789", MCI, nil},
		"mci no prefix":    {"9123456789", MCI, nil},
		"mci +98":          {"+989123456789", MCI, nil},
		"mci 0098":         {"00989123456789", MCI, nil},
		"irancell":         {"09351234567", Irancell, nil},
		"irancell td-lte":  {"09411234567", Irancell, nil},
		"rightel":          {"09201234567", RightTel, nil},
		"rightel new 924":  {"09241234567", RightTel, nil},
		"taliya":           {"09321234567", Taliya, nil},
		"shatel mobile":    {"09981234567", ShatelMobile, nil},
		"aptel":            {"09991234567", ApTel, nil},
		"tele kish":        {"09341234567", TeleKish, nil},
		"espadan":          {"09311234567", Espadan, nil},
		"too short":        {"0912345", "", ErrInvalidFormat},
		"too long":         {"091234567890", "", ErrInvalidFormat},
		"not a mobile":     {"02112345678", "", ErrInvalidFormat},
		"empty":            {"", "", ErrInvalidFormat},
		"unknown operator": {"09771234567", "", ErrInvalidPrefix},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			details, err := GetPhoneDetails(tt.number)
			if !errors.Is(err, tt.err) {
				t.Fatalf("GetPhoneDetails(%q) error = %v, want %v", tt.number, err, tt.err)
			}
			if tt.err != nil {
				return
			}
			if got := details.GetOperator(); got != tt.operator {
				t.Errorf("operator = %q, want %q", got, tt.operator)
			}
			if details.GetBase() == "" {
				t.Errorf("GetBase() is empty for %q", tt.number)
			}
			if details.GetSimTypeList() == nil {
				t.Errorf("GetSimTypeList() is nil for %q", tt.number)
			}
			// Exercise the remaining accessors; they may legitimately be empty.
			details.GetProvinceList()
			details.GetModel()
		})
	}
}

func TestOperatorDetailsAccessors(t *testing.T) {
	details, err := GetPhoneDetails("09141234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := details.GetBase(); got != "آذربایجان غربی" {
		t.Errorf("GetBase() = %q", got)
	}
	if got := details.GetProvinceList(); len(got) != 3 {
		t.Errorf("GetProvinceList() = %q, want three provinces", got)
	}
	if got := details.GetModel(); got != "" {
		t.Errorf("GetModel() = %q, want empty", got)
	}

	child, err := GetPhoneDetails("09041234567")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := child.GetModel(); got == "" {
		t.Error("GetModel() is empty for the 904 prefix, which has a model")
	}
}

func TestOperatorDetails(t *testing.T) {
	for operator, want := range allTables {
		t.Run(string(operator), func(t *testing.T) {
			got := operator.Details()
			if len(got) != len(want) {
				t.Errorf("Details() has %d prefixes, want %d", len(got), len(want))
			}
		})
	}
	if got := Operator("nope").Details(); got != nil {
		t.Errorf("Details() of an unknown operator = %v, want nil", got)
	}
}

func TestGetPrefixDetails(t *testing.T) {
	if _, err := GetPrefixDetails("912"); err != nil {
		t.Errorf("GetPrefixDetails(912): %v", err)
	}
	if _, err := GetPrefixDetails("777"); !errors.Is(err, ErrInvalidPrefix) {
		t.Errorf("GetPrefixDetails(777) error = %v, want ErrInvalidPrefix", err)
	}
}

func TestPhoneNumberNormalizer(t *testing.T) {
	tests := map[string]struct {
		number, prefix, want string
		wantErr              bool
	}{
		"to +98":        {"09122221811", "+98", "+989122221811", false},
		"to 0":          {"+989122221811", "0", "09122221811", false},
		"from bare":     {"9122221811", "0", "09122221811", false},
		"from 0098":     {"00989122221811", "+98", "+989122221811", false},
		"strip to bare": {"09122221811", "", "9122221811", false},
		"invalid":       {"12345", "+98", "", true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := PhoneNumberNormalizer(tt.number, tt.prefix)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("PhoneNumberNormalizer(%q, %q) = %q, want %q", tt.number, tt.prefix, got, tt.want)
			}
		})
	}
}

func TestGetOperatorPrefix(t *testing.T) {
	if got, err := GetOperatorPrefix("09123456789"); err != nil || got != "912" {
		t.Errorf("GetOperatorPrefix = %q, %v; want \"912\", nil", got, err)
	}
	if _, err := GetOperatorPrefix("123"); !errors.Is(err, ErrInvalidFormat) {
		t.Errorf("GetOperatorPrefix of an invalid number: %v", err)
	}
}

func TestGetPhonePrefix(t *testing.T) {
	tests := map[string]string{
		"09123456789":    "0",
		"+989123456789":  "+98",
		"00989123456789": "0098",
		"989123456789":   "98",
		"9123456789":     "",
	}
	for number, want := range tests {
		if got := GetPhonePrefix(number); got != want {
			t.Errorf("GetPhonePrefix(%q) = %q, want %q", number, got, want)
		}
	}
}

// TestPrefixTablesAreWellFormed guards the operator tables: every key is three
// digits, no prefix is claimed twice, and each entry is internally consistent.
func TestPrefixTablesAreWellFormed(t *testing.T) {
	shape := regexp.MustCompile(`^\d{3}$`)
	owner := map[string]Operator{}

	for operator, table := range allTables {
		for prefix, details := range table {
			if !shape.MatchString(prefix) {
				t.Errorf("%s: prefix %q is not three digits", operator, prefix)
			}
			if other, dup := owner[prefix]; dup {
				t.Errorf("prefix %q is claimed by both %s and %s", prefix, other, operator)
			}
			owner[prefix] = operator

			if details.operator != operator {
				t.Errorf("%s: prefix %q reports operator %q", operator, prefix, details.operator)
			}
			if details.base == "" {
				t.Errorf("%s: prefix %q has no base province", operator, prefix)
			}
			if len(details.simTypes) == 0 {
				t.Errorf("%s: prefix %q lists no SIM type", operator, prefix)
			}
		}
	}

	// Every prefix must be reachable through the public lookup.
	for prefix := range owner {
		if _, err := GetPrefixDetails(prefix); err != nil {
			t.Errorf("prefix %q is in a table but does not resolve: %v", prefix, err)
		}
	}
}
