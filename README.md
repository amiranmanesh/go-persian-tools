<div align="center">
	<img src="https://raw.githubusercontent.com/persian-tools/persian-tools/master/images/logo.png" width="180" alt="Persian Tools" />
	<h1>Go Persian Tools</h1>
	<p><em>An anthology of tools for working with Persian (Iranian) data in Go.</em></p>
</div>

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/amiranmanesh/go-persian-tools.svg)](https://pkg.go.dev/github.com/amiranmanesh/go-persian-tools)
[![CI](https://github.com/amiranmanesh/go-persian-tools/actions/workflows/ci.yml/badge.svg)](https://github.com/amiranmanesh/go-persian-tools/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/amiranmanesh/go-persian-tools/branch/master/graph/badge.svg)](https://codecov.io/gh/amiranmanesh/go-persian-tools)
[![Go Report Card](https://goreportcard.com/badge/github.com/amiranmanesh/go-persian-tools)](https://goreportcard.com/report/github.com/amiranmanesh/go-persian-tools)
[![Go Version](https://img.shields.io/github/go-mod/go-version/amiranmanesh/go-persian-tools)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

## Features

| Package | What it does |
| --- | --- |
| [`bank`](./bank) | Validate card numbers (Luhn), resolve the issuing bank, and validate Sheba (IBAN) codes with bank details. |
| [`bill`](./bill) | Determine a utility bill's type, amount and barcode, and validate its id. |
| [`digit`](./digit) | Add/remove thousands separators and convert digits to Persian words. |
| [`nationalid`](./nationalid) | Validate national numbers (code-e Melli) and resolve the issuing city/province. |
| [`phonenumbers`](./phonenumbers) | Validate Iranian mobile numbers, normalize prefixes and resolve operator details. |

## Install

```bash
go get github.com/amiranmanesh/go-persian-tools@latest
```

Requires **Go 1.23+**. Import only the sub-packages you need.

## Usage

### bank

```go
import "github.com/amiranmanesh/go-persian-tools/bank"

// Card number -> bank slug
name, err := bank.CardInfo("6037701689095443") // "keshavarzi", nil
_, err = bank.CardInfo("6219861034529008")      // "", bank.ErrInvalidCard

// Sheba (IBAN)
sheba := bank.ShebaCode{Code: "IR820540102680020817909002"}
if ok := sheba.IsValid(); ok {
    info := sheba.IsSheba()
    fmt.Println(info.Name)        // Parsian Bank
    fmt.Println(info.PersianName) // بانک پارسیان
}
```

`IsSheba` returns a `bank.ShebaResult`:

```go
type ShebaResult struct {
    Name                   string
    Code                   string
    NickName               string
    PersianName            string
    AccountNumber          string
    AccountNumberAvailable bool
    FormattedAccountNumber string
    Process                func(str string) ShebaProcess
}
```

### bill

```go
import "github.com/amiranmanesh/go-persian-tools/bill"

params := bill.BillParams{BillId: 1117753200140, PaymentId: 12070160}

bill.GetBillType(params)  // "تلفن ثابت"
bill.GetCurrency(params)  // 120000
bill.GetBarCode(params)   // "111775320014000012070160"
bill.VerifyBillID(params) // true
```

### digit

```go
import "github.com/amiranmanesh/go-persian-tools/digit"

digit.DigitToWord("۱۵۶۷۸۹") // "صد پنجاه و شش هزار هفتصد هشتاد و نه"
digit.DigitToWord("-10")     // "منفی ده"

digit.AddCommas(14555478854)          // "14,555,478,854"
n, err := digit.RemoveCommas("4,555") // 4555, nil
```

### nationalid

```go
import "github.com/amiranmanesh/go-persian-tools/nationalid"

nationalid.Validate("0067749828") // true
nationalid.Validate("0684159415") // false

place := nationalid.GetPlaceByIranNationalId("0499370899")
fmt.Println(place.City, place.Province) // شهرری تهران
```

### phonenumbers

```go
import "github.com/amiranmanesh/go-persian-tools/phonenumbers"

phonenumbers.IsPhoneValid("09122221811") // true

details, err := phonenumbers.GetPhoneDetails("09123456789")
if err == nil {
    fmt.Println(details.GetOperator())     // MCI
    fmt.Println(details.GetProvinceList()) // [...]
}

// Normalize the prefix (e.g. to +98)
phonenumbers.PhoneNumberNormalizer("09122221811", "+98") // "+989122221811", nil
```

A runnable end-to-end demo lives in [`examples/`](./examples):

```bash
go run ./examples
```

## Development

```bash
make test    # run tests
make race    # tests with the race detector
make cover   # tests + coverage profile
make lint    # golangci-lint (https://golangci-lint.run)
make example # run the example program
```

## Contributing

Contributions are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md). Please keep
`gofmt`, `go vet` and the tests green, and add tests for new behavior.

## License

Released under the [MIT License](LICENSE).
