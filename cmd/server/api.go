package main

import (
	"fmt"
	"library-api/internal/book"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

//	@title			Library Api
//	@version		1.0
//	@description	API for creating and retreiving books

//	@contact.name	Mario
//	@contact.url	http://example.org
//	@contact.email	mario.calin@mindcurv.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
func main() {
	log.Println("Starting api")

	// ---- BOOK ----
	log.Println("Creating context...")

	bookRepo := book.CreateBookRepositoryInstance()
	log.Println("bookRepo created")

	bookEventPropagator := book.CreateBookEventPropagatorInstance()
	log.Println("bookEventPropagator created")

	bookEventConsumer := book.CreateBooEventConsumerInstance()
	log.Println("bookEventConsumer created")

	bookSvc := book.NewService(bookRepo, bookEventPropagator)
	log.Println("bookEventConsumer created")

	bookEventConsumer.BindBookCreated(func(bc *book.BookCreated) {
		bok, _ := bookSvc.GetBookByID(bc.Id)
		log.Println("A book has been created", bok)
	})

	go bookEventConsumer.StartConsuming()
	defer bookEventConsumer.Close()
	defer bookEventPropagator.Close()

	router := gin.Default()
	bookController := book.NewBookController(bookSvc)
	registerHandlers(router, bookController)

	// ---- STATUS ----
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go startServer(router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Aplicación stopped")
}

func startServer(gin *gin.Engine) {
	// Iniciar servidor en rutina
	if err := gin.Run(":8080"); err != nil {
		panic(fmt.Sprintf("Error running HTTP Server: %s", err.Error()))
	}
}

func registerHandlers(r *gin.Engine, bookController *book.BookController) {
	r.GET("/book", bookController.GetAllBooks)
	r.GET("/book/:id", bookController.GetBookById)
	r.POST("/book", bookController.CreateBook)
}
