//go:build unit || !integration

package application

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"library-api/internal/book/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type ServiceTestSuite struct {
	suite.Suite
	repository *MockRepository
	propagator *MockBookEventPropagator
	sut        Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.repository = NewMockRepository(s.T())
	s.propagator = NewMockBookEventPropagator(s.T())

	s.sut = NewService(s.repository, s.propagator)
}

func (s *ServiceTestSuite) Test_Create_Book() {
	// Arrange
	s.repository.On("Create", mock.Anything).Once().Return(nil)
	s.propagator.On("PropagateBookCreated", mock.Anything).Once().Return(nil)

	name, _ := domain.NewName("Test Book")
	publishedDate := domain.NewPublishDate(time.Now())
	categories, _ := domain.NewCategories([]string{"category"})
	description, _ := domain.NewDescription("This is a test book")

	command := CreateCommand{
		Name:        name,
		PublishDate: publishedDate,
		Categories:  categories,
		Description: description,
	}

	// Act
	bookCreated, err := s.sut.CreateBook(&command)
	s.Require().NoError(err)

	// Assert
	bookArgument := s.repository.Calls[0].Arguments.Get(0).(*domain.Book)
	s.Assert().EqualValues(bookCreated, bookArgument)

	bookCreatedEvent := s.propagator.Calls[0].Arguments.Get(0).(*CreatedEvent)
	s.Assert().Equal(bookCreatedEvent.Id, bookArgument.ID)
}

func (s *ServiceTestSuite) Test_Update_Book_Name() {
	s.repository.On("Update", mock.Anything).Once().Return(nil)

	id := uuid.New()

	name, _ := domain.NewName("Test Book")
	publishedDate := domain.NewPublishDate(time.Now())
	categories, _ := domain.NewCategories([]string{"category"})
	description, _ := domain.NewDescription("This is a test book")

	existingBook := domain.Book{
		ID:          &id,
		Name:        name,
		PublishDate: publishedDate,
		Categories:  categories,
		Description: description,
	}

	s.repository.On("FindByID", existingBook.ID).Once().Return(&existingBook, nil)

	newName, _ := domain.NewName("Updated Test Book")

	command := CreateCommand{
		Name: newName,
	}

	// Act
	err := s.sut.UpdateBook(existingBook.ID, &command)

	// Assert
	s.Require().NoError(err)

	bookArgument := s.repository.Calls[1].Arguments.Get(0).(*domain.Book)
	s.Equal(command.Name.Value(), bookArgument.Name.Value())
	s.Equal(existingBook.Categories.Value(), bookArgument.Categories.Value())
	s.Equal(existingBook.Description.Value(), bookArgument.Description.Value())
	s.Equal(existingBook.PublishDate.Value(), bookArgument.PublishDate.Value())
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
