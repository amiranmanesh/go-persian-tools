package nationalid

import (
	"reflect"
	"testing"
)

func TestGetPlaceByIranNationalId(t *testing.T) {
	gpbnID := GetPlaceByIranNationalID("0499370899")
	if gpbnID.City != "شهرری" || gpbnID.Province != "تهران" || !reflect.DeepEqual(gpbnID.Codes, []string{"048", "049"}) {
		t.Errorf("Result is false : %v", gpbnID)
	}
	wrongResult := GetPlaceByIranNationalID("059499370899")
	if wrongResult.City != "" || wrongResult.Province != "" || len(wrongResult.Codes) > 0 {
		t.Errorf("Result is false : %v", wrongResult)
	}
}

func TestValidate(t *testing.T) {
	verifyIranianNationalID := Validate("0067749828")
	verifyIranianNationalIDFalse := Validate("0684159415")
	verifyIranianNationalExceptionIDTrue := Validate("1111111111")
	verifyIranianNationalExceptionIDFalse := Validate("9999999999")

	if !verifyIranianNationalID {
		t.Errorf("Result is false : %v", verifyIranianNationalID)
	}
	if verifyIranianNationalIDFalse {
		t.Errorf("Result is false : %v", verifyIranianNationalIDFalse)
	}
	if !verifyIranianNationalExceptionIDTrue {
		t.Errorf("Result is True : %v", verifyIranianNationalExceptionIDTrue)
	}
	if verifyIranianNationalExceptionIDFalse {
		t.Errorf("Result is false : %v", verifyIranianNationalExceptionIDFalse)
	}
}

// Digits 4-9 are the serial number. Zeros in digits 4-6 alone are an ordinary
// serial; only a serial that is entirely zero is refused.
func TestValidateSerialNumber(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{"zeros in digits 4-6, serial not zero", "0490001238", true},
		{"zeros in digits 4-6, other office", "1700004565", true},
		{"whole serial is zero, valid check digit", "0490000002", false},
		{"whole serial is zero, other office", "1700000004", false},
		{"serial not zero, wrong check digit", "0490001237", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Validate(tt.code); got != tt.want {
				t.Errorf("Validate(%q) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}
