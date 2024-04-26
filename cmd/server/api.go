package main

import (
	"fmt"
	"library-api/internal/book"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// Book context
	bookRepo := book.CreateBookRepositoryInstance()
	bookEventPropagator := book.CreateBookEventPropagatorInstance()
	bookEventConsumer := book.CreateBooEventConsumerInstance()

	bookSvc := book.NewService(bookRepo, bookEventPropagator)

	bookEventConsumer.BindBookCreated(func(bc *book.BookCreated) {
		bok, _ := bookSvc.GetBookByID(bc.Id)
		fmt.Println("A book has been created", bok)
	})

	go bookEventConsumer.StartConsuming()

	router := gin.Default()
	book.RegisterHandlers(router, bookSvc)

	// Ping context
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	go startServer(router)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Aplicación detenida")
}

func startServer(gin *gin.Engine) {
	// Iniciar servidor en rutina
	if err := gin.Run(":8080"); err != nil {
		panic(fmt.Sprintf("Error running HTTP Server: %s", err.Error()))
	}
}
