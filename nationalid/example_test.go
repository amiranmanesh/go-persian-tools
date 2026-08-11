package nationalid_test

import (
	"fmt"

	"github.com/amiranmanesh/go-persian-tools/nationalid"
)

func ExampleValidate() {
	fmt.Println(nationalid.Validate("0067749828"))
	// Output: true
}

func ExampleGetPlaceByIranNationalID() {
	place := nationalid.GetPlaceByIranNationalID("0499370899")
	fmt.Printf("%s, %s\n", place.City, place.Province)
	// Output: شهرری, تهران
}
