package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("did_db/dids", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Put([]byte("Key #1"), []byte("Value #1"), nil)

	data, err := db.Get([]byte("Key #1"), nil)
	strData := string(data[:len(data)])
	fmt.Println("Data: ", strData)
}
