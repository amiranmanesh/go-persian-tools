// Package phonenumbers validates Iranian mobile numbers, normalizes their
// prefix and resolves operator details from the number's prefix.
package phonenumbers

import "errors"

// Operator identifies an Iranian mobile network operator.
type Operator string

// SimType tells a permanent (post-paid) SIM from a credit (pre-paid) one.
type SimType string

// Errors returned by the phone number helpers.
var (
	// ErrInvalidFormat is returned when a phone number is not a valid Iranian
	// mobile number.
	ErrInvalidFormat = errors.New("phonenumbers: invalid phone number format")
	// ErrInvalidPrefix is returned when an operator prefix is not recognized.
	ErrInvalidPrefix = errors.New("phonenumbers: invalid prefix")
)

// Prefixes lists the dialing prefixes an Iranian mobile number may carry.
var Prefixes = []string{"+98", "98", "0098", "0"}

// The operators and SIM types this package recognizes.
const (
	ShatelMobile Operator = "ShatelMobile"
	MCI          Operator = "MCI"
	Irancell     Operator = "Irancell"
	Taliya       Operator = "Taliya"
	RightTel     Operator = "RightTel"

	Permanent SimType = "Permanent"
	Credit    SimType = "Credit"
)

// OperatorDetails describes the operator behind a number prefix: its base
// province, the provinces it covers, the SIM types it issues and any
// sub-model the prefix is reserved for.
type OperatorDetails struct {
	base     string
	province []string
	model    string
	operator Operator
	simTypes []SimType
}

// Details returns the prefix table of the operator, or nil when the operator
// is not recognized.
func (o Operator) Details() map[string]OperatorDetails {
	switch o {
	case MCI:
		return MCIMap
	case Taliya:
		return TALIYA
	case RightTel:
		return RIGHTTEL
	case Irancell:
		return IRANCELL
	case ShatelMobile:
		return SHATELMOBILE
	default:
		return nil
	}
}

// GetProvinceList returns the provinces the prefix covers beyond its base.
func (od *OperatorDetails) GetProvinceList() []string {
	return od.province
}

// GetBase returns the base province of the prefix.
func (od *OperatorDetails) GetBase() string {
	return od.base
}

// GetModel returns the sub-model the prefix is reserved for, if any.
func (od *OperatorDetails) GetModel() string {
	return od.model
}

// GetOperator returns the operator the prefix belongs to.
func (od *OperatorDetails) GetOperator() Operator {
	return od.operator
}

// GetSimTypeList returns the SIM types issued on the prefix.
func (od *OperatorDetails) GetSimTypeList() []SimType {
	return od.simTypes
}
