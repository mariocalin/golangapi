//go:build unit || !integration

package server

import (
	"bytes"
	"library-api/common"
	"library-api/internal/book/application"
	"library-api/internal/book/domain"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestCreateBookHandler(t *testing.T) {
	mockSvc := application.NewMockBookService(t)
	realDateHandler := common.NewDateHandler()
	controller := NewBookController(mockSvc, realDateHandler)

	router := setupRouter()
	router.POST("/book", controller.CreateBook)

	t.Run("Success", func(t *testing.T) {

		bookDate, _ := realDateHandler.DateParse("2022-01-01")

		id := uuid.New()

		name, _ := domain.NewName("Example")
		description, _ := domain.NewDescription("Example")
		categories, _ := domain.NewCategories([]string{"fiction"})
		publishDate := domain.NewPublishDate(bookDate)

		book := domain.Book{
			ID:          &id,
			Name:        name,
			Description: description,
			Categories:  categories,
			PublishDate: publishDate,
		}

		mockSvc.On("CreateBook", mock.Anything).Return(&book, nil)

		body := `{"name":"Example","publish_date":"2022-01-01","categories":["fiction"],"description":"A book description"}`
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
