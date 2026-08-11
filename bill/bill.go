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

// BillParams describes a bill to be inspected.
type BillParams struct {
	BillId    int
	PaymentId int
	Currency  Currency
	Barcode   string
}

func bills(id int) billType {
	return billTypes[id]
}

// GetBillType returns the Persian label of the bill's service type, or
// "unknown" when the type digit is not recognized.
func GetBillType(billParams BillParams) string {
	billIdString := strconv.Itoa(billParams.BillId)
	billIdLen := len(billIdString)
	billId, _ := strconv.Atoi(billIdString[billIdLen-2 : billIdLen-1])
	if billId == 0 {
		return bills(10).bType
	}
	return bills(billId - 1).bType
}

// VerifyBillID reports whether the bill id is valid: it checks the trailing
// control digit and that the bill maps to a known type.
func VerifyBillID(billParams BillParams) bool {
	newBillId := strconv.Itoa(billParams.BillId)
	result := false

	if len(newBillId) < 6 {
		return false
	}

	controlBit := newBillId[len(newBillId)-1:]
	newBillId = newBillId[:len(newBillId)-1]

	c := calTheBit(newBillId)
	controlInt, _ := strconv.Atoi(controlBit)
	result = c == controlInt

	billType := GetBillType(billParams)

	return result && billType != "unknown"
}

// GetBarCode builds the payment barcode from the bill id and payment id.
func GetBarCode(billParams BillParams) string {
	billID := strconv.Itoa(billParams.BillId)
	paymentID := strconv.Itoa(billParams.PaymentId)
	return billID + "000" + paymentID
}

// GetCurrency returns the bill amount, scaled to Toman or Rial per the
// Currency flags (Rial, or an unset Currency, uses the 1000 multiplier).
func GetCurrency(billParams BillParams) int {
	currency := 100
	if billParams.Currency.Rial || !billParams.Currency.Toman {
		currency = 1000
	}

	payment := strconv.Itoa(billParams.PaymentId)
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
