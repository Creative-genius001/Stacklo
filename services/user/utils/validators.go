package utils

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidEmail(email string) bool {
	// Define the email regex
	emailRegex := regexp.MustCompile(`(?i)^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

	// Check the URL field against the regex
	return emailRegex.MatchString(email)
}

func IsValidPhoneNumber(numberToParse, countryCode string) (isValid bool, formattedNum string, err error) {
	if len(countryCode) <= 1 {
		countryCode = "NG"
	}

	// Parse the phone number
	metadata, err := phonenumbers.Parse(numberToParse, countryCode)
	if err != nil {
		return false, "", err
	}

	// Format the phone number using the national format and remove spaces
	formattedNum = strings.ReplaceAll(phonenumbers.Format(metadata, phonenumbers.E164), " ", "")

	// Check if the phone number is valid
	isValid = phonenumbers.IsValidNumber(metadata)

	return isValid, formattedNum, nil
}
