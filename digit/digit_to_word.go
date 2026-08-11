package digit

import (
	"strconv"

	"github.com/mavihq/persian"
	ntw "moul.io/number-to-words"
)

// DigitToWord converts a numeric string (Persian or ASCII digits, optionally
// negative) into its Persian word representation, e.g. "156789" becomes
// "صد پنجاه و شش هزار هفتصد هشتاد و نه". It returns an empty string when the
// input is not a valid integer.
func DigitToWord(word string) string {
	n, err := strconv.Atoi(persian.ToEnglishDigits(word))
	if err != nil {
		return ""
	}
	return ntw.IntegerToIrIr(n)
}
