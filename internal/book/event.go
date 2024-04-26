package book

type BookCreated struct {
	Id BookId
}

type BookEventPropagator interface {
	PropagateBookCreated(bookCreated *BookCreated)
}
