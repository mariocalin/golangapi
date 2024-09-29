package internal

import (
	"context"
	"errors"
)

type TransactionCtxKeyType string

const TransactionKey TransactionCtxKeyType = "tx"

var ErrTransactionNotFound = errors.New("transaction not found")

// TransactionManager is an interface that defines transaction management methods.
type TransactionManager interface {
	// RunInTransaction creates a transaction and executes the given function within it.
	// If the function returns an error, the transaction is rolled back.
	// If the function returns nil, the transaction is committed.
	// The function is passed a context that is a child of the context with the transaction.
	RunInTransaction(ctx context.Context, function func(ctx context.Context) error) error
}
