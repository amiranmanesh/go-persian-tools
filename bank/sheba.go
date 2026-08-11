package bank

import (
	"errors"
	"math/big"
	"regexp"
	"strconv"
)

// Sheba (IBAN) validation errors, wrapped by iso7064Mod97_10.
var (
	errShebaCheckDigits = errors.New("bank: IBAN has incorrect check digits")
	errShebaChecksum    = errors.New("bank: IBAN checksum is not valid")
)

var (
	// shebaPattern matches the overall structure of an Iranian IBAN and
	// captures the 3-digit bank code in the second group.
	shebaPattern = regexp.MustCompile(`(IR[0-9]{2}([0-9]{3})[0-9]{19})`)
	// shebaShapePattern matches a well-formed 26-character Iranian IBAN.
	shebaShapePattern = regexp.MustCompile(`(IR[0-9]{24})`)
)

// ShebaCode wraps an Iranian IBAN ("Sheba") string.
type ShebaCode struct {
	Code string
}

// IsSheba validates the IBAN and, when valid, returns the matching bank's
// details. It returns a zero ShebaResult when the code is invalid or the bank
// is unknown. Use IsValid when only a boolean is needed.
func (s ShebaCode) IsSheba() ShebaResult {
	if !s.IsValid() {
		return ShebaResult{}
	}

	code := shebaPattern.FindStringSubmatch(s.Code)
	if code == nil {
		return ShebaResult{}
	}

	bank := shebaHashTable(code[2])
	if bank.Name == "" {
		return ShebaResult{}
	}

	return bank
}

// IsValid reports whether the code is a structurally valid Iranian IBAN that
// passes the ISO 7064 MOD-97-10 checksum.
func (s ShebaCode) IsValid() bool {
	shebaCode := s.Code

	if len(shebaCode) != 26 || !shebaShapePattern.MatchString(shebaCode) {
		return false
	}

	d1 := []rune(shebaCode)[0] - 65 + 10
	d2 := []rune(shebaCode)[1] - 65 + 10

	rearranged := shebaCode[4:] + strconv.Itoa(int(d1)) + strconv.Itoa(int(d2)) + shebaCode[2:4]

	return iso7064Mod97_10(rearranged) == nil
}

// iso7064Mod97_10 validates the IBAN checksum per ISO 7064 MOD-97-10.
func iso7064Mod97_10(iban string) error {
	bigVal, ok := new(big.Int).SetString(iban, 10)
	if !ok {
		return errShebaCheckDigits
	}
	if new(big.Int).Mod(bigVal, big.NewInt(97)).Int64() != 1 {
		return errShebaChecksum
	}
	return nil
}
