<div align="center">
	<h1>Go Persian Tools</h1>
	<p><em>مجموعه‌ای از ابزارهای کار با داده‌های فارسی و ایرانی در زبان Go</em></p>
	<p>
		<a href="README.md">English</a>
	</p>
</div>

<div align="center">

[![Go Reference](https://pkg.go.dev/badge/github.com/amiranmanesh/go-persian-tools.svg)](https://pkg.go.dev/github.com/amiranmanesh/go-persian-tools)
[![CI](https://github.com/amiranmanesh/go-persian-tools/actions/workflows/ci.yml/badge.svg)](https://github.com/amiranmanesh/go-persian-tools/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/amiranmanesh/go-persian-tools)](https://goreportcard.com/report/github.com/amiranmanesh/go-persian-tools)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

**بدون هیچ وابستگی.** فایل `go.mod` این ماژول هیچ وابستگی‌ای ندارد و اصلاً فایل
`go.sum` وجود ندارد؛ همه‌چیز فقط با کتابخانه‌ی استاندارد Go نوشته شده است.

## پکیج‌ها

| پکیج | کاری که انجام می‌دهد |
| --- | --- |
| [`text`](./text) | یکسان‌سازی متن فارسی، اصلاح حروف عربی، اصلاح چیدمان کیبورد، فینگلیش |
| [`digit`](./digit) | تبدیل ارقام، جداسازی سه‌رقمی، تبدیل عدد به حروف، قالب‌بندی تومان و ریال |
| [`bank`](./bank) | اعتبارسنجی شماره کارت، تشخیص بانک، اعتبارسنجی شبا |
| [`nationalid`](./nationalid) | اعتبارسنجی کد ملی و تشخیص شهر و استان محل صدور |
| [`phonenumbers`](./phonenumbers) | اعتبارسنجی شماره موبایل، یکسان‌سازی پیش‌شماره، تشخیص اپراتور |
| [`bill`](./bill) | تشخیص نوع، مبلغ و بارکد قبض و اعتبارسنجی شناسه‌ی آن |

## نصب

```bash
go get github.com/amiranmanesh/go-persian-tools@latest
```

نیازمند **Go 1.22 به بالا**. فقط همان زیرپکیجی را import کنید که لازم دارید.

## نمونه‌ی استفاده

### text — متن

```go
import "github.com/amiranmanesh/go-persian-tools/text"

// اصلاح حروف عربی تا مقایسه‌ی متن درست کار کند
text.FixArabic("علي كريم")  // علی کریم
text.Normalize("مقالهٔ من") // مقاله من

// بازیابی متنی که با چیدمان اشتباه کیبورد تایپ شده
text.SwitchToPersianKey("sghl") // سلام
text.SwitchToEnglishKey("اثغ")  // hey

text.Finglish("سلام")               // salam
text.Reverse("سلام")                // مالس
text.CheckIsEnglish("ali")          // true
text.OnlyPersianAlpha("123شاهینhi") // شاهین
```

`Normalize` برای ساختن کلید مقایسه است، نه برای نمایش: شکل‌های مختلف الف و همزه
را یکی می‌کند، اعراب را حذف می‌کند و نیم‌فاصله را به فاصله تبدیل می‌کند. متن
اصلی را برای نمایش به کاربر نگه دارید.

### digit — ارقام و مبلغ

```go
import "github.com/amiranmanesh/go-persian-tools/digit"

digit.ToPersianDigits("123salam456") // ۱۲۳salam۴۵۶
digit.ToEnglishDigits("۰۹۱۲٣٤٥٦٧٨٩") // 09123456789
digit.OnlyNumbers("شماره: 0912-345")  // 0912345

digit.AddCommas(14555478854)            // 14,555,478,854
n, err := digit.RemoveCommas("۱۲۳،۴۵۶") // 123456
digit.ToWords(156789)                   // صد و پنجاه و شش هزار و هفتصد و هشتاد و نه
digit.ToWord("-10")                     // منفی ده

digit.Currency("1234567") // ۱،۲۳۴،۵۶۷
digit.Toman("1234567")    // ۱،۲۳۴،۵۶۷ تومان
digit.Rial("1234567")     // ۱،۲۳۴،۵۶۷ ﷼
```

### bank — بانک

```go
import "github.com/amiranmanesh/go-persian-tools/bank"

name, err := bank.CardInfo("6037701689095443") // "keshavarzi"
// خطاها: bank.ErrInvalidCard و bank.ErrBankNotFound

sheba := bank.ShebaCode{Code: "IR820540102680020817909002"}
if sheba.IsValid() {
    info := sheba.IsSheba()
    fmt.Println(info.PersianName) // بانک پارسیان
}
```

بانک‌هایی که در بانک سپه ادغام شده‌اند (انصار، قوامین، حکمت ایرانیان، مهر
اقتصاد و مؤسسه‌ی کوثر) همچنان شناسایی می‌شوند، چون شباهای صادرشده پیش از ادغام
هنوز معتبرند؛ این موارد با فیلد `MergedInto` مشخص شده‌اند.

### nationalid — کد ملی

```go
import "github.com/amiranmanesh/go-persian-tools/nationalid"

nationalid.Validate("0067749828") // true

place := nationalid.GetPlaceByIranNationalID("0499370899")
fmt.Println(place.City, place.Province) // شهرری تهران
```

### phonenumbers — شماره موبایل

```go
import "github.com/amiranmanesh/go-persian-tools/phonenumbers"

phonenumbers.IsPhoneValid("09122221811") // true

details, err := phonenumbers.GetPhoneDetails("09123456789")
if err == nil {
    fmt.Println(details.GetOperator()) // MCI
}

phonenumbers.PhoneNumberNormalizer("09122221811", "+98") // +989122221811
```

اپراتورهای همراه اول، ایرانسل، رایتل، تالیا، شاتل‌موبایل، آپتل، تله‌کیش و
اسپادان پشتیبانی می‌شوند.

### bill — قبض

```go
import "github.com/amiranmanesh/go-persian-tools/bill"

params := bill.Params{BillID: 1117753200140, PaymentID: 12070160}

bill.GetBillType(params)  // تلفن ثابت
bill.GetCurrency(params)  // 120000
bill.GetBarCode(params)   // 111775320014000012070160
bill.VerifyBillID(params) // true
```

## خط فرمان

همه‌ی پکیج‌ها از طریق یک باینری هم در دسترس‌اند و می‌توانید آن را در خط لوله
(pipeline) هم استفاده کنید:

```bash
go install github.com/amiranmanesh/go-persian-tools/cmd/persian-tools@latest

persian-tools normalize "علي كريم"          # علی کریم
persian-tools words 156789                   # صد و پنجاه و شش هزار و هفتصد و هشتاد و نه
persian-tools currency -unit toman 1234567   # ۱،۲۳۴،۵۶۷ تومان
persian-tools key-to-persian sghl            # سلام
persian-tools national-id 0067749828         # true

cat names.txt | persian-tools normalize > keys.txt
```

یا بدون نصب:

```bash
docker run --rm ghcr.io/amiranmanesh/go-persian-tools:latest normalize "علي كريم"
```

## توسعه

```bash
make check   # هر چیزی که CI اجرا می‌کند
make test    # تست‌ها با race detector
make cover   # گزارش پوشش تست
make bench   # بنچمارک‌ها
make fuzz    # تست‌های fuzz
make help    # فهرست همه‌ی دستورها
```

## مشارکت

راهنمای مشارکت در [CONTRIBUTING.md](CONTRIBUTING.md) است. برای گزارش مشکل
امنیتی [SECURITY.md](SECURITY.md) را ببینید.

## مجوز

تحت مجوز [MIT](LICENSE) منتشر شده است.
