package usecase

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── SanitizeInput ────────────────────────────────────────

func TestSanitizeInput_TrimWhitespace(t *testing.T) {
	got := SanitizeInput("  hello  ")
	assert.Equal(t, "hello", got)
}

func TestSanitizeInput_LimitLength(t *testing.T) {
	long := strings.Repeat("A", 200)
	got := SanitizeInput(long)
	assert.Equal(t, MaxInputLength, len(got))
}

func TestSanitizeInput_RemoveControlChars(t *testing.T) {
	got := SanitizeInput("hello\x00world")
	assert.Equal(t, "helloworld", got)
}

func TestSanitizeInput_PreserveNewline(t *testing.T) {
	got := SanitizeInput("line1\nline2")
	assert.Equal(t, "line1\nline2", got, "newlines should be preserved")
}

func TestSanitizeInput_Empty(t *testing.T) {
	got := SanitizeInput("")
	assert.Empty(t, got)
}

// ─── IsValidEquipmentCode ─────────────────────────────────

func TestIsValidEquipmentCode_Valid(t *testing.T) {
	cases := []string{"SSH12345", "MED001", "AB12345", "ABC_12345", "XY-999"}
	for _, c := range cases {
		assert.True(t, IsValidEquipmentCode(c), "expected valid: %s", c)
	}
}

func TestIsValidEquipmentCode_TooShort(t *testing.T) {
	assert.False(t, IsValidEquipmentCode("A1"), "2-char code should be invalid")
}

func TestIsValidEquipmentCode_InvalidChars(t *testing.T) {
	assert.False(t, IsValidEquipmentCode("SSH 123"), "code with spaces should be invalid")
}

func TestIsValidEquipmentCode_NoDigits(t *testing.T) {
	assert.False(t, IsValidEquipmentCode("ABCDEF"), "code without digits should be invalid")
}

// ─── IsValidTicketNumber ──────────────────────────────────

func TestIsValidTicketNumber_Valid(t *testing.T) {
	assert.True(t, IsValidTicketNumber("TK-2024001"))
}

func TestIsValidTicketNumber_NoPrefix(t *testing.T) {
	assert.False(t, IsValidTicketNumber("2024001"), "ticket without TK- prefix should be invalid")
}

func TestIsValidTicketNumber_TooFewDigits(t *testing.T) {
	assert.False(t, IsValidTicketNumber("TK-123"), "TK-123 (only 3 digits) should be invalid")
}

// ─── IsAlphanumericWithSeparators ─────────────────────────

func TestIsAlphanumericWithSeparators_Valid(t *testing.T) {
	cases := []string{"ABC123", "hello-world", "foo_bar", "MED-001"}
	for _, c := range cases {
		assert.True(t, IsAlphanumericWithSeparators(c), "expected valid: %s", c)
	}
}

func TestIsAlphanumericWithSeparators_Empty(t *testing.T) {
	assert.False(t, IsAlphanumericWithSeparators(""), "empty string should return false")
}

func TestIsAlphanumericWithSeparators_SpecialChars(t *testing.T) {
	cases := []string{"hello world", "foo@bar", "abc!123", "สวัสดี"}
	for _, c := range cases {
		assert.False(t, IsAlphanumericWithSeparators(c), "expected invalid: %s", c)
	}
}

// ─── ValidateAndSanitizeSerial ────────────────────────────

func TestValidateAndSanitizeSerial_Valid(t *testing.T) {
	got, ok := ValidateAndSanitizeSerial("  SSH12345  ")
	assert.True(t, ok)
	assert.Equal(t, "SSH12345", got)
}

func TestValidateAndSanitizeSerial_TooShort(t *testing.T) {
	_, ok := ValidateAndSanitizeSerial("AB")
	assert.False(t, ok, "2-char input should be invalid")
}

func TestValidateAndSanitizeSerial_InvalidChars(t *testing.T) {
	_, ok := ValidateAndSanitizeSerial("สวัสดี")
	assert.False(t, ok, "Thai characters should be invalid for serial")
}
