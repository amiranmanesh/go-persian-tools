// Package bill provides helpers to parse and validate Iranian utility bills
// (bill id / payment id pairs) and to derive their type, amount and barcode.
package bill

import (
	"strconv"
)

// Currency selects the unit used when computing a bill amount.
type Currency struct {
	Toman bool
	Rial  bool
}

type billType struct {
	bNum  int
	bType string
}

// billTypes maps the bill's type digit to its Persian label.
var billTypes = map[int]billType{
	0:  {1, "آب"},
	1:  {2, "برق"},
	2:  {3, "گاز"},
	3:  {4, "تلفن ثابت"},
	4:  {5, "تلفن همراه"},
	5:  {6, "عوارض شهرداری"},
	6:  {8, "سازمان مالیات"},
	7:  {9, "جرایم راهنمایی و رانندگی"},
	10: {0, "unknown"},
}

// Params describes a bill to be inspected.
type Params struct {
	BillID    int
	PaymentID int
	Currency  Currency
	Barcode   string
}

// BillParams is the former name of [Params].
//
// Deprecated: use [Params]. Note that the BillId and PaymentId fields were
// renamed to BillID and PaymentID at the same time, to match Go's naming
// conventions.
type BillParams = Params //nolint:revive // deprecated alias kept for compatibility

func bills(id int) billType {
	return billTypes[id]
}

// GetBillType returns the Persian label of the bill's service type, or
// "unknown" when the type digit is not recognized.
func GetBillType(billParams Params) string {
	billIDString := strconv.Itoa(billParams.BillID)
	billIDLen := len(billIDString)
	billID, _ := strconv.Atoi(billIDString[billIDLen-2 : billIDLen-1])
	if billID == 0 {
		return bills(10).bType
	}
	return bills(billID - 1).bType
}

// VerifyBillID reports whether the bill id is valid: it checks the trailing
// control digit and that the bill maps to a known type.
func VerifyBillID(billParams Params) bool {
	newBillID := strconv.Itoa(billParams.BillID)
	result := false

	if len(newBillID) < 6 {
		return false
	}

	controlBit := newBillID[len(newBillID)-1:]
	newBillID = newBillID[:len(newBillID)-1]

	c := calTheBit(newBillID)
	controlInt, _ := strconv.Atoi(controlBit)
	result = c == controlInt

	billType := GetBillType(billParams)

	return result && billType != "unknown"
}

// GetBarCode builds the payment barcode from the bill id and payment id.
func GetBarCode(billParams Params) string {
	billID := strconv.Itoa(billParams.BillID)
	paymentID := strconv.Itoa(billParams.PaymentID)
	return billID + "000" + paymentID
}

// GetCurrency returns the bill amount, scaled to Toman or Rial per the
// Currency flags (Rial, or an unset Currency, uses the 1000 multiplier).
func GetCurrency(billParams Params) int {
	currency := 100
	if billParams.Currency.Rial || !billParams.Currency.Toman {
		currency = 1000
	}

	payment := strconv.Itoa(billParams.PaymentID)
	payAmount, _ := strconv.Atoi(payment[0 : len(payment)-5])

	var amount = payAmount * currency

	return amount
}

func calTheBit(num string) int {
	sum := 0
	Base := 2
	for i := 0; i < len(num); i++ {
		if Base == 8 {
			Base = 2
		}

		subString := num[len(num)-1-i : len(num)-i]
		subStringInt, _ := strconv.Atoi(subString)
		sum += subStringInt * Base
		Base++
	}
	sum %= 11
	if sum < 2 {
		sum = 0
	} else {
		sum = 11 - sum
	}
	return sum
}
