package phonenumbers

import (
	"reflect"
	"testing"
)

func TestIsPhoneValid(t *testing.T) {
	tests := []struct {
		phoneNumber string
		isValid     bool
	}{
		{"9122221811", true},
		{"09122221811", true},
		{"+989122221811", true},
		{"12903908", false},
		{"901239812390812908", false},
	}

	for _, tt := range tests {
		ok := IsPhoneValid(tt.phoneNumber)
		if ok != tt.isValid {
			t.Errorf("IsPhoneValid(%s) got %t, expected %t", tt.phoneNumber, ok, tt.isValid)
		}
	}
}

func TestGetPhonePrefixOperator(t *testing.T) {
	expectedDetails904 := &OperatorDetails{
		base:     "کشوری",
		province: []string{},
		simTypes: []SimType{Credit},
		operator: Irancell,
		model:    "سیم\u200cکارت کودک",
	}

	if details, err := GetPrefixDetails("904"); err != nil {
		t.Errorf("expected no error, got %v", err)
	} else if !reflect.DeepEqual(details, expectedDetails904) {
		t.Errorf("expected %v, got %v", expectedDetails904, details)
	}

	expectedDetails910 := &OperatorDetails{
		base:     "کشوری",
		province: []string{},
		simTypes: []SimType{Permanent, Credit},
		operator: MCI,
	}

	if details, err := GetPrefixDetails("910"); err != nil {
		t.Errorf("expected no error, got %v", err)
	} else if !reflect.DeepEqual(details, expectedDetails910) {
		t.Errorf("expected %v, got %v", expectedDetails910, details)
	}

	if _, err := GetPrefixDetails("9100"); err == nil {
		t.Error("expected error, got none")
	}
}
