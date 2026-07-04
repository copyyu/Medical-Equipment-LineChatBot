package exporturl

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

// parseSignedURL extracts the query values a handler would read from a link.
func parseSignedURL(t *testing.T, signed string) (dept, filter, exp, sig string) {
	t.Helper()
	u, err := url.Parse(signed)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	q := u.Query()
	return q.Get("dept_id"), q.Get("filter"), q.Get("exp"), q.Get("sig")
}

func TestSignedURL_VerifyRoundTrip(t *testing.T) {
	Init("test-secret")
	dept7 := uint(7)

	cases := []struct {
		name   string
		deptID *uint
		filter string
	}{
		{"all departments, default filter", nil, ""},
		{"all departments, this_year", nil, "this_year"},
		{"single department", &dept7, ""},
		{"single department, next_year", &dept7, "next_year"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			signed := SignedURL("https://example.com", tc.deptID, tc.filter)
			if !strings.HasPrefix(signed, "https://example.com"+Path+"?") {
				t.Fatalf("unexpected url: %s", signed)
			}
			dept, filter, exp, sig := parseSignedURL(t, signed)
			if err := Verify(dept, filter, exp, sig, time.Now()); err != nil {
				t.Fatalf("expected valid link, got %v", err)
			}
		})
	}
}

func TestVerify_RejectsTampering(t *testing.T) {
	Init("test-secret")
	dept7 := uint(7)
	signed := SignedURL("https://example.com", &dept7, "this_year")
	dept, filter, exp, sig := parseSignedURL(t, signed)

	// A valid link must verify.
	if err := Verify(dept, filter, exp, sig, time.Now()); err != nil {
		t.Fatalf("baseline should verify: %v", err)
	}

	// Changing the department must invalidate the signature (privilege / scope tamper).
	if err := Verify("9", filter, exp, sig, time.Now()); err != ErrInvalidSignature {
		t.Errorf("tampered dept_id: want ErrInvalidSignature, got %v", err)
	}
	// Changing the filter must invalidate.
	if err := Verify(dept, "all", exp, sig, time.Now()); err != ErrInvalidSignature {
		t.Errorf("tampered filter: want ErrInvalidSignature, got %v", err)
	}
	// Extending the expiry must invalidate (exp is signed).
	if err := Verify(dept, filter, "9999999999", sig, time.Now()); err != ErrInvalidSignature {
		t.Errorf("tampered exp: want ErrInvalidSignature, got %v", err)
	}
	// Garbage signature must invalidate.
	if err := Verify(dept, filter, exp, "deadbeef", time.Now()); err != ErrInvalidSignature {
		t.Errorf("garbage sig: want ErrInvalidSignature, got %v", err)
	}
}

func TestVerify_Expired(t *testing.T) {
	Init("test-secret")
	// Sign a link that expired one hour ago.
	signed := signedURLAt("https://example.com", nil, "all", time.Now().Add(-time.Hour))
	_, filter, exp, sig := parseSignedURL(t, signed)
	if err := Verify("", filter, exp, sig, time.Now()); err != ErrExpired {
		t.Errorf("want ErrExpired, got %v", err)
	}
}

func TestVerify_NotConfigured(t *testing.T) {
	key = nil // simulate Init never called
	if err := Verify("", "all", "123", "abc", time.Now()); err != ErrNotConfigured {
		t.Errorf("want ErrNotConfigured, got %v", err)
	}
}

func TestVerify_InvalidExp(t *testing.T) {
	Init("test-secret")
	if err := Verify("", "all", "not-a-number", "abc", time.Now()); err != ErrInvalidExp {
		t.Errorf("want ErrInvalidExp, got %v", err)
	}
}

func TestVerify_DefaultFilterMatchesEmpty(t *testing.T) {
	Init("test-secret")
	// A link built with an empty filter is written as filter=all; a handler that
	// reads the default "all" must still verify.
	signed := SignedURL("https://example.com", nil, "")
	_, _, exp, sig := parseSignedURL(t, signed)
	if err := Verify("", "all", exp, sig, time.Now()); err != nil {
		t.Errorf("default filter should verify: %v", err)
	}
}
