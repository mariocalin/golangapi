package application

import (
	"errors"
	"github.com/google/uuid"
	"library-api/internal/book/domain"
)

type service struct {
	repo   BookRepository
	events EventPropagator
}

func NewService(repo BookRepository, events EventPropagator) Service {
	return &service{
		repo:   repo,
		events: events,
	}
}

func (s *service) GetBooks() ([]domain.Book, error) {
	return s.repo.FindAll()
}

func (s *service) GetBookByID(id *domain.Id) (*domain.Book, error) {
	return s.repo.FindByID(id)
}

func (s *service) CreateBook(command *CreateCommand) (*domain.Book, error) {

	if !command.isValidForCreate() {
		return nil, errors.New("cannot create book")
	}

	id := uuid.New()

	newBook := domain.Book{
		ID:          &id,
		Name:        command.Name,
		PublishDate: command.PublishDate,
		Categories:  command.Categories,
		Description: command.Description,
	}

	s.repo.Create(&newBook)
	s.events.PropagateBookCreated(&CreatedEvent{Id: newBook.ID, Name: newBook.Name})

	return &newBook, nil
}

func (s *service) UpdateBook(id *domain.Id, command *CreateCommand) error {
	existingBook, notFound := s.repo.FindByID(id)
	if notFound != nil {
		return notFound
	}

	if command.Name != nil {
		existingBook.Name = command.Name
	}

	if command.PublishDate != nil {
		existingBook.PublishDate = command.PublishDate
	}

	if command.Categories != nil {
		existingBook.Categories = command.Categories
	}

	if command.Description != nil {
		existingBook.Description = command.Description
	}

	s.repo.Update(existingBook)

	return nil
}
