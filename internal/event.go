package internal

import "context"

type BookCreated struct {
	Id   BookId
	Name BookName
}

type BookEventPropagator interface {
	PropagateBookCreated(ctx context.Context, event BookCreated) error
	Close(ctx context.Context)
}

type BookEventConsumer interface {
	BindBookCreated(func(ctx context.Context, event BookCreated))
	StartConsuming(ctx context.Context)
	Close(ctx context.Context)
}
