package server

import (
	"context"
	"errors"
	"fmt"
	"library-api/internal"
	"library-api/kit/date"
	"library-api/kit/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockery --name=BookAdapter --output=mocks --filename=book_adapter_mock.go --with-expecter=true
type BookAdapter interface {
	GetBooks(ctx context.Context) ([]internal.Book, error)
	GetBookByID(ctx context.Context, id string) (internal.Book, error)
	CreateBook(ctx context.Context, request CreateBookRequest) (internal.Book, error)
}

type BookController struct {
	adapter     BookAdapter
	dateHandler date.Handler
}

func NewBookController(adapter BookAdapter, dateHandler date.Handler) BookController {
	return BookController{
		adapter:     adapter,
		dateHandler: dateHandler,
	}
}

// GetAllBooks godoc
//
//	@Summary		Get all persisted books
//	@Description	get all books that are stored in the system
//	@ID				get-all-books
//	@Tags			books
//	@Produce		json
//	@Success		200	{object}	BookResource
//	@Router			/book [get]
func (c *BookController) GetAllBooks(gc *gin.Context) {
	log.Info("Calling getAllBooksHandler")

	books, err := c.adapter.GetBooks(gc.Request.Context())
	if err != nil {
		fmt.Println(err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	gc.JSON(http.StatusOK, c.createBookResources(books))
}

func (c *BookController) GetBookById(gc *gin.Context) {
	book, err := c.adapter.GetBookByID(gc.Request.Context(), gc.Param("id"))
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			gc.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		fmt.Println(err.Error())
		gc.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	gc.JSON(http.StatusOK, c.createBookResource(book))
}

// CreateBook godoc
//
//	@Summary		Create a book
//	@Description	Creat a book with required parameters
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			request		body	CreateBookRequest	true	"Book"
//
// @Success		200		{object}	BookResource
// @Router			/book [post]
func (c *BookController) CreateBook(gc *gin.Context) {
	var req CreateBookRequest
	if err := gc.ShouldBindJSON(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the book using the service
	book, err := c.adapter.CreateBook(gc.Request.Context(), req)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": "Failed to add book"})
		return
	}

	// Return the created book
	gc.JSON(http.StatusCreated, c.createBookResource(book))
}

type CreateBookRequest struct {
	Name        string   `json:"name" binding:"required"`
	PublishDate string   `json:"publish_date" binding:"required"`
	Categories  []string `json:"categories" binding:"required"`
	Description string   `json:"description" binding:"required"`
}

type BookResource struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	PublishDate string   `json:"publish_date"`
	Categories  []string `json:"categories"`
	Description string   `json:"description"`
}

func (c *BookController) createBookResource(book internal.Book) BookResource {
	return BookResource{
		Id:          book.ID.String(),
		Name:        book.Name.Value(),
		PublishDate: c.dateHandler.DateToString(book.PublishDate.Value()),
		Categories:  book.Categories.Value(),
		Description: book.Description.Value(),
	}
}

func (c *BookController) createBookResources(books []internal.Book) []BookResource {
	resources := make([]BookResource, 0, len(books))

	for _, book := range books {
		resources = append(resources, c.createBookResource(book))
	}

	return resources
}
