package sqlite3

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"library-api/internal"
)

type TransactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// RunInTransaction implements internal.TransactionManager.
func (t TransactionManager) RunInTransaction(ctx context.Context, function func(ctx context.Context) error) error {
	databaseTransaction, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	ctx = context.WithValue(ctx, internal.TransactionKey, databaseTransaction)

	if err := function(ctx); err != nil {
		if rollbackErr := databaseTransaction.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}

		return fmt.Errorf("error running transaction: %w", err)
	}

	if err := databaseTransaction.Commit(); err != nil {
		if rollbackErr := databaseTransaction.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}

		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
