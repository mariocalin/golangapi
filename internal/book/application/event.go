package application

import (
	"library-api/internal/book/domain"
)

type CreatedEvent struct {
	Id   *domain.Id
	Name *domain.Name
}

//go:generate mockery --name=EventPropagator --outpkg=book --dir=. --filename=event_propagator_mock.go --output=. --structname=MockEventPropagator
type EventPropagator interface {
	PropagateBookCreated(bookCreated *CreatedEvent)
	Close()
}

//go:generate mockery --name=EventConsumer --outpkg=book --dir=. --filename=event_consumer_mock.go --structname=MockEventConsumer
type EventConsumer interface {
	BindBookCreated(func(*CreatedEvent))
	StartConsuming()
	Close()
}
