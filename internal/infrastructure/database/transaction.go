package database

import (
	"context"

	"medical-webhook/internal/domain/port"

	"gorm.io/gorm"
)

// txContextKey is the unexported key under which the active *gorm.DB
// transaction is stored in a context.
type txContextKey struct{}

// ContextWithTx returns a child context carrying tx as the active transaction.
func ContextWithTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// DBFromContext returns the transaction stored in ctx if one is present,
// otherwise the provided fallback handle. Repositories call this so they
// transparently participate in a transaction started by the TxManager, while
// still working normally when no transaction is active.
func DBFromContext(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if ctx != nil {
		if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok && tx != nil {
			return tx
		}
	}
	return fallback
}

// txManager is the GORM-backed implementation of port.TxManager.
type txManager struct {
	db *gorm.DB
}

// NewTxManager creates a new transaction manager.
func NewTxManager(db *gorm.DB) port.TxManager {
	return &txManager{db: db}
}

// WithTransaction runs fn inside a single GORM transaction. It delegates to
// gorm's Transaction helper, which commits on success and rolls back on either
// an error or a panic (re-panicking after the rollback), so no connection is
// leaked. The active transaction is injected into the context passed to fn so
// repositories can pick it up via DBFromContext.
func (m *txManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ContextWithTx(ctx, tx))
	})
}
