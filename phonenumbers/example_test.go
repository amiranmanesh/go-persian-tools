package phonenumbers_test

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/phonenumbers"
)

func ExampleIsPhoneValid() {
	fmt.Println(phonenumbers.IsPhoneValid("09122221811"))
	// Output: true
}

func ExampleGetPhoneDetails() {
	details, err := phonenumbers.GetPhoneDetails("09123456789")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(details.GetOperator())
	// Output: MCI
}
