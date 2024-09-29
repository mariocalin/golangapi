package create

import (
	"context"
	"fmt"
	"library-api/internal"
	"time"
)

type UseCaseRequest struct {
	Name        string
	PublishDate time.Time
	Categories  []string
	Description string
}

type UseCase struct {
	txManager  internal.TransactionManager
	repository internal.BookRepository
	propagator internal.BookEventPropagator
}

func NewUseCase(repository internal.BookRepository, txManager internal.TransactionManager, propagator internal.BookEventPropagator) UseCase {
	return UseCase{
		txManager:  txManager,
		repository: repository,
		propagator: propagator,
	}
}

func (uc UseCase) Execute(ctx context.Context, request UseCaseRequest) (internal.Book, error) {
	book, err := internal.BookFactory.CreateBook(request.Name, request.PublishDate, request.Categories, request.Description)
	if err != nil {
		return internal.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	err = uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := uc.repository.Create(ctx, book); err != nil {
			return err
		}

		bookCreatedEvent := internal.BookCreated{
			Id:   book.ID,
			Name: book.Name,
		}

		if err := uc.propagator.PropagateBookCreated(ctx, bookCreatedEvent); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return internal.Book{}, fmt.Errorf("failed to create book: %w", err)
	}

	return book, nil
}
