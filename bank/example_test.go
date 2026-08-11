package bank_test

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/bank"
)

func ExampleCardInfo() {
	name, err := bank.CardInfo("6037701689095443")
	fmt.Println(name, err)
	// Output: keshavarzi <nil>
}

func ExampleShebaCode_IsSheba() {
	result := bank.ShebaCode{Code: "IR820540102680020817909002"}.IsSheba()
	fmt.Println(result.Name)
	// Output: Parsian Bank
}
