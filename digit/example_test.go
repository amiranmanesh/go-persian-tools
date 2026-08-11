package digit_test

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/digit"
)

func ExampleToWord() {
	fmt.Println(digit.ToWord("156789"))
	// Output: صد و پنجاه و شش هزار و هفتصد و هشتاد و نه
}

func ExampleToWords() {
	fmt.Println(digit.ToWords(1402))
	// Output: هزار و چهارصد و دو
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

func ExampleToPersianDigits() {
	fmt.Println(digit.ToPersianDigits("123salam456"))
	// Output: ۱۲۳salam۴۵۶
}

func ExampleToPersianDigitsFromInt() {
	fmt.Println(digit.ToPersianDigitsFromInt(1402))
	// Output: ۱۴۰۲
}

func ExampleToEnglishDigits() {
	// Persian and Arabic-Indic digits are both accepted.
	fmt.Println(digit.ToEnglishDigits("۰۹۱۲٣٤٥٦٧٨٩"))
	// Output: 09123456789
}

func ExampleOnlyEnglishNumbers() {
	fmt.Println(digit.OnlyEnglishNumbers("شماره: 0912-345-6789"))
	// Output: 09123456789
}

func ExampleCurrency() {
	fmt.Println(digit.Currency("1234567"))
	// Output: ۱،۲۳۴،۵۶۷
}

func ExampleToman() {
	fmt.Println(digit.Toman("1234567"))
	// Output: ۱،۲۳۴،۵۶۷ تومان
}

func ExampleRial() {
	fmt.Println(digit.Rial("1234567"))
	// Output: ۱،۲۳۴،۵۶۷ ﷼
}
