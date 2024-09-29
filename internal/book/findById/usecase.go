package findById

import (
	"context"
	"fmt"
	"library-api/internal"

	"github.com/google/uuid"
)

type Request struct {
	Id string
}

type UseCase struct {
	txManager  internal.TransactionManager
	repository internal.BookRepository
}

func NewUseCase(repository internal.BookRepository, txManager internal.TransactionManager) UseCase {
	return UseCase{repository: repository, txManager: txManager}
}

func (uc UseCase) Execute(ctx context.Context, request Request) (internal.Book, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return internal.Book{}, fmt.Errorf("failed to parse id: %w", err)
	}

	var book internal.Book

	err = uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		var inner error

		book, inner = uc.repository.FindByID(ctx, internal.BookId(id))
		if inner != nil {
			return err
		}

		return nil
	})

	return book, nil
}
