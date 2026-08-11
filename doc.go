// Package persiantools is an anthology of tools for working with Persian
// (Iranian) data in Go.
//
// The module has no dependencies. It is organized as a set of focused
// sub-packages:
//
//   - bank         validate card numbers, resolve banks, and check Sheba (IBAN) codes
//   - bill         parse and validate Iranian utility bills
//   - digit        convert digit sets, group and spell numbers, format money
//   - nationalid   validate national numbers (code-e Melli) and resolve place
//   - phonenumbers validate and inspect Iranian mobile numbers
//   - text         normalize Persian text, fix keyboard layouts, romanize
//
// Import only the sub-packages you need, for example:
//
//	import "github.com/amiranmanesh/go-persian-tools/nationalid"
//
// A command-line interface over every package is available too:
//
//	go install github.com/amiranmanesh/go-persian-tools/cmd/persian-tools@latest
package persiantools
