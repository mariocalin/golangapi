package book

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_BookService(t *testing.T) {
	t.Run("Create a book", func(t *testing.T) {
		// Arrange
		repo := NewMockBookRepository(t)
		event := NewMockBookEventPropagator(t)
		service := NewService(repo, event)

		repo.On("Create", mock.Anything).Return(nil)
		event.On("PropagateBookCreated", mock.Anything).Return(nil)

		command := BookCommand{
			Name:        &Name{"Test Book"},
			PublishDate: &PublishDate{time.Now()},
			Categories:  &Categories{[]string{"category"}},
			Description: &Description{"This is a test book"},
		}

		// Act
		_, err := service.CreateBook(&command)

		// Assert
		assert.NoError(t, err)

		repo.AssertCalled(t, "Create", mock.Anything)
		event.AssertCalled(t, "PropagateBookCreated", mock.Anything)

		bookArgument := repo.Calls[0].Arguments.Get(0).(*Book)
		assert.Equal(t, command.Name.Value(), bookArgument.Name.Value())
		assert.Equal(t, command.Categories.Value(), bookArgument.Categories.Value())
		assert.Equal(t, command.Description.Value(), bookArgument.Description.Value())
		assert.Equal(t, command.PublishDate.Value(), bookArgument.PublishDate.Value())
		assert.NotNil(t, bookArgument.ID)

		bookCreatedEvent := event.Calls[0].Arguments.Get(0).(*BookCreated)

		assert.Equal(t, bookCreatedEvent.Id, bookArgument.ID)

		repo.AssertExpectations(t)
		event.AssertExpectations(t)
	})

	t.Run("Update the name of the book", func(t *testing.T) {
		repo := NewMockBookRepository(t)
		event := NewMockBookEventPropagator(t)
		service := NewService(repo, event)

		repo.On("Update", mock.Anything).Return(nil)

		id := uuid.New()

		existingBook := Book{
			ID:          &id,
			Name:        &Name{"Test Book"},
			PublishDate: &PublishDate{time.Now()},
			Categories:  &Categories{[]string{"category"}},
			Description: &Description{"This is a test book"},
		}

		repo.On("FindByID", existingBook.ID).Return(&existingBook, nil)

		command := BookCommand{
			Name: &Name{"Updated Test Book"},
		}

		// Act
		err := service.UpdateBook(existingBook.ID, &command)

		// Assert
		assert.NoError(t, err)
		repo.AssertCalled(t, "FindByID", mock.Anything)
		repo.AssertCalled(t, "Update", mock.Anything)

		bookArgument := repo.Calls[1].Arguments.Get(0).(*Book)
		fmt.Println(bookArgument)
		assert.Equal(t, command.Name.Value(), bookArgument.Name.Value())
		assert.Equal(t, existingBook.Categories.Value(), bookArgument.Categories.Value())
		assert.Equal(t, existingBook.Description.Value(), bookArgument.Description.Value())
		assert.Equal(t, existingBook.PublishDate.Value(), bookArgument.PublishDate.Value())

		assert.NotNil(t, bookArgument.ID)

		repo.AssertExpectations(t)
		event.AssertExpectations(t)
	})

}
