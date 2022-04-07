package validators

import (
	"fmt"
	"regexp"
)

// var phoneMatch = regexp.MustCompile(`(\d{2}\d{2}9\d{8})`).
var phoneMatch = regexp.MustCompile(`[0-9]+`)

func NormalizePhoneNumber(phone string) (string, error) {
	pn := string(phoneMatch.Find([]byte(phone)))
	if pn == "" {
		return "", fmt.Errorf("'%s' is not a valid phone number", phone)
	}

	return pn, nil
}
