# Security Policy

## Supported versions

The latest minor release is supported. Fixes are published as a new release
rather than backported.

## Scope

This package takes a string and returns a string. It has no dependencies, does
no I/O, spawns nothing and parses no untrusted format, so the realistic risk is
narrow: a crafted input that makes a function panic, loop without terminating,
or allocate out of proportion to its input.

Reports about any of those are welcome. So is anything about normalization
being used as a security boundary — for example, two different inputs that
`Normalize` folds onto the same string in a way that could defeat a uniqueness
check on usernames. That behavior is by design for search and sorting, but if
you have found a case where it is surprising enough to be dangerous, please say
so.

## Reporting

Use [GitHub's private vulnerability reporting](https://github.com/amiranmanesh/go-persian-tools/security/advisories/new).
Please do not open a public issue first.

Include the input that triggers it, the version, and what happens. You should
get a reply within a week.

## Verifying what you depend on

Releases are tagged and served through the Go module proxy, so the checksum
database pins what you get:

```sh
go get github.com/amiranmanesh/go-persian-tools@latest
go mod verify
```
