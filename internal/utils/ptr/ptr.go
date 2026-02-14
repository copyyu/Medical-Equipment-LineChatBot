package ptr

// StringPtr returns a pointer to the given string.
// Used throughout the application to convert string values to *string for optional fields.
func StringPtr(s string) *string {
	return &s
}
