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

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestCreateBookHandler(t *testing.T) {
	mockSvc := NewMockBookService(t)
	controller := NewBookController(mockSvc)

	router := setupRouter()
	router.POST("/book", controller.CreateBook)

	t.Run("Success", func(t *testing.T) {

		bookDate, _ := time.Parse(time.DateOnly, "2022-01-01")

		id := uuid.New()

		book := Book{
			ID:          &id,
			Name:        &Name{"Example"},
			Description: &Description{"Example"},
			Categories:  &Categories{[]string{"fiction"}},
			PublishDate: &PublishDate{bookDate},
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
