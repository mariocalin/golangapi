package book

import "log"

type channelBookConsumer struct {
	bookCreatedCallbacks []func(*BookCreated)
	bookCreatedChannel   chan *BookCreated
}

func NewChannelBookConsumer(bookCreatedChannel chan *BookCreated) *channelBookConsumer {
	return &channelBookConsumer{
		bookCreatedChannel: bookCreatedChannel,
	}
}

func (consumer *channelBookConsumer) BindBookCreated(callback func(*BookCreated)) {
	consumer.bookCreatedCallbacks = append(consumer.bookCreatedCallbacks, callback)
}

func (consumer *channelBookConsumer) StartConsuming() {
	log.Println("Listening to events")

	for bookCreated := range consumer.bookCreatedChannel {
		log.Println("Event received", bookCreated)

		for _, callback := range consumer.bookCreatedCallbacks {
			callback(bookCreated)
		}
	}
}

func (consumer *channelBookConsumer) Close() {
	log.Println("Consumer will not close the channel")
}

type channelBookEventPropagator struct {
	bookCreatedChannel chan *BookCreated
}

func NewChannelBookEventPropagator(bookCreatedChannel chan *BookCreated) *channelBookEventPropagator {
	return &channelBookEventPropagator{
		bookCreatedChannel: bookCreatedChannel,
	}
}

func (propagator *channelBookEventPropagator) PropagateBookCreated(bookCreated *BookCreated) {
	propagator.bookCreatedChannel <- bookCreated
}

func (propagator *channelBookEventPropagator) Close() {
	log.Println("Closing propagator gracefully")
	close(propagator.bookCreatedChannel)
}
