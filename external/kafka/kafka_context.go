package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// Local kafka info
var (
	brokers = []string{"localhost:9092"}
	topic   = "test_topic"
)

type KafkaContext struct {
	Consumer sarama.Consumer
	Producer sarama.AsyncProducer
	Topic    string
}

func (context *KafkaContext) close() {
	context.Producer.Close()
	context.Consumer.Close()
}

func CreateKafkaContext() *KafkaContext {
	producer := createProducer()
	consumer := createConsumer()

	return &KafkaContext{
		Consumer: consumer,
		Producer: producer,
		Topic:    topic,
	}
}

func createProducer() sarama.AsyncProducer {
	// Configuración del productor de Kafka
	producer, err := sarama.NewAsyncProducer(brokers, nil)
	if err != nil {
		panic(fmt.Sprintf("Error creating Kafka producer: %s", err))
	}

	return producer
}

func createConsumer() sarama.Consumer {
	// Configuración del consumidor de Kafka
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(fmt.Sprintf("Error creando Kafka consumer: %s", err))
	}

	return consumer
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
