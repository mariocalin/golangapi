package findAll

import (
	"context"
	"fmt"
	"library-api/internal"
)

type UseCase struct {
	repository internal.BookRepository
	txManager  internal.TransactionManager
}

func NewUseCase(repository internal.BookRepository, txManager internal.TransactionManager) UseCase {
	return UseCase{
		repository: repository,
		txManager:  txManager,
	}
}

func (uc UseCase) Execute(ctx context.Context) ([]internal.Book, error) {
	var books []internal.Book
	if err := uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		var innerErr error
		books, innerErr = uc.repository.FindAll(ctx)

		if innerErr != nil {
			return innerErr
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}

	return books, nil
}
