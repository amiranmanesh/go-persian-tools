# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Rewrote `text.Finglish`.** It now works in three passes — split into words,
  read the letters into phonemes, then group the phonemes into syllables —
  instead of walking a single state machine across the whole input. That makes
  the rules expressible: long vowels (خوب is "khub", not "khob"), the silent vav
  in خوا (خواهر is "khahar"), a final ه as a vowel (خانه is "khane"), the
  initial ای and او pairs (ایران is "iran"), و and ی resolved from context
  (اهواز is "ahvaz"), and a vowel in every syllable (رفتن is "raftan", not
  "raftn").

  Measured against a reference word list, the hit rate goes from 47.0% to 65.2%
  on the words the rules were developed against, and from 38.3% to 53.3% on a
  held-out list that was never used while tuning them. Both lists ship with the
  tests, and `TestFinglishAccuracy` fails if either rate drops.

  Words no longer leak state across a space, so a word transliterates the same
  way in a sentence as it does alone. Roughly 1.9x slower per call (261ns
  against 139ns) for the extra pass.

## [1.1.0] - 2026-08-11

This release merges `github.com/amiranmanesh/persian` into this module and
removes every external dependency.

### Added

- **`text` package** — Persian text utilities merged in from
  `github.com/amiranmanesh/persian`: `FixArabic` and `Normalize` for folding
  Arabic look-alikes, `SwitchToPersianKey` / `SwitchToEnglishKey` for keyboard
  layout recovery, `Finglish` for romanization, plus `Reverse`,
  `CheckIsEnglish`, `OnlyPersianAlpha` and the `ZWNJ` / `ZWJ` constants.
- **`digit` package additions** — conversion between the ASCII, Persian and
  Arabic-Indic digit sets (`ToPersianDigits`, `ToEnglishDigits`,
  `ToPersianDigitsFromInt`), the `OnlyNumbers` family of filters, Persian money
  formatting (`Currency`, `Toman`, `Rial`), and `ToWords` / `ParseInt`.
- **`persian-tools` CLI** — one binary over every package, usable as a one-off
  or as a pipeline filter, published as prebuilt binaries for six platforms and
  as a container image.
- Three virtual mobile operators and a new prefix: ApTel (`999`), Tele Kish
  (`934`), Espadan (`931`) and RightTel's `924`.
- `ShebaResult.MergedInto`, marking the five institutions absorbed by Bank Sepah
  between 2019 and 2020. Their codes remain in the table because IBANs issued
  before the merger still have to resolve.
- Card prefixes `505426` (Gardeshgari) and `636797` (Central Bank).
- Eleven fuzz targets, benchmarks for every hot path, and dataset tests that
  check all 477 national id prefixes resolve and that no mobile prefix is
  claimed twice. Coverage is 96.4%.

### Changed

- **The module has no dependencies.** `moul.io/number-to-words`,
  `github.com/dustin/go-humanize` and `github.com/mavihq/persian` are replaced
  by implementations in this repository, and `go.sum` is gone.
- The minimum Go version is 1.22. CI now tests Go 1.22 through 1.26 on Linux,
  macOS and Windows, and runs golangci-lint, fuzzing and govulncheck.
- Renamed to Go's naming conventions, with the old names kept as deprecated
  aliases: `bill.BillParams` → `bill.Params`, `digit.DigitToWord` →
  `digit.ToWord`, `nationalid.IPlaceByNationalId` → `nationalid.Place`,
  `nationalid.GetPlaceByIranNationalId` → `GetPlaceByIranNationalID`.
- `digit.RemoveCommas` now also accepts Persian digits and the Persian
  thousands separator, and is roughly four times faster.

### Fixed

- **`DigitToWord` now emits grammatically correct Persian.** 156789 reads
  `صد و پنجاه و شش هزار و هفتصد و هشتاد و نه`; the previous library dropped the
  `و` conjunctions.
- `GetPhoneDetails` rejected a valid bare number such as `9123456789`, because
  `GetOperatorPrefix` insisted on a dialing prefix that `IsPhoneValid` does not
  require.
- Six national id rows had been appended to the table under Isfahan's province
  code. Zanjan, Abhar and Khorramdarreh now resolve to a new Zanjan province
  entry, and Marand, Malekan and Mianeh to East Azerbaijan.
- Azna's national id code carried a leading space, which leaked into the
  returned `Codes` slice.

### Removed

- `examples/`, whose demo is covered by the CLI and the runnable `Example`
  functions.

### Deprecated

- `bill.BillParams`, `digit.DigitToWord`, `nationalid.IPlaceByNationalId` and
  `nationalid.GetPlaceByIranNationalId`. They still work; prefer the new names.

### Breaking

- The `BillId` and `PaymentId` fields of `bill.Params` are now `BillID` and
  `PaymentID`. Struct fields cannot be aliased, so this is the one rename that
  is not source-compatible.
- `digit.DigitToWord` output changed, as described under Fixed.

## [1.0.1] - 2026-08-11

### Added

- Publish a multi-arch container image of the demo to the GitHub Container
  Registry (`ghcr.io/amiranmanesh/go-persian-tools`) on every release, built
  from a distroless base and pinned action SHAs.

## [1.0.0] - 2026-08-11

### Security

- **Removed a malicious GitHub Actions workflow** (`.github/workflows/go.yml`)
  that base64-decoded and executed a credential-exfiltration script targeting
  cloud metadata, CI/CD tokens and local secrets. If you ever ran this workflow,
  rotate any secrets exposed to the repository's Actions.

### Changed

- **Module path** is now `github.com/amiranmanesh/go-persian-tools`.
- Renamed packages to idiomatic Go: `national_id` → `nationalid`, and the
  `phone_numbers` directory now matches its `phonenumbers` package.
- `digit.RemoveCommas` now returns `(int64, error)` instead of `int`, avoiding
  overflow on large values and surfacing parse errors.
- `bank`'s exported result type is now `ShebaResult` (was the unexported
  `shebaResultHash`); added `ShebaCode.IsValid()`.
- Introduced sentinel errors: `bank.ErrInvalidCard`, `bank.ErrBankNotFound`,
  `phonenumbers.ErrInvalidPrefix`.

### Added

- Doc comments across all exported identifiers and package overviews.
- Unified CI (lint + multi-version test + coverage) and a tag-based release
  workflow; `.golangci.yml` configuration.
- `LICENSE` (MIT), `CONTRIBUTING.md` and this changelog.

### Removed

- Dockerfile, `.docker_push`, CircleCI config and the demo binary at the module
  root (this is a library, not an application).
