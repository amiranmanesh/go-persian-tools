package digit_test

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/digit"
)

func ExampleDigitToWord() {
	fmt.Println(digit.DigitToWord("156789"))
	// Output: صد پنجاه و شش هزار هفتصد هشتاد و نه
}

func ExampleAddCommas() {
	fmt.Println(digit.AddCommas(14555478854))
	// Output: 14,555,478,854
}

func ExampleRemoveCommas() {
	n, _ := digit.RemoveCommas("4,555,478,854")
	fmt.Println(n)
	// Output: 4555478854
}
