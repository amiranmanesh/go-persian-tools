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
