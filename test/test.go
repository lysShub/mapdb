package main

import (
	"fmt"
	"time"

	"github.com/lysShub/mapdb"
)

func main() {
	db := new(mapdb.Db)
	db.Init()
	time.Sleep(time.Millisecond * 20)

	db.Ct("1", map[string]string{
		"a": "1a",
		"b": "1b",
		"c": "1c",
	}, time.Second)

	fmt.Println(db.R("1", "b"))
	time.Sleep(time.Millisecond * 1010)
	fmt.Println(db.R("1", "b"))

}
