<div align="center">
	<h1>Go Persian Tools</h1>
	<p><em>An anthology of tools for working with Persian (Iranian) data in Go.</em></p>
	<p>
		<a href="README.fa.md">فارسی</a>
	</p>
</div>

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/amiranmanesh/go-persian-tools.svg)](https://pkg.go.dev/github.com/amiranmanesh/go-persian-tools)
[![CI](https://github.com/amiranmanesh/go-persian-tools/actions/workflows/ci.yml/badge.svg)](https://github.com/amiranmanesh/go-persian-tools/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/amiranmanesh/go-persian-tools)](https://goreportcard.com/report/github.com/amiranmanesh/go-persian-tools)
[![Go Version](https://img.shields.io/github/go-mod/go-version/amiranmanesh/go-persian-tools)](go.mod)
[![Dependencies](https://img.shields.io/badge/dependencies-none-0f9d94)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

**No dependencies.** The module's `go.mod` lists nothing, and there is no `go.sum`.
Everything here is the standard library and this repository.

## Packages

| Package | What it does |
| --- | --- |
| [`text`](./text) | Normalize Persian text, fold Arabic look-alikes, fix keyboard layouts, romanize |
| [`digit`](./digit) | Convert digit sets, group and spell numbers, format Toman and Rial |
| [`bank`](./bank) | Validate card numbers (Luhn), resolve banks, validate Sheba (IBAN) codes |
| [`nationalid`](./nationalid) | Validate national numbers (code-e Melli), resolve city and province |
| [`phonenumbers`](./phonenumbers) | Validate Iranian mobiles, normalize prefixes, resolve operators |
| [`bill`](./bill) | Determine a utility bill's type, amount and barcode; validate its id |

## Install

```bash
go get github.com/amiranmanesh/go-persian-tools@latest
```

Requires **Go 1.22+**. Import only the sub-packages you need.

## Usage

### text

```go
import "github.com/amiranmanesh/go-persian-tools/text"

// Fold Arabic look-alikes so text compares reliably.
text.FixArabic("علي كريم")  // علی کریم
text.Normalize("مقالهٔ من") // مقاله من  (idempotent: safe to apply twice)

// Recover text typed on the wrong keyboard layout.
text.SwitchToPersianKey("sghl") // سلام
text.SwitchToEnglishKey("اثغ")  // hey

text.Finglish("سلام")           // salam
text.Reverse("سلام")            // مالس
text.CheckIsEnglish("ali")      // true
text.OnlyPersianAlpha("123شاهینhi") // شاهین
```

`Normalize` is meant for comparison keys, not display: it folds alef and hamza
variants, drops vocalization marks, and turns zero-width joiners into spaces.
Keep the original for showing back to the user.

### digit

```go
import "github.com/amiranmanesh/go-persian-tools/digit"

// Digit sets — Persian, Arabic-Indic and ASCII.
digit.ToPersianDigits("123salam456")   // ۱۲۳salam۴۵۶
digit.ToEnglishDigits("۰۹۱۲٣٤٥٦٧٨٩")   // 09123456789
digit.OnlyNumbers("شماره: 0912-345")   // 0912345

// Grouping and words.
digit.AddCommas(14555478854)            // 14,555,478,854
n, err := digit.RemoveCommas("۱۲۳،۴۵۶") // 123456
digit.ToWords(156789)                   // صد و پنجاه و شش هزار و هفتصد و هشتاد و نه
digit.ToWord("-10")                     // منفی ده

// Money.
digit.Currency("1234567") // ۱،۲۳۴،۵۶۷
digit.Toman("1234567")    // ۱،۲۳۴،۵۶۷ تومان
digit.Rial("1234567")     // ۱،۲۳۴،۵۶۷ ﷼
```

### bank

```go
import "github.com/amiranmanesh/go-persian-tools/bank"

name, err := bank.CardInfo("6037701689095443") // "keshavarzi", nil
// errors.Is(err, bank.ErrInvalidCard) / bank.ErrBankNotFound

sheba := bank.ShebaCode{Code: "IR820540102680020817909002"}
if sheba.IsValid() {
    info := sheba.IsSheba()
    fmt.Println(info.Name, info.PersianName) // Parsian Bank بانک پارسیان
}
```

### nationalid

```go
import "github.com/amiranmanesh/go-persian-tools/nationalid"

nationalid.Validate("0067749828") // true

place := nationalid.GetPlaceByIranNationalID("0499370899")
fmt.Println(place.City, place.Province) // شهرری تهران
```

### phonenumbers

```go
import "github.com/amiranmanesh/go-persian-tools/phonenumbers"

phonenumbers.IsPhoneValid("09122221811") // true

details, err := phonenumbers.GetPhoneDetails("09123456789")
if err == nil {
    fmt.Println(details.GetOperator()) // MCI
}

phonenumbers.PhoneNumberNormalizer("09122221811", "+98") // +989122221811
```

### bill

```go
import "github.com/amiranmanesh/go-persian-tools/bill"

params := bill.Params{BillID: 1117753200140, PaymentID: 12070160}

bill.GetBillType(params)  // تلفن ثابت
bill.GetCurrency(params)  // 120000
bill.GetBarCode(params)   // 111775320014000012070160
bill.VerifyBillID(params) // true
```

## Command line

Every package is also reachable from one binary, usable as a one-off or as a
pipeline filter:

```bash
go install github.com/amiranmanesh/go-persian-tools/cmd/persian-tools@latest

persian-tools normalize "علي كريم"          # علی کریم
persian-tools words 156789                   # صد و پنجاه و شش هزار و هفتصد و هشتاد و نه
persian-tools currency -unit toman 1234567   # ۱،۲۳۴،۵۶۷ تومان
persian-tools key-to-persian sghl            # سلام
persian-tools national-id 0067749828         # true   (exit 0; invalid exits 1)
persian-tools sheba IR820540102680020817909002

cat names.txt | persian-tools normalize > keys.txt
```

Or without installing anything:

```bash
docker run --rm ghcr.io/amiranmanesh/go-persian-tools:latest normalize "علي كريم"
```

Prebuilt binaries for Linux, macOS and Windows are attached to every
[release](https://github.com/amiranmanesh/go-persian-tools/releases).

## Development

```bash
make check   # fmt, vet, lint and test — everything CI runs
make test    # tests with the race detector
make cover   # coverage profile and HTML report
make bench   # benchmarks
make fuzz    # every fuzz target
make help    # list all targets
```

## Contributing

Contributions are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md) and the
[code of conduct](CODE_OF_CONDUCT.md). Keep `gofmt`, `go vet` and the tests
green, and add tests for new behavior. To report a security issue, see
[SECURITY.md](SECURITY.md).

## License

Released under the [MIT License](LICENSE).
