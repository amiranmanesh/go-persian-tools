// Command example demonstrates the go-persian-tools sub-packages.
//
// Run it with:
//
//	go run ./examples
package main

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/bank"
	"github.com/amiranmanesh/go-persian-tools/bill"
	"github.com/amiranmanesh/go-persian-tools/digit"
	"github.com/amiranmanesh/go-persian-tools/nationalid"
	"github.com/amiranmanesh/go-persian-tools/phonenumbers"
)

func main() {
	// --- bill ---
	params := bill.BillParams{BillId: 1117753200140, PaymentId: 12070160}
	fmt.Println("Bill type:  ", bill.GetBillType(params))
	fmt.Println("Bill amount:", bill.GetCurrency(params))
	fmt.Println("Bill barcode:", bill.GetBarCode(params))
	fmt.Println("Bill valid: ", bill.VerifyBillID(bill.BillParams{BillId: 7748317800142, PaymentId: 1770160}))

	// --- bank ---
	if name, err := bank.CardInfo("6037701689095443"); err == nil {
		fmt.Println("Card bank:  ", name)
	}
	sheba := bank.ShebaCode{Code: "IR820540102680020817909002"}.IsSheba()
	fmt.Println("Sheba bank: ", sheba.Name, "-", sheba.PersianName)

	// --- digit ---
	fmt.Println("Digit->word:", digit.DigitToWord("۱۵۶۷۸۹"))
	fmt.Println("Add commas: ", digit.AddCommas(14555478854))
	if n, err := digit.RemoveCommas("4,555,522,212"); err == nil {
		fmt.Println("Remove commas:", n)
	}

	// --- nationalid ---
	fmt.Println("National id valid:", nationalid.Validate("0067749828"))
	place := nationalid.GetPlaceByIranNationalId("0499370899")
	fmt.Printf("National id place: %s, %s\n", place.City, place.Province)

	// --- phonenumbers ---
	fmt.Println("Phone valid:", phonenumbers.IsPhoneValid("09122221811"))
	if details, err := phonenumbers.GetPhoneDetails("09123456789"); err == nil {
		fmt.Println("Phone operator:", details.GetOperator())
	}
}
