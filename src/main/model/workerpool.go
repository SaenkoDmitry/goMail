package workerpool

import (
	"fmt"
	"sync"
	"time"

	tarantool "github.com/tarantool/go-tarantool"
)

type Task interface {
	Execute()
}

type Pool struct {
	mu    sync.Mutex
	size  int
	tasks chan Task
	kill  chan struct{}
	wg    sync.WaitGroup
}

type TarantoolTask struct {
	command    string
	tuple_id   uint32
	name_space string
	user_id    uint64
	data       interface{}
}

func NewPool(size int) *Pool {
	pool := &Pool{
		tasks: make(chan Task, 128),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.kill:
			return
		}
	}
}

func (p *Pool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}

func (p *Pool) Close() {
	close(p.tasks)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Exec(task Task) {
	p.tasks <- task
}

func (e TarantoolTask) Execute() {
	//name_spaceT := convertToNameInTarantool(e.name_space, e.user_id)
	name_spaceT := "examples"
	tarantoolConn := InitTarantool()
	_, err := tarantoolConn.Insert(name_spaceT, []interface{}{e.tuple_id, e.data})
	if err != nil {
		fmt.Println(err.Error())
	}
}

//func convertToNameInTarantool(name string, user_id uint64) string {
//	return name + fmt.Sprintf("%v", user_id)
//}

func InitTarantool() *tarantool.Connection {

	server := "127.0.0.1:3302"
	opts := tarantool.Opts{
		Timeout:       500 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          "test",
		Pass:          "12345",
	}
	conn, err := tarantool.Connect(server, opts)
	if err != nil {
		panic(err)
	}
	return conn
}
