package book

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
	for bookCreated := range consumer.bookCreatedChannel {
		for _, callback := range consumer.bookCreatedCallbacks {
			callback(bookCreated)
		}
	}
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
