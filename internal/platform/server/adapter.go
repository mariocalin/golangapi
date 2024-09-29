package server

import (
	"context"
	"library-api/internal"
	"library-api/internal/book/create"
	"library-api/internal/book/findAll"
	"library-api/internal/book/findById"
	"library-api/internal/book/update"
)

type bookAdapter struct {
	createUseCase   create.UseCase
	findAllUseCase  findAll.UseCase
	findByIdUseCase findById.UseCase
	updateUseCase   update.UseCase
}

func NewBookAdapter(
	createUseCase create.UseCase,
	findAllUseCase findAll.UseCase,
	findByIdUseCase findById.UseCase,
	updateUseCase update.UseCase,
) BookAdapter {
	return bookAdapter{
		createUseCase:   createUseCase,
		findAllUseCase:  findAllUseCase,
		findByIdUseCase: findByIdUseCase,
		updateUseCase:   updateUseCase,
	}
}

func (a bookAdapter) GetBooks(ctx context.Context) ([]internal.Book, error) {
	return a.findAllUseCase.Execute(ctx)
}

func (a bookAdapter) GetBookByID(ctx context.Context, id string) (internal.Book, error) {
	return a.findByIdUseCase.Execute(ctx, findById.Request{Id: id})
}

func (a bookAdapter) CreateBook(ctx context.Context, request CreateBookRequest) (internal.Book, error) {
	publishDate, err := parseDateYearMonthDay(request.PublishDate)
	if err != nil {
		return internal.Book{}, err
	}

	usecaseRequest := create.UseCaseRequest{
		Name:        request.Name,
		PublishDate: publishDate,
		Categories:  request.Categories,
		Description: request.Description,
	}

	return a.createUseCase.Execute(ctx, usecaseRequest)
}
