package mapdb_test

import (
	"strconv"
	"testing"

	"github.com/lysShub/mapdb"
)

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Write()
	}

	// 一次循环有30次写入，一次删除操作
	// 6561 ns/op	     161 B/op	      20 allocs/op // 不记录日志
	// 573078 ns/op	    1229 B/op	      25 allocs/op // 记录日志225
}

var db *mapdb.Db

func init() {
	var err error

	db, err = mapdb.NewMapDb(func(d *mapdb.Db) {
		d.Name = "test"
		d.Log = true
	})
	if err != nil {
		panic(err)
	}
}

var C = map[string]string{
	"a": "1a",
	"b": "1b",
	"c": "1c",
}

var index int = 0

// Write
func Write() {
	tmp := index + 10
	for ; index < tmp; index++ {
		db.UpdateRow(strconv.Itoa(index), C)
		db.DeleteRow(strconv.Itoa(index))
	}
}
