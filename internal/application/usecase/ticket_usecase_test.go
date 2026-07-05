package usecase

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestIsDuplicateKeyErr(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"gorm sentinel", gorm.ErrDuplicatedKey, true},
		{"wrapped gorm sentinel", errors.Join(errors.New("create ticket"), gorm.ErrDuplicatedKey), true},
		{"pg duplicate key message", errors.New(`ERROR: duplicate key value violates unique constraint "idx_tickets_ticket_no"`), true},
		{"unique constraint message", errors.New("UNIQUE constraint failed: tickets.ticket_no"), true},
		{"unrelated error", errors.New("connection refused"), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isDuplicateKeyErr(tc.err); got != tc.want {
				t.Errorf("isDuplicateKeyErr(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
