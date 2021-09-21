package tq

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type TQ struct {

	// 达到任务执行时间时返回对应的Ts.P; 请使用for-range及时读取, 否则会阻塞以致影响后续任务
	MQ chan interface{}

	addChan chan Ts    // 增加任务管道
	works   []*work    // 记录任务
	lock    sync.Mutex //
	capChan int        // chan容量
	writer  io.Writer  //
}

// work 表示一个工作
type work struct {
	c       chan Ts   // 任务队列
	endTime time.Time // 队列中最后任务执行时间
}

// Ts 表示一个任务
type Ts struct {
	T time.Time   // 设定执行时间
	P interface{} // 执行时MQ返回的数据
}

// NewTQ
// 	@ioOut: 被丢弃任务输出日志；对于MQ中未被及时读取的的数据，唯一的操作是将其丢弃，但避免静默处理，因此需要打日志
func NewTQ(ioOut ...io.Writer) *TQ {
	var t = new(TQ)
	if len(ioOut) != 0 {
		t.writer = io.MultiWriter(ioOut...)
	} else {
		t.writer = os.Stdout
	}
	t.run()
	return t
}

// Run 启动
func (t *TQ) run() {
	t.capChan = 128
	t.MQ = make(chan interface{}, t.capChan)
	t.addChan = make(chan Ts, t.capChan)
	t.works = make([]*work, 0, 64)

	// 分发任务
	go func() {
		var r Ts
		var flag bool

		for r = range t.addChan {

			flag = false
			for i := 0; i < len(t.works); i++ {
				if r.T.After(t.works[i].endTime) && len(t.works[i].c) < cap(t.works[i].c) {
					t.works[i].endTime = r.T
					t.works[i].c <- r
					flag = true
					break
				}
			}

			// 需要新建工作
			if !flag {
				var w = new(work)
				w.c, w.endTime = make(chan Ts, 1024), r.T

				t.lock.Lock()
				t.works = append(t.works, w)
				t.lock.Unlock()
				w.c <- r

				go t.exec(w) // 运行work

				// 维护works，释放过多空闲的works
				if len(t.works) > 16 {
					t.lock.Lock()
					for i := 3; i < len(t.works); i++ {
						if len(t.works[i].c) == 0 {
							close(t.works[i].c)
							t.works = append(t.works[:i], t.works[i+1:]...)
						}
					}
					t.lock.Unlock()
				}
			}

		}

	}()
}

// Add 增加任务
// 	存在阻塞的可能！！
func (t *TQ) Add(r Ts) error {
	t.addChan <- r
	return nil
}

// Drop 销毁
func (t *TQ) Drop() {
	for i := 0; i < len(t.works); i++ {
		close(t.works[i].c)
	}
	close(t.MQ)
}

// exec 执行work
func (t *TQ) exec(w *work) {
	defer func() {
		recover() // 从w.c读取到任务，在执行通知之前先执行了Drop；会导致panic
	}()
	var ts Ts

	for ts = range w.c {

		time.Sleep(time.Until(ts.T)) //延时

		t.lock.Lock()
		if len(t.MQ) < t.capChan {
			t.MQ <- ts.P
		} else {
			t.writer.Write([]byte(fmt.Sprintln(<-t.MQ)))
			t.MQ <- ts.P
		}
		t.lock.Unlock()
	}
}
