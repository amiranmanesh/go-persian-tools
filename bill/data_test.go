package bill

import "testing"

func TestBillTypes(t *testing.T) {
	tests := map[string]struct {
		billID int
		want   string
	}{
		"water":        {1117753200110, "آب"},
		"electricity":  {1117753200120, "برق"},
		"gas":          {1117753200130, "گاز"},
		"landline":     {1117753200140, "تلفن ثابت"},
		"mobile":       {9100074409150, "تلفن همراه"},
		"municipality": {1117753200160, "عوارض شهرداری"},
		"unknown zero": {1117753200100, "unknown"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := GetBillType(Params{BillID: tt.billID}); got != tt.want {
				t.Errorf("GetBillType(%d) = %q, want %q", tt.billID, got, tt.want)
			}
		})
	}
}

func TestVerifyBillIDEdgeCases(t *testing.T) {
	tests := map[string]struct {
		billID int
		want   bool
	}{
		"valid":        {7748317800142, true},
		"too short":    {12345, false},
		"zero":         {0, false},
		"unknown type": {1117753200100, false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := VerifyBillID(Params{BillID: tt.billID}); got != tt.want {
				t.Errorf("VerifyBillID(%d) = %v, want %v", tt.billID, got, tt.want)
			}
		})
	}
}

func TestGetCurrencyUnits(t *testing.T) {
	const paymentID = 12070160
	tests := map[string]struct {
		currency Currency
		want     int
	}{
		"unset defaults to rial": {Currency{}, 120000},
		"rial":                   {Currency{Rial: true}, 120000},
		"toman":                  {Currency{Toman: true}, 12000},
		"both prefers rial":      {Currency{Toman: true, Rial: true}, 120000},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := GetCurrency(Params{BillID: 1117753200140, PaymentID: paymentID, Currency: tt.currency})
			if got != tt.want {
				t.Errorf("GetCurrency(%v) = %d, want %d", tt.currency, got, tt.want)
			}
		})
	}
}

// TestDeprecatedAliasCompiles keeps the old type name usable.
func TestDeprecatedAliasCompiles(t *testing.T) {
	// Constructing through the alias is the point: it must still name the
	// same type, and the renamed fields must be reachable through it.
	params := BillParams{BillID: 1117753200140, PaymentID: 12070160}
	if GetBillType(params) != "تلفن ثابت" {
		t.Error("the deprecated BillParams alias no longer behaves like Params")
	}
}
