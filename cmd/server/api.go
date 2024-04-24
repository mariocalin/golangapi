package main

import (
	"fmt"
	"library-api/internal/book"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// Local kafka info
var (
	brokers = []string{"localhost:9092"}
	topic   = "test_topic"
)

func main() {
	// Book context
	bookRepo := book.NewSqlite3BookRepository("data.sqlite3")
	bookSvc := book.NewService(bookRepo)

	// Configuración del productor de Kafka
	producer, err := sarama.NewSyncProducer(brokers, nil)
	if err != nil {
		panic(fmt.Sprintf("Error creating Kafka producer: %s", err))
	}
	defer producer.Close()

	// Configuración del consumidor de Kafka
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(fmt.Sprintf("Error creando Kafka consumer: %s", err))
	}
	defer consumer.Close()

	go consumeMessages(consumer)

	go func() {
		router := gin.Default()
		book.RegisterHandlers(router, bookSvc)

		// Ping context
		router.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Iniciar servidor
		if err := router.Run(":8080"); err != nil {
			panic(fmt.Sprintf("Error running HTTP Server: %s", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Aplicación detenida")
}

func consumeMessages(consumer sarama.Consumer) {
	// Crear un consumidor de la partición 0 del topic
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %s\n", err)
		return
	}
	defer partitionConsumer.Close()

	// Proceso para escuchar continuamente mensajes de Kafka
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Mensaje recibido: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error consumiendo mensaje: %s\n", err)
		}
	}
}
