package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}

func ValidateUUID(uuidStr string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(uuidStr))
}

func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

func ValidateLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func ValidateLongitude(lng float64) bool {
	return lng >= -180 && lng <= 180
}
