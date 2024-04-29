package book

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BookController struct {
	svc BookService
}

func NewBookController(svc BookService) *BookController {
	return &BookController{
		svc: svc,
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
func (controller *BookController) GetAllBooks(c *gin.Context) {
	log.Println("Calling getAllBooksHandler")
	books, err := controller.svc.GetBooks()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, createBookResources(books))
}

func (controller *BookController) GetBookById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad id"})
		return
	}

	book, err := controller.svc.GetBookByID(id)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if book == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, createBookResource(book))
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
func (controller *BookController) CreateBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse publish date
	name, err := NewName(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}

	description, err := NewDescription(req.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid description"})
		return
	}

	categories, err := NewCategories(req.Categories)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid categories"})
		return
	}

	publishDate, err := time.Parse(time.DateOnly, req.PublishDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publish date format"})
		return
	}

	// Create command
	command := CreateBookCommand{
		Name:        *name,
		Description: *description,
		Categories:  *categories,
		PublishDate: *NewPublishDate(publishDate)}

	// Add the book using the service
	book, err := controller.svc.CreateBook(&command)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}

	// Return the created book
	c.JSON(http.StatusCreated, createBookResource(book))
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

func createBookResource(book *Book) BookResource {
	return BookResource{
		Id:          book.ID.String(),
		Name:        book.Name.Value(),
		PublishDate: book.PublishDate.Value().Format(time.DateOnly),
		Categories:  book.Categories.Value(),
		Description: book.Description.Value(),
	}
}

func createBookResources(books []Book) []BookResource {
	var resources []BookResource

	for _, book := range books {
		resources = append(resources, createBookResource(&book))
	}

	return resources
}
