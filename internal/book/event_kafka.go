package book

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type kafkaBookEventPropagator struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewKafkaBookEventPropagator(producer sarama.AsyncProducer, topic string) BookEventPropagator {
	return &kafkaBookEventPropagator{
		producer: producer,
		topic:    topic,
	}
}

func (p *kafkaBookEventPropagator) propagateBookCreated(bookCreated *BookCreated) {
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
