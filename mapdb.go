package mapdb

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lysShub/mapdb/store"
	"github.com/lysShub/tq"
)

/* 使用map数据结构实现的缓存简单数据库 */

type Db struct {
	// 一个实例存储一个表
	// 支持行的TTL

	Name string // 名称, 必须参数
	Log  bool   //在数据TTL删除之前进行持久化, 采用的boltdb记录

	m map[string]map[string]string

	lock sync.RWMutex // 写入锁
	q    *tq.TQ       // 时间任务队列, 用于TTL

	s *store.Store
}

// NewMapDb
func NewMapDb(config func(*Db) *Db) (*Db, error) {
	var d = new(Db)
	d = config(d)
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// init 初始化
func (d *Db) init() error {
	if d.Name == "" {
		return errors.New("invalid Db.Name")
	}
	if d.Log {
		path := getExePath() + `/` + d.Name
		var err error
		if d.s, err = store.OpenDb(path); err != nil {
			return err
		}
	}
	d.m = make(map[string]map[string]string)

	d.q = tq.NewTQ() // 时间任务队列
	var r interface{}
	go func() {
		for r = range d.q.MQ {
			if v, ok := r.(string); ok {
				d.lock.RLock()
				delete(d.m, v)
				d.lock.RUnlock()
			}
		}
	}()
	return nil
}

// R 查，没有将会返回空字符串
func (d *Db) R(id, field string) string {
	d.lock.RLock()
	var r string = d.m[id][field]
	d.lock.RUnlock()
	return r
}

// U 更新值
func (d *Db) U(id, field, value string) {
	d.lock.RLock()
	if d.m[id] == nil {
		d.m[id] = map[string]string{}
		d.m[id][field] = value
	} else {
		d.m[id][field] = value
	}
	d.lock.RUnlock()
}

// UpdateRow 更新一行
func (d *Db) UpdateRow(id string, t map[string]string, ttl ...time.Duration) {

	if d.m[id] != nil {
		d.lock.RLock()
		for k, v := range t {
			d.m[id][k] = v
		}
		d.lock.RUnlock()
	} else {
		d.lock.RLock()
		d.m[id] = t
		d.lock.RUnlock()
	}

	// ttl
	if len(ttl) != 0 {
		d.q.Add(tq.Ts{
			T: time.Now().Add(ttl[0]),
			P: id,
		})
	}
}

// DeleteRow 删除一行
func (d *Db) DeleteRow(id string) {
	d.lock.RLock()
	delete(d.m, id)
	d.lock.RUnlock()
}

// ExitRow 行是否存在
func (d *Db) ExitRow(id string) bool {
	d.lock.RLock()
	if d.m[id] == nil {
		d.lock.RUnlock()
		return false
	}
	d.lock.RUnlock()
	return true
}

// Drop 销毁
func (d *Db) Drop() {
	// d.done = true

	d.q.Drop()
	d.m = nil
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
