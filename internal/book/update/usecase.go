package update

import (
	"context"
	"fmt"
	"library-api/internal"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	Id          string
	Name        *string
	PublishDate *time.Time
	Categories  *[]string
	Description *string
}

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

func (uc UseCase) Execute(ctx context.Context, request Request) (internal.Book, error) {
	if uc.isRequestUpdateEmpty(request) {
		return internal.Book{}, fmt.Errorf("request is empty")
	}

	id, err := uuid.Parse(request.Id)
	if err != nil {
		return internal.Book{}, fmt.Errorf("failed to parse id: %w", err)
	}

	book, err := uc.repository.FindByID(ctx, id)
	if err != nil {
		return internal.Book{}, err
	}

	if err := uc.updateNameIfNotNil(&book, request.Name); err != nil {
		return internal.Book{}, fmt.Errorf("failed to update book name: %w", err)
	}

	if err := uc.updatePublishDateIfNotNil(&book, request.PublishDate); err != nil {
		return internal.Book{}, fmt.Errorf("failed to update book publish date: %w", err)
	}

	if err := uc.updateCategoriesIfNotNil(&book, request.Categories); err != nil {
		return internal.Book{}, fmt.Errorf("failed to update book categories: %w", err)
	}

	if err := uc.updateDescriptionIfNotNil(&book, request.Description); err != nil {
		return internal.Book{}, fmt.Errorf("failed to update book description: %w", err)
	}

	if err := uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		return uc.repository.Update(ctx, book)
	}); err != nil {
		return internal.Book{}, fmt.Errorf("failed to update book: %w", err)
	}

	return book, nil
}

func (uc UseCase) isRequestUpdateEmpty(request Request) bool {
	return request.Name == nil && request.PublishDate == nil && request.Categories == nil && request.Description == nil
}

func (uc UseCase) updateNameIfNotNil(book *internal.Book, name *string) error {
	if name == nil {
		return nil
	}

	bookName, err := internal.NewBookName(*name)
	if err != nil {
		return fmt.Errorf("failed to create book name %s: %w", *name, err)
	}

	book.Name = bookName

	return nil
}

func (uc UseCase) updatePublishDateIfNotNil(book *internal.Book, publishDate *time.Time) error {
	if publishDate == nil {
		return nil
	}

	bookPublishDate, err := internal.NewBookPublishDate(*publishDate)
	if err != nil {
		return fmt.Errorf("failed to create book publish date %s: %w", publishDate, err)
	}

	book.PublishDate = bookPublishDate

	return nil
}

func (uc UseCase) updateCategoriesIfNotNil(book *internal.Book, categories *[]string) error {
	if categories == nil {
		return nil
	}

	bookCategories, err := internal.NewBookCategories(*categories)
	if err != nil {
		return fmt.Errorf("failed to create book categories %s: %w", categories, err)
	}

	book.Categories = bookCategories

	return nil
}

func (uc UseCase) updateDescriptionIfNotNil(book *internal.Book, description *string) error {
	if description == nil {
		return nil
	}

	bookDescription, err := internal.NewBookDescription(*description)
	if err != nil {
		return fmt.Errorf("failed to create book description %s: %w", *description, err)
	}

	book.Description = bookDescription

	return nil
}
