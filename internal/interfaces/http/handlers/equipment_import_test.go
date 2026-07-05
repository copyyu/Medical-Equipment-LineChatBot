package handlers

import "testing"

func TestIsValidExcelFile(t *testing.T) {
	cases := []struct {
		filename string
		want     bool
	}{
		{"report.xlsx", true},
		{"report.xls", true},
		{"Report.XLSX", true}, // uppercase (Windows) must be accepted
		{"REPORT.XLS", true},
		{"data.Xlsx", true},
		{"report.csv", false},
		{"report", false},
		{"", false},
		{"xlsx", false}, // no dot/name
		{"archive.xlsx.zip", false},
	}
	for _, tc := range cases {
		t.Run(tc.filename, func(t *testing.T) {
			if got := isValidExcelFile(tc.filename); got != tc.want {
				t.Errorf("isValidExcelFile(%q) = %v, want %v", tc.filename, got, tc.want)
			}
		})
	}
}
