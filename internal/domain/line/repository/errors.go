package repository

import "errors"

// ErrDuplicate is returned by repositories when a write violates a unique
// constraint. It lets application code detect conflicts (e.g. to retry) without
// depending on the underlying ORM/driver error types.
var ErrDuplicate = errors.New("duplicate record")
