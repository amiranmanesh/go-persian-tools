package nationalid

import (
	"strings"
)

// Place holds the city, province and matching code prefixes resolved from a
// national id.
type Place struct {
	Codes    []string
	City     string
	Province string
}

// IPlaceByNationalId is the former name of [Place].
//
// Deprecated: use [Place].
type IPlaceByNationalId = Place //nolint:revive // deprecated alias kept for compatibility

func getAllCities(code string) []nationalCode {
	var allCities []nationalCode
	for _, s := range getNationalCodes() {
		findCity := strings.Contains(s.code, code)
		if findCity {
			allCities = append(allCities, s)
		}
	}
	return allCities
}

func getProvince(allCities []nationalCode) provinceCode {
	var findProvince provinceCode
	for _, s := range getProvincesCode() {
		if s.code == allCities[0].parentCode {
			findProvince = s
		}
	}
	return findProvince
}

// GetPlaceByIranNationalID resolves the city and province that issued the given
// 10-digit national id. It returns a zero [Place] when the id has the wrong
// length or its prefix matches no known place.
func GetPlaceByIranNationalID(nationalID string) Place {
	if len(nationalID) == 10 {
		code := nationalID[0:3]
		allCities := getAllCities(code)

		if len(allCities) > 0 {
			findProvince := getProvince(allCities)
			codeString := allCities[0].code

			// strings.Split returns a single-element slice when there is no
			// separator, so it handles both the "048-049" and "170" cases.
			return Place{
				City:     allCities[0].city,
				Province: findProvince.city,
				Codes:    strings.Split(codeString, "-"),
			}
		}
		return Place{}
	}

	return Place{}
}

// GetPlaceByIranNationalId is the former name of [GetPlaceByIranNationalID].
//
// Deprecated: use [GetPlaceByIranNationalID].
func GetPlaceByIranNationalId(nationalID string) Place { //nolint:revive // deprecated alias kept for compatibility
	return GetPlaceByIranNationalID(nationalID)
}
