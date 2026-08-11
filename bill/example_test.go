package bill_test

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/bill"
)

func ExampleGetBillType() {
	params := bill.Params{BillID: 1117753200140, PaymentID: 12070160}
	fmt.Println(bill.GetBillType(params))
	// Output: تلفن ثابت
}

func ExampleGetCurrency() {
	params := bill.Params{BillID: 1117753200140, PaymentID: 12070160}
	fmt.Println(bill.GetCurrency(params))
	// Output: 120000
}
