package usecase

import "testing"

func TestClampPagination(t *testing.T) {
	cases := []struct {
		name              string
		page, limit       int
		wantPage, wantLim int
	}{
		{"valid values pass through", 3, 50, 3, 50},
		{"zero page -> 1", 0, 20, 1, 20},
		{"negative page -> 1", -5, 20, 1, 20},
		{"zero limit -> default", 1, 0, 1, defaultPageLimit},
		{"negative limit -> default", 1, -1, 1, defaultPageLimit},
		{"huge limit -> max", 1, 100000000, 1, maxPageLimit},
		{"limit at max stays", 1, maxPageLimit, 1, maxPageLimit},
		{"limit just over max clamps", 1, maxPageLimit + 1, 1, maxPageLimit},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotPage, gotLim := clampPagination(tc.page, tc.limit)
			if gotPage != tc.wantPage || gotLim != tc.wantLim {
				t.Errorf("clampPagination(%d,%d) = (%d,%d), want (%d,%d)",
					tc.page, tc.limit, gotPage, gotLim, tc.wantPage, tc.wantLim)
			}
		})
	}
}
