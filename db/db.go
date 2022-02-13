package db

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

// index:$time:$id:attr = value
// chat:$time:$id:attr = value
type Database struct {
	db *badger.DB
}

func Open() *Database {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	return &Database{db}
}
func (db Database) Close() {
	//db.db.RunValueLogGC()
}
func (db Database) InsertChatMessage() {
	//db.db.Update(func(txn *badger.Txn) error {
	//	txn.
	//})
}
