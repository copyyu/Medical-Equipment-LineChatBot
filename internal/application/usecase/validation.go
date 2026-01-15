package usecase

import (
	"regexp"
	"strings"
	"unicode"
)

// Input validation patterns
var (
	// equipmentCodePattern matches equipment codes like SSH12345, MED-001, etc.
	equipmentCodePattern = regexp.MustCompile(`^[A-Za-z]{2,5}[-_]?\d{3,10}$`)

	// ticketPattern matches ticket numbers like TK-2024001
	ticketPattern = regexp.MustCompile(`^TK[-_]?\d{6,10}$`)
)

// SanitizeInput cleans and validates user input
// - Trims whitespace
// - Limits length to MaxInputLength
// - Removes potentially dangerous characters
func SanitizeInput(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Limit length
	if len(input) > MaxInputLength {
		input = input[:MaxInputLength]
	}

	// Remove control characters (except newline and tab for display purposes)
	var cleaned strings.Builder
	for _, r := range input {
		if unicode.IsPrint(r) || r == '\n' || r == '\t' {
			cleaned.WriteRune(r)
		}
	}

	return cleaned.String()
}

// IsValidEquipmentCode checks if the input matches equipment code format
// Examples: SSH12345, MED-001, ABC_12345
func IsValidEquipmentCode(input string) bool {
	if len(input) < MinSerialLength || len(input) > MaxInputLength {
		return false
	}
	return equipmentCodePattern.MatchString(input)
}

// IsValidTicketNumber checks if the input matches ticket number format
// Example: TK-2024001
func IsValidTicketNumber(input string) bool {
	return ticketPattern.MatchString(input)
}

// IsAlphanumericWithSeparators checks if string contains only letters, numbers, and common separators
func IsAlphanumericWithSeparators(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return false
		}
	}
	return true
}

// ValidateAndSanitizeSerial validates and sanitizes equipment serial/code input
// Returns sanitized string and whether it's valid
func ValidateAndSanitizeSerial(input string) (string, bool) {
	sanitized := SanitizeInput(input)

	if len(sanitized) < MinSerialLength {
		return "", false
	}

	// Check if it's alphanumeric with allowed separators
	if !IsAlphanumericWithSeparators(sanitized) {
		return "", false
	}

	return sanitized, true
}
