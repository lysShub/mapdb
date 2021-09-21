package main

import (
	"fmt"
	"time"

	"github.com/lysShub/mapdb"
	"github.com/lysShub/mapdb/store"
)

func main() {

	s, err := store.OpenDb(`D:\OneDrive\code\go\mapdb\test\test.db`)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	fmt.Println(s.ReadRow("1"))
	return

	db, err := mapdb.NewMapDb(func(d *mapdb.Db) *mapdb.Db {
		d.Name = "test.db"
		d.Log = true
		return d
	})
	if err != nil {
		panic(err)
	}

	db.UpdateRow("1", map[string]string{
		"a": "a",
		"b": "b",
		"c": "c",
	}, time.Second)

	db.UpdateRow("2", map[string]string{
		"a1": "1a",
		"b1": "1b",
		"c1": "1c",
	}, time.Second)

	defer db.Drop()
}

func main1() {
	db, err := mapdb.NewMapDb(func(d *mapdb.Db) *mapdb.Db {
		d.Name = "test"
		return d
	})
	fmt.Println(err)

	db.UpdateRow("1", map[string]string{
		"a": "1a",
		"b": "1b",
		"c": "1c",
	}, time.Second)

	fmt.Println(db.R("1", "b"))
	time.Sleep(time.Millisecond * 1010)
	fmt.Println(db.R("1", "b"))
}
