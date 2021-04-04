package mapdb

import (
	"sync"
	"time"

	"github.com/lysShub/tq"
)

/* 使用map数据结构实现的缓存简单数据库 */

// type Handle map[string]map[string]string

// 时间任务队列
var Q *tq.TQ

type Db struct {
	// 使用map存储数据，仅单线程安全
	H map[string]map[string]string

	lock sync.RWMutex // 写入锁
}

// Init 阻塞函数，请用协程启动此程序以初始化
func (d *Db) Init() {
	d.H = make(map[string]map[string]string)

	Q = new(tq.TQ)
	go Q.Run()
	tq.InitEnd.Wait()

	var r interface{}
	go func() {
		for {
			r = (<-(Q.MQ))
			v, ok := r.(string)
			if ok {
				delete(d.H, v)
			}
		}
	}()

}

// D 删除
func (d *Db) D(id string) {
	d.lock.RLock()
	delete(d.H, id)
	d.lock.RUnlock()
}

// R 读取(没有将会返回空字符串)
func (d *Db) R(id, field string) string {
	return d.H[id][field]
}

// U 更新值(表不存在将不会记录)
func (d *Db) U(id, field, value string) {
	var t map[string]string = make(map[string]string)
	t = d.H[id]
	if t == nil {
		return
	}
	t[field] = value
	d.lock.RLock()
	d.H[id] = t
	d.lock.RUnlock()

}

// Ut 更新表(表不存在将不会记录)
func (d *Db) Ut(id string, t map[string]string) {
	if d.H[id] == nil {
		return
	}
	d.lock.RLock()
	d.H[id] = t
	d.lock.RUnlock()
}

// Ct 创造表(ttl及表的生存时间，请使用UTC时间)
func (d *Db) Ct(id string, t map[string]string, ttl ...time.Duration) {

	var ct time.Time
	d.lock.RLock()
	d.H[id] = t
	if len(ttl) != 0 {
		ct = time.Now().Add(ttl[0])
	}
	d.lock.RUnlock()

	// ttl
	if len(ttl) != 0 {
		Q.Add(tq.Ts{
			T: ct,
			P: id,
		})
	}

}

// Et 表是否存在(异常情况返回false)
func (d *Db) Et(id string) bool {
	if d.H[id] == nil {
		return false
	}
	return true
}