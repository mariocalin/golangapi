package memory

import (
	"library-api/internal/book/application"
	"log"
)

type ChannelBookConsumer struct {
	bookCreatedCallbacks []func(*application.CreatedEvent)
	bookCreatedChannel   chan *application.CreatedEvent
}

func NewChannelBookConsumer(bookCreatedChannel chan *application.CreatedEvent) *ChannelBookConsumer {
	return &ChannelBookConsumer{
		bookCreatedChannel: bookCreatedChannel,
	}
}

func (consumer *ChannelBookConsumer) BindBookCreated(callback func(*application.CreatedEvent)) {
	consumer.bookCreatedCallbacks = append(consumer.bookCreatedCallbacks, callback)
}

func (consumer *ChannelBookConsumer) StartConsuming() {
	log.Println("Listening to events")

	for bookCreated := range consumer.bookCreatedChannel {
		log.Println("Event received", bookCreated)

		for _, callback := range consumer.bookCreatedCallbacks {
			callback(bookCreated)
		}
	}
}

func (consumer *ChannelBookConsumer) Close() {
	log.Println("Consumer will not close the channel")
}

type channelBookEventPropagator struct {
	bookCreatedChannel chan *application.CreatedEvent
}

func NewChannelBookEventPropagator(bookCreatedChannel chan *application.CreatedEvent) *channelBookEventPropagator {
	return &channelBookEventPropagator{
		bookCreatedChannel: bookCreatedChannel,
	}
}

func (propagator *channelBookEventPropagator) PropagateBookCreated(bookCreated *application.CreatedEvent) {
	propagator.bookCreatedChannel <- bookCreated
}

func (propagator *channelBookEventPropagator) Close() {
	log.Println("Closing propagator gracefully")
	close(propagator.bookCreatedChannel)
}
