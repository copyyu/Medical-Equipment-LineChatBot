package service

import (
	"net/http"
	"testing"

	apperrors "medical-webhook/internal/utils/errors"
)

// The admin service returns its own sentinel vars; the HTTP layer maps them with
// errors.Is. If those sentinels aren't the same values the mapper knows about,
// every admin error falls through to a generic 500 (regression guard for the
// duplicate-sentinel bug).
func TestAdminSentinelsMapToHTTPStatus(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"invalid credentials -> 401", ErrInvalidCredentials, http.StatusUnauthorized},
		{"invalid token -> 401", ErrInvalidToken, http.StatusUnauthorized},
		{"account inactive -> 403", ErrAdminInactive, http.StatusForbidden},
		{"username exists -> 409", ErrUsernameExists, http.StatusConflict},
		{"email exists -> 409", ErrEmailExists, http.StatusConflict},
		{"admin not found -> 404", ErrAdminNotFound, http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, _, _ := apperrors.MapErrorToResponse(tc.err)
			if status != tc.want {
				t.Errorf("MapErrorToResponse(%v) status = %d, want %d", tc.err, status, tc.want)
			}
		})
	}
}
