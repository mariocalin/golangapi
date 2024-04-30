package book

type BookCreated struct {
	Id   *BookId
	Name *Name
}

type BookEventPropagator interface {
	PropagateBookCreated(bookCreated *BookCreated)
	Close()
}

type BookEventConsumer interface {
	BindBookCreated(func(*BookCreated))
	StartConsuming()
	Close()
}
