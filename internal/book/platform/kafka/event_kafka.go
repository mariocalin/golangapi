package kafka

import (
	"encoding/json"
	"library-api/internal/book/application"

	"github.com/IBM/sarama"
)

type kafkaBookEventPropagator struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewKafkaBookEventPropagator(producer sarama.AsyncProducer, topic string) application.EventPropagator {
	return &kafkaBookEventPropagator{
		producer: producer,
		topic:    topic,
	}
}

func (p *kafkaBookEventPropagator) PropagateBookCreated(bookCreated *application.CreatedEvent) {
	payload := struct {
		Id string `json:"id"`
	}{
		Id: bookCreated.Id.String(),
	}

	msgBytes, _ := json.Marshal(payload)

	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}
}

func (p *kafkaBookEventPropagator) Close() {
	p.producer.Close()
}
