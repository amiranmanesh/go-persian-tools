package bill_test

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/bill"
)

func ExampleGetBillType() {
	params := bill.BillParams{BillId: 1117753200140, PaymentId: 12070160}
	fmt.Println(bill.GetBillType(params))
	// Output: تلفن ثابت
}

func ExampleGetCurrency() {
	params := bill.BillParams{BillId: 1117753200140, PaymentId: 12070160}
	fmt.Println(bill.GetCurrency(params))
	// Output: 120000
}
