package book

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateBookRequest struct {
	Name        string   `json:"name" binding:"required"`
	PublishDate string   `json:"publish_date" binding:"required"`
	Categories  []string `json:"categories" binding:"required"`
	Description string   `json:"description" binding:"required"`
}

func getAllBooksHandler(svc BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := svc.GetBooks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, books)
	}
}

func getBookByIdHandler(svc BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.Parse(c.Param("id"))
		book, err := svc.GetBookByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if book == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusOK, book)
	}
}

func createBookHandler(svc BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON request to CreateBookRequest struct
		var req CreateBookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse publish date
		publishDate, err := time.Parse(time.RFC3339, req.PublishDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid publish date format"})
			return
		}

		// Create Book object
		newBook := Book{
			ID:          uuid.New(),
			Name:        req.Name,
			PublishDate: publishDate,
			Categories:  req.Categories,
			Description: req.Description,
		}

		// Add the book using the service
		err = svc.CreateBook(&newBook)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
			return
		}

		// Return the created book
		c.JSON(http.StatusCreated, newBook)
	}
}

func RegisterHandlers(r *gin.Engine, svc BookService) {
	r.GET("/book", getAllBooksHandler(svc))
	r.GET("/book/:id", getBookByIdHandler(svc))
	r.POST("/book", createBookHandler(svc))
}
