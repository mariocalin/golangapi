package book

type BookCreated struct {
	Id   BookId
	Name Name
}

type BookEventPropagator interface {
	PropagateBookCreated(bookCreated *BookCreated)
}

type BookEventConsumer interface {
	BindBookCreated(func(*BookCreated))
	StartConsuming()
}
