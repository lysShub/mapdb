package store

import (
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

type Store struct {
	// 一个文件只存储一个表

	Path   string // 数据存储路径
	handle *bolt.DB
}

// OpenDb open
// 	如果数据文件已经存在，则继续写入数据
func OpenDb(dbFilePath string) (*Store, error) {
	var db *bolt.DB
	var err error
	db, err = bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	} else {
		var b = new(Store)
		b.handle = db
		return b, nil
	}
}

// Close 关闭
func (b *Store) Close() {
	b.handle.Close()
}

// UpdateRow 更新行
func (d *Store) UpdateRow(id string, p map[string]string) error {

	return d.handle.Update(func(tx *bolt.Tx) error {
		if b, err := tx.CreateBucketIfNotExists([]byte(id)); err != nil {
			return err
		} else {
			for k, v := range p {
				if err = b.Put([]byte(k), []byte(v)); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (d *Store) DeleteRow(id string) error {
	return d.handle.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(id))
	})
}

func (b *Store) ReadRow(id string) map[string]string {
	var r map[string]string = make(map[string]string)
	b.handle.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(id)); b == nil {
			r = nil // 没有此行
		} else {
			for k, v := b.Cursor().First(); k != nil; k, v = b.Cursor().Next() {
				r[string(k)] = string(v)
			}
		}
		return nil
	})

	return r
}

// 方法可执行文件(不包括)所在路径
func getExePath() string {
	ex, err := os.Executable()
	if err != nil {
		exReal, err := filepath.EvalSymlinks(ex)
		if err != nil {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				return "./"
			}
			return dir
		}
		return filepath.Dir(exReal)
	}
	return filepath.Dir(ex)
}
