// Package nationalid validates Iranian national numbers (code-e Melli) and
// resolves the issuing city and province from a national id.
package nationalid

import (
	"regexp"
	"strconv"
)

// nationalIDPattern matches exactly 10 ASCII digits.
var nationalIDPattern = regexp.MustCompile(`^\d{10}$`)

var notAllowedDigits = map[string]bool{
	"0000000000": true,
	"2222222222": true,
	"3333333333": true,
	"4444444444": true,
	"5555555555": true,
	"6666666666": true,
	"7777777777": true,
	"8888888888": true,
	"9999999999": true,
}

// Validate reports whether code is a valid Iranian national number
// (code-e Melli): 10 digits passing the official check-digit algorithm.
func Validate(code string) bool {
	if notAllowedDigits[code] {
		return false
	}
	if !nationalIDPattern.MatchString(code) {
		return false
	}
	code = ("0000" + code)[len(code)+4-10:]
	stringCode, _ := strconv.Atoi(code[3:6])
	if stringCode == 0 {
		return false
	}
	lastNumber, _ := strconv.Atoi(code[9:10])
	sum := 0

	for i := 0; i < 9; i++ {
		intCode, _ := strconv.Atoi(code[i : i+1])
		sum += intCode * (10 - i)
	}
	sum = sum % 11
	result := false
	if sum < 2 && lastNumber == sum || sum >= 2 && lastNumber == 11-sum {
		result = true
	}
	return result
}
