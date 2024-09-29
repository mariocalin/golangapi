package memory

import (
	"context"
	"library-api/internal"
)

type channelBookConsumer struct {
	bookCreatedCallbacks []func(ctx context.Context, event internal.BookCreated)
	bookCreatedChannel   chan internal.BookCreated
}

func NewChannelBookConsumer(bookCreatedChannel chan internal.BookCreated) *channelBookConsumer {
	return &channelBookConsumer{
		bookCreatedChannel: bookCreatedChannel,
	}
}

func (consumer *channelBookConsumer) BindBookCreated(callback func(ctx context.Context, event internal.BookCreated)) {
	consumer.bookCreatedCallbacks = append(consumer.bookCreatedCallbacks, callback)
}

func (consumer *channelBookConsumer) StartConsuming(ctx context.Context) {
	for bookCreated := range consumer.bookCreatedChannel {

		for _, callback := range consumer.bookCreatedCallbacks {
			callback(ctx, bookCreated)
		}
	}
}

func (consumer *channelBookConsumer) Close(ctx context.Context) {
	close(consumer.bookCreatedChannel)
}

type ChannelBookEventPropagator struct {
	bookCreatedChannel chan internal.BookCreated
}

func NewChannelBookEventPropagator(bookCreatedChannel chan internal.BookCreated) *ChannelBookEventPropagator {
	return &ChannelBookEventPropagator{
		bookCreatedChannel: bookCreatedChannel,
	}
}

func (propagator *ChannelBookEventPropagator) PropagateBookCreated(ctx context.Context, bookCreated internal.BookCreated) error {
	select {
	case propagator.bookCreatedChannel <- bookCreated:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (propagator *ChannelBookEventPropagator) Close(ctx context.Context) {
	close(propagator.bookCreatedChannel)
}
