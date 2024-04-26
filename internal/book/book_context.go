package book

var event = "CHANNEL"

var bookChannel chan *BookCreated = make(chan *BookCreated)

func CreateBookRepositoryInstance() BookRepository {
	return NewSqlite3BookRepository("data.sqlite3")
}

func CreateBookEventPropagatorInstance() BookEventPropagator {
	if event == "CHANNEL" {
		return NewChannelBookEventPropagator(bookChannel)
	}

	panic("BookEventPropagator not initialized")
}

func CreateBooEventConsumerInstance() BookEventConsumer {
	if event == "CHANNEL" {
		return NewChannelBookConsumer(bookChannel)
	}

	panic("BookEventConsumer not initialized")
}
