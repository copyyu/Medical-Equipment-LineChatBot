package usecase

// Validation constants
const (
	MinSerialLength = 3
	MaxInputLength  = 100
)

// Pagination bounds shared by list endpoints. Clamping protects against a
// missing/negative/huge limit (which would otherwise offset badly, load the
// whole table into memory, or divide by zero when computing total pages).
const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

// clampPagination normalizes a requested page/limit to safe bounds.
func clampPagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultPageLimit
	} else if limit > maxPageLimit {
		limit = maxPageLimit
	}
	return page, limit
}
