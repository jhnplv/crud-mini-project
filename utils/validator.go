package utils

import (
	"strings"
)

func IsEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
