package book

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Creamos una implementación de Mock para el repositorio de libros
type mockBookRepository struct {
	mock.Mock
}

// Implementamos el método Create del repositorio de libros mock
func (m *mockBookRepository) Create(book *Book) error {
	args := m.Called(book)
	return args.Error(0)
}

// Implementamos el método FindAll del repositorio de libros mock
func (m *mockBookRepository) FindAll() ([]Book, error) {
	args := m.Called()
	return args.Get(0).([]Book), args.Error(1)
}

// Implementamos el método FindByID del repositorio de libros mock
func (m *mockBookRepository) FindByID(id BookId) (*Book, error) {
	args := m.Called(id)
	return args.Get(0).(*Book), args.Error(1)
}

func TestCreateBook(t *testing.T) {
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
}
