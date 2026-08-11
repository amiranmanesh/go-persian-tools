// Package persiantools is an anthology of tools for working with Persian
// (Iranian) data in Go.
//
// It is organized as a set of focused sub-packages:
//
//   - bank         validate card numbers, resolve banks, and check Sheba (IBAN) codes
//   - bill         parse and validate Iranian utility bills
//   - digit        format numbers and convert digits to Persian words
//   - nationalid   validate national numbers (code-e Melli) and resolve place
//   - phonenumbers validate and inspect Iranian mobile numbers
//
// Import only the sub-packages you need, for example:
//
//	import "github.com/amiranmanesh/go-persian-tools/nationalid"
package persiantools
