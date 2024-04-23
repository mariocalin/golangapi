package book

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBookRepository struct {
	mock.Mock
}

func (m *mockBookRepository) Create(book *Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *mockBookRepository) FindAll() ([]Book, error) {
	args := m.Called()
	return args.Get(0).([]Book), args.Error(1)
}

func (m *mockBookRepository) FindByID(id BookId) (*Book, error) {
	args := m.Called(id)
	return args.Get(0).(*Book), args.Error(1)
}

func (m *mockBookRepository) Update(book *Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func Test_CreateBook(t *testing.T) {
	// Arrange
	repo := new(mockBookRepository)
	repo.On("Create", mock.Anything).Return(nil)

	service := NewService(repo)

	command := CreateBookCommand{
		Name:        Name{"Test Book"},
		PublishDate: PublishDate{time.Now()},
		Categories:  Categories{[]string{"category"}},
		Description: Description{"This is a test book"},
	}

	// Act
	_, err := service.CreateBook(&command)

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "Create", mock.Anything)

	bookArgument := repo.Calls[0].Arguments.Get(0).(*Book)
	assert.Equal(t, command.Name.Value(), bookArgument.Name.Value())
	assert.Equal(t, command.Categories.Value(), bookArgument.Categories.Value())
	assert.Equal(t, command.Description.Value(), bookArgument.Description.Value())
	assert.Equal(t, command.PublishDate.Value(), bookArgument.PublishDate.Value())
	assert.NotNil(t, bookArgument.ID)
	repo.AssertExpectations(t)
}

func Test_UpdateBookName(t *testing.T) {
	// Arrange
	repo := new(mockBookRepository)
	repo.On("Update", mock.Anything).Return(nil)

	existingBook := Book{
		ID:          uuid.New(),
		Name:        Name{"Test Book"},
		PublishDate: PublishDate{time.Now()},
		Categories:  Categories{[]string{"category"}},
		Description: Description{"This is a test book"},
	}

	repo.On("FindByID", existingBook.ID).Return(&existingBook, nil)

	service := NewService(repo)

	command := UpdateBookCommand{
		Name: &Name{"Updated Test Book"},
	}

	// Act
	err := service.UpdateBook(existingBook.ID, &command)

	// Assert
	assert.NoError(t, err)
	repo.AssertCalled(t, "FindByID", mock.Anything)
	repo.AssertCalled(t, "Update", mock.Anything)

	calls := repo.Calls

	bookArgument := calls[1].Arguments.Get(0).(*Book)
	fmt.Println(bookArgument)
	assert.Equal(t, command.Name.Value(), bookArgument.Name.Value())
	assert.Equal(t, existingBook.Categories.Value(), bookArgument.Categories.Value())
	assert.Equal(t, existingBook.Description.Value(), bookArgument.Description.Value())
	assert.Equal(t, existingBook.PublishDate.Value(), bookArgument.PublishDate.Value())

	assert.NotNil(t, bookArgument.ID)
	repo.AssertExpectations(t)
}
