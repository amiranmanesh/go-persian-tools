package bill

import "testing"

// Every test builds its own Params. The suite used to share one package-level
// value and mutate it, which made the results depend on the order the tests
// happened to run in — and `go test -shuffle=on` duly broke it.

func TestGetBillTypeCases(t *testing.T) {
	tests := map[string]struct {
		params Params
		want   string
	}{
		"landline": {Params{BillID: 1117753200140, PaymentID: 12070160}, "تلفن ثابت"},
		"mobile":   {Params{BillID: 9100074409151, PaymentID: 12908190}, "تلفن همراه"},
		"unknown":  {Params{BillID: 7748317800105, PaymentID: 1770160}, "unknown"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := GetBillType(tt.params); got != tt.want {
				t.Errorf("GetBillType(%d) = %q, want %q", tt.params.BillID, got, tt.want)
			}
		})
	}
}

func TestGetCurrencyRialAndToman(t *testing.T) {
	params := Params{BillID: 1117753200140, PaymentID: 1770160}

	if got := GetCurrency(params); got != 17000 {
		t.Errorf("GetCurrency with the default currency = %d, want 17000", got)
	}

	params.Currency = Currency{Toman: true}
	if got := GetCurrency(params); got != 1700 {
		t.Errorf("GetCurrency in Toman = %d, want 1700", got)
	}
}

func TestVerifyBillIDCases(t *testing.T) {
	tests := map[string]struct {
		params Params
		want   bool
	}{
		"valid":   {Params{BillID: 7748317800142, PaymentID: 1770160}, true},
		"invalid": {Params{BillID: 2234322344613, PaymentID: 1070189}, false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := VerifyBillID(tt.params); got != tt.want {
				t.Errorf("VerifyBillID(%d) = %v, want %v", tt.params.BillID, got, tt.want)
			}
		})
	}
}

func TestGetBarCode(t *testing.T) {
	tests := map[string]struct {
		params Params
		want   string
	}{
		"first":  {Params{BillID: 7748317800142, PaymentID: 1770160}, "77483178001420001770160"},
		"second": {Params{BillID: 9174639504124, PaymentID: 12908197}, "917463950412400012908197"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := GetBarCode(tt.params); got != tt.want {
				t.Errorf("GetBarCode(%d) = %q, want %q", tt.params.BillID, got, tt.want)
			}
		})
	}
}
