package book

type BookCreated struct {
	Id BookId
}

type BookEventPropagator interface {
	propagateBookCreated(bookCreated *BookCreated)
}
