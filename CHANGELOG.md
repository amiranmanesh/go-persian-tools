# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- Bumped the minimum Go version to 1.23.

### Added

- Doc comments across all exported identifiers and package overviews.
- Runnable demo under `examples/` and testable `Example` functions.
- Unified CI (lint + multi-version test + coverage) and a tag-based release
  workflow; `.golangci.yml` configuration.
- `LICENSE` (MIT), `CONTRIBUTING.md` and this changelog.

### Removed

- Dockerfile, `.docker_push`, CircleCI config and the demo binary at the module
  root (this is a library, not an application).
