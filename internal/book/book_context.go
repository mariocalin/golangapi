package book

import (
	"library-api/common"
	"library-api/external/sqlite3"
	"library-api/internal/book/application"
	"library-api/internal/book/platform/memory"
	"library-api/internal/book/platform/sqlite"
	"sync"
)

var (
	once        sync.Once
	bookChannel chan *application.CreatedEvent
)

func CreateBookRepositoryInstance(sqliContext *sqlite3.Sqlite3Context, dateHandler *common.DateHandler) application.BookRepository {
	return sqlite.NewSqlite3BookRepository(sqliContext.Db, dateHandler)
}

func CreateBookEventPropagatorInstance() application.EventPropagator {
	once.Do(initChannel)
	return memory.NewChannelBookEventPropagator(bookChannel)
}

func CreateBooEventConsumerInstance() application.EventConsumer {
	once.Do(initChannel)
	return memory.NewChannelBookConsumer(bookChannel)
}

func initChannel() {
	bookChannel = make(chan *application.CreatedEvent)
}
