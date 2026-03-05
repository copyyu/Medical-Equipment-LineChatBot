package usecase

import (
	"strings"
	"testing"
)

// ─── SanitizeInput ────────────────────────────────────────

func TestSanitizeInput_TrimWhitespace(t *testing.T) {
	got := SanitizeInput("  hello  ")
	if got != "hello" {
		t.Errorf("expected 'hello', got '%s'", got)
	}
}

func TestSanitizeInput_LimitLength(t *testing.T) {
	long := strings.Repeat("A", 200)
	got := SanitizeInput(long)
	if len(got) != MaxInputLength {
		t.Errorf("expected max length %d, got %d", MaxInputLength, len(got))
	}
}

func TestSanitizeInput_RemoveControlChars(t *testing.T) {
	got := SanitizeInput("hello\x00world")
	if got != "helloworld" {
		t.Errorf("expected 'helloworld', got '%s'", got)
	}
}

func TestSanitizeInput_PreserveNewline(t *testing.T) {
	got := SanitizeInput("line1\nline2")
	if got != "line1\nline2" {
		t.Errorf("newlines should be preserved, got '%s'", got)
	}
}

func TestSanitizeInput_Empty(t *testing.T) {
	got := SanitizeInput("")
	if got != "" {
		t.Errorf("expected empty string, got '%s'", got)
	}
}

// ─── IsValidEquipmentCode ─────────────────────────────────

func TestIsValidEquipmentCode_Valid(t *testing.T) {
	cases := []string{"SSH12345", "MED001", "AB12345", "ABC_12345", "XY-999"}
	for _, c := range cases {
		if !IsValidEquipmentCode(c) {
			t.Errorf("expected valid: %s", c)
		}
	}
}

func TestIsValidEquipmentCode_TooShort(t *testing.T) {
	if IsValidEquipmentCode("A1") {
		t.Error("2-char code should be invalid (below MinSerialLength)")
	}
}

func TestIsValidEquipmentCode_InvalidChars(t *testing.T) {
	if IsValidEquipmentCode("SSH 123") {
		t.Error("code with spaces should be invalid")
	}
}

func TestIsValidEquipmentCode_NoDigits(t *testing.T) {
	if IsValidEquipmentCode("ABCDEF") {
		t.Error("code without digits should be invalid")
	}
}

// ─── IsValidTicketNumber ──────────────────────────────────

func TestIsValidTicketNumber_Valid(t *testing.T) {
	if !IsValidTicketNumber("TK-2024001") {
		t.Error("TK-2024001 should be valid")
	}
}

func TestIsValidTicketNumber_NoPrefix(t *testing.T) {
	if IsValidTicketNumber("2024001") {
		t.Error("ticket without TK- prefix should be invalid")
	}
}

func TestIsValidTicketNumber_TooFewDigits(t *testing.T) {
	if IsValidTicketNumber("TK-123") {
		t.Error("TK-123 (only 3 digits) should be invalid")
	}
}

// ─── IsAlphanumericWithSeparators ─────────────────────────

func TestIsAlphanumericWithSeparators_Valid(t *testing.T) {
	cases := []string{"ABC123", "hello-world", "foo_bar", "MED-001"}
	for _, c := range cases {
		if !IsAlphanumericWithSeparators(c) {
			t.Errorf("expected valid: %s", c)
		}
	}
}

func TestIsAlphanumericWithSeparators_Empty(t *testing.T) {
	if IsAlphanumericWithSeparators("") {
		t.Error("empty string should return false")
	}
}

func TestIsAlphanumericWithSeparators_SpecialChars(t *testing.T) {
	cases := []string{"hello world", "foo@bar", "abc!123", "สวัสดี"}
	for _, c := range cases {
		if IsAlphanumericWithSeparators(c) {
			t.Errorf("expected invalid: %s", c)
		}
	}
}

// ─── ValidateAndSanitizeSerial ────────────────────────────

func TestValidateAndSanitizeSerial_Valid(t *testing.T) {
	got, ok := ValidateAndSanitizeSerial("  SSH12345  ")
	if !ok || got != "SSH12345" {
		t.Errorf("expected (SSH12345, true), got (%s, %v)", got, ok)
	}
}

func TestValidateAndSanitizeSerial_TooShort(t *testing.T) {
	_, ok := ValidateAndSanitizeSerial("AB")
	if ok {
		t.Error("2-char input should be invalid")
	}
}

func TestValidateAndSanitizeSerial_InvalidChars(t *testing.T) {
	_, ok := ValidateAndSanitizeSerial("สวัสดี")
	if ok {
		t.Error("Thai characters should be invalid for serial")
	}
}
