package handlers

import "regexp"

var (
	passportNumberPattern = regexp.MustCompile(`^\d{4}\s\d{6}$`)
	passportSeriePattern  = regexp.MustCompile(`^\d{4}$`)
	passportNumPattern    = regexp.MustCompile(`^\d{6}$`)
)

// ValidatePassportNumber validates the format of the passport number
func ValidatePassportNumber(passportNumber string) bool {
	return passportNumberPattern.MatchString(passportNumber)
}

// ValidatePassportSerieAndNumber validates the passport series and number separately
func ValidatePassportSerieAndNumber(passportSerie, passportNumber string) bool {
	return passportSeriePattern.MatchString(passportSerie) && passportNumPattern.MatchString(passportNumber)
}
