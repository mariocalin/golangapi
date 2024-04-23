package book

import (
	"bytes"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) GetBooks() ([]Book, error) {
	args := m.Called()
	return args.Get(0).([]Book), args.Error(1)
}

func (m *MockBookService) GetBookByID(id BookId) (*Book, error) {
	args := m.Called(id)
	return args.Get(0).(*Book), args.Error(1)
}

func (m *MockBookService) CreateBook(cmd *CreateBookCommand) (*Book, error) {
	args := m.Called(cmd)
	return args.Get(0).(*Book), args.Error(1)
}

func (m *MockBookService) UpdateBook(id BookId, cmd *UpdateBookCommand) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestCreateBookHandler(t *testing.T) {
	mockSvc := new(MockBookService)
	router := setupRouter()
	router.POST("/book", createBookHandler(mockSvc))

	t.Run("Success", func(t *testing.T) {

		bookDate, _ := time.Parse(time.RFC3339, "2022-01-01T15:04:05Z")

		book := Book{
			ID:          uuid.New(),
			Name:        Name{"Example"},
			Description: Description{"Example"},
			Categories:  Categories{[]string{"fiction"}},
			PublishDate: PublishDate{bookDate},
		}

		mockSvc.On("CreateBook", mock.Anything).Return(&book, nil)

		body := `{"name":"Example","publish_date":"2022-01-01T15:04:05Z","categories":["fiction"],"description":"A book description"}`
		req := httptest.NewRequest("POST", "/book", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid Data", func(t *testing.T) {
		body := `{"name":"","publish_date":"not-a-date","categories":[],"description":""}`
		req := httptest.NewRequest("POST", "/book", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
	})
}
