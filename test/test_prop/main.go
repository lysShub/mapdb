package main

import (
	"errors"

	"github.com/lysShub/mapdb"
)

var C = map[string]string{
	"a": "1a",
	"b": "1b",
	"c": "1c",
}

var db *mapdb.Db

func init() {
	db = mapdb.NewMapDb()
}

// 11 次操作
func Comprehensive(id string) error {
	db.Ut(id, C)

	db.U(id, "a", "2a")
	db.U(id, "b", "2b")
	db.U(id, "c", "2c")

	if db.R(id, "a") != "2a" || db.R(id, "b") != "2b" || db.R(id, "c") != "2c" {
		return errors.New("error")
	}

	db.D(id)

	if db.Et(id) {
		return errors.New("error2")
	}
	return nil
}
