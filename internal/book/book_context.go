package book

import (
	"library-api/common"
	"library-api/external/sqlite3"
	"sync"
)

var (
	once        sync.Once
	bookChannel chan *BookCreated
)

func CreateBookRepositoryInstance(sqliContext *sqlite3.Sqlite3Context, dateHandler *common.DateHandler) BookRepository {
	return NewSqlite3BookRepository(sqliContext.Db, dateHandler)
}

func CreateBookEventPropagatorInstance() BookEventPropagator {
	once.Do(initChannel)
	return NewChannelBookEventPropagator(bookChannel)
}

func CreateBooEventConsumerInstance() BookEventConsumer {
	once.Do(initChannel)
	return NewChannelBookConsumer(bookChannel)
}

func initChannel() {
	bookChannel = make(chan *BookCreated)
}
