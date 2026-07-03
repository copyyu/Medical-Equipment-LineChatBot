// Package port defines application-facing interfaces (ports) that the
// application layer depends on, with concrete implementations living in the
// infrastructure layer. Keeping these interfaces here preserves the clean
// architecture dependency rule: application depends on abstractions, not on
// GORM/Redis/etc.
package port

import "context"

// TxManager runs a set of operations inside a single database transaction.
//
// The callback receives a context that carries the active transaction. Any
// repository call made with that context participates in the same transaction;
// if the callback returns an error (or panics) the transaction is rolled back,
// otherwise it is committed. The interface is deliberately free of any ORM type
// so the application layer stays decoupled from the persistence technology.
type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
