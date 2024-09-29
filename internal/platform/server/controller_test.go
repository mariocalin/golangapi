package server_test

import (
	"bytes"
	"library-api/internal"
	"library-api/internal/platform/server"
	"library-api/internal/platform/server/mocks"
	"library-api/kit/date"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ControllerTestSuite struct {
	suite.Suite
	adapter     *mocks.BookAdapter
	datehandler date.Handler
	router      *gin.Engine
	sut         server.BookController
}

func TestControllerTestSuite(t *testing.T) {
	if integrationEnabled := os.Getenv("RUN_INTEGRATION_TESTS"); integrationEnabled != "1" {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(ControllerTestSuite))
}

func (s *ControllerTestSuite) SetupSuite() {
	s.adapter = mocks.NewBookAdapter(s.T())
	s.datehandler = date.NewLocalHandler()
	s.sut = server.NewBookController(s.adapter, s.datehandler)

	s.router = setupRouter()
	s.router.POST("/book", s.sut.CreateBook)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func (s *ControllerTestSuite) TestCreateBookHandlerOk() {
	date, err := time.Parse("2006-01-02", "2021-01-01")
	s.Require().NoError(err)

	book, err := internal.BookFactory.CreateBook("Book Name", date, []string{"category1", "category2"}, "Book Description")
	s.Require().NoError(err)

	s.adapter.EXPECT().CreateBook(mock.Anything, mock.Anything).Return(book, nil)

	body := `{"name":"Book Name","publish_date":"2021-01-01","categories":["category1","category2"],"description":"Book Description"}`

	req := httptest.NewRequest("POST", "/book", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)
}

func (s *ControllerTestSuite) TestCreateBookHandlerInvalidData() {
	body := `{"name":"","publish_date":"not-a-date","categories":[],"description":""}`
	req := httptest.NewRequest("POST", "/book", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
}
