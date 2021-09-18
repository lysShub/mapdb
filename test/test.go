package main

import (
	"fmt"
	"time"

	"github.com/lysShub/mapdb"
)

func main() {

	// d := new(bolt.Bolt)
	// fmt.Println(d.OpenDb(`D:\Desktop\garbage\data.db`))

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
