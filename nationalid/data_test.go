package nationalid

import (
	"strconv"
	"strings"
	"testing"
)

// TestPlaceLookups pins the resolutions that the dataset corrections changed,
// alongside a few that must not move.
func TestPlaceLookups(t *testing.T) {
	tests := map[string]struct {
		id             string
		city, province string
		codes          []string
	}{
		"tehran":       {"0499370899", "شهرری", "تهران", []string{"048", "049"}},
		"zanjan":       {"4280000000", "زنجان", "زنجان", []string{"428", "427"}},
		"abhar":        {"6150000000", "ابهر", "زنجان", []string{"615"}},
		"marand":       {"1580000000", "مرند", "آذربایجان شرقی", []string{"158"}},
		"malekan":      {"5070000000", "ملکان", "آذربایجان شرقی", []string{"507"}},
		"mianeh":       {"1520000000", "میانه", "آذربایجان شرقی", []string{"152", "153"}},
		"azna":         {"4830000000", "ازنا", "لرستان", []string{"483", "484"}},
		"tabriz":       {"1360000000", "تبریز", "آذربایجان شرقی", []string{"136", "137", "138"}},
		"estahban":     {"2520000000", "استهبان", "فارس", []string{"252"}},
		"eghlid":       {"2530000000", "اقلید", "فارس", []string{"253"}},
		"unknown code": {"9990000000", "", "", nil},
		"too short":    {"123", "", "", nil},
		"too long":     {"059499370899", "", "", nil},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetPlaceByIranNationalID(tt.id)
			if got.City != tt.city || got.Province != tt.province {
				t.Errorf("place = %q/%q, want %q/%q", got.City, got.Province, tt.city, tt.province)
			}
			if len(got.Codes) != len(tt.codes) {
				t.Fatalf("codes = %q, want %q", got.Codes, tt.codes)
			}
			for i, code := range got.Codes {
				if code != tt.codes[i] {
					t.Errorf("codes[%d] = %q, want %q", i, code, tt.codes[i])
				}
			}
		})
	}
}

// TestDeprecatedAliasAgrees keeps the old spelling working.
func TestDeprecatedAliasAgrees(t *testing.T) {
	old := GetPlaceByIranNationalId("0499370899") //nolint:staticcheck // exercising the deprecated alias on purpose
	current := GetPlaceByIranNationalID("0499370899")
	if old.City != current.City || old.Province != current.Province {
		t.Errorf("deprecated alias returned %v, want %v", old, current)
	}
}

// TestDatasetIsWellFormed guards the shape of the tables themselves, so a bad
// hand edit fails here rather than silently returning the wrong province.
func TestDatasetIsWellFormed(t *testing.T) {
	provinces := make(map[int]string)
	for _, p := range getProvincesCode() {
		if _, dup := provinces[p.code]; dup {
			t.Errorf("duplicate province code %d", p.code)
		}
		if strings.TrimSpace(p.city) != p.city || p.city == "" {
			t.Errorf("province %d has untrimmed or empty name %q", p.code, p.city)
		}
		provinces[p.code] = p.city
	}

	for _, c := range getNationalCodes() {
		if strings.TrimSpace(c.city) != c.city || c.city == "" {
			t.Errorf("city %q is untrimmed or empty", c.city)
		}
		if strings.TrimSpace(c.code) != c.code {
			t.Errorf("city %q has an untrimmed code %q", c.city, c.code)
		}
		for _, part := range strings.Split(c.code, "-") {
			if len(part) != 3 {
				t.Errorf("city %q has code part %q that is not three digits", c.city, part)
				continue
			}
			if _, err := strconv.Atoi(part); err != nil {
				t.Errorf("city %q has non-numeric code part %q", c.city, part)
			}
		}
		if _, ok := provinces[c.parentCode]; !ok {
			t.Errorf("city %q points at unknown province %d", c.city, c.parentCode)
		}
	}
}

// TestEveryCityIsReachable checks that each code prefix in the table actually
// resolves through the public API.
func TestEveryCityIsReachable(t *testing.T) {
	for _, c := range getNationalCodes() {
		prefix := strings.Split(c.code, "-")[0]
		place := GetPlaceByIranNationalID(prefix + "0000000")
		if place.City == "" {
			t.Errorf("code %q (%s) resolves to nothing", prefix, c.city)
		}
		if place.Province == "" {
			t.Errorf("code %q (%s) resolves to no province", prefix, c.city)
		}
	}
}

func TestValidateEdgeCases(t *testing.T) {
	tests := map[string]struct {
		code string
		want bool
	}{
		"valid":               {"0067749828", true},
		"invalid checksum":    {"0684159415", false},
		"repeated ones":       {"1111111111", true},
		"repeated nines":      {"9999999999", false},
		"repeated zeros":      {"0000000000", false},
		"empty":               {"", false},
		"too short":           {"123456789", false},
		"too long":            {"12345678901", false},
		"letters":             {"abcdefghij", false},
		"persian digits":      {"۰۰۶۷۷۴۹۸۲۸", false},
		"zero middle section": {"0000012345", false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := Validate(tt.code); got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}
