// Package phonenumbers validates Iranian mobile numbers, normalizes their
// prefix and resolves operator details from the number's prefix.
package phonenumbers

import "errors"

type Operator string
type SimType string

// Errors returned by the phone number helpers.
var (
	// ErrInvalidFormat is returned when a phone number is not a valid Iranian
	// mobile number.
	ErrInvalidFormat = errors.New("phonenumbers: invalid phone number format")
	// ErrInvalidPrefix is returned when an operator prefix is not recognized.
	ErrInvalidPrefix = errors.New("phonenumbers: invalid prefix")
)

// List of valid prefixes
var Prefixes = []string{"+98", "98", "0098", "0"}

const (
	ShatelMobile Operator = "ShatelMobile"
	MCI          Operator = "MCI"
	Irancell     Operator = "Irancell"
	Taliya       Operator = "Taliya"
	RightTel     Operator = "RightTel"

	Permanent SimType = "Permanent"
	Credit    SimType = "Credit"
)

type OperatorDetails struct {
	base     string
	province []string
	model    string
	operator Operator
	simTypes []SimType
}

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

func (od *OperatorDetails) GetProvinceList() []string {
	return od.province
}

func (od *OperatorDetails) GetBase() string {
	return od.base
}

func (od *OperatorDetails) GetModel() string {
	return od.model
}

func (od *OperatorDetails) GetOperator() Operator {
	return od.operator
}

func (od *OperatorDetails) GetSimTypeList() []SimType {
	return od.simTypes
}
