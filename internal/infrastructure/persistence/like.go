package persistence

import "strings"

// escapeLike escapes LIKE/ILIKE wildcards in user-supplied search text so the
// characters are matched literally instead of acting as patterns. PostgreSQL
// uses backslash as the default LIKE escape character, so backslash itself must
// be escaped first.
func escapeLike(s string) string {
	return strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	).Replace(s)
}
