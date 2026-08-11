# Contributing

Thanks for your interest in improving **go-persian-tools**!

## Getting started

1. Fork and clone the repository.
2. Make sure you have **Go 1.23+** installed.
3. Install [golangci-lint](https://golangci-lint.run/welcome/install/) for linting.

## Development workflow

```bash
make test    # run the test suite
make race    # run tests with the race detector
make cover   # generate a coverage profile
make lint    # run golangci-lint
make fmt     # format the code
```

Before opening a pull request, please make sure:

- `gofmt`, `go vet ./...` and `make test` are clean.
- New behavior is covered by tests (unit tests and, where it aids docs,
  testable `Example` functions).
- Exported identifiers have doc comments.

## Commit messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/)
(`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:` …). Renovate keeps
dependencies up to date via automated PRs.

## Releases

Releases are cut by pushing a semantic-version tag:

```bash
git tag v1.2.3
git push origin v1.2.3
```

The `Release` workflow runs the tests and publishes a GitHub Release with
auto-generated notes.

## Reporting security issues

Please do **not** open a public issue for security problems. Report them
privately to the maintainers instead.
