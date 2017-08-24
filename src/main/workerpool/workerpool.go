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

var MainPool *Pool

type Pool struct {
	Mu    sync.Mutex
	Size  int
	Tasks chan Task
	Kill  chan struct{}
	Wg    sync.WaitGroup
}

type TarantoolTask struct {
	Command    string
	Tuple_id   uint64
	Name_space string
	User_id    uint64
	Data       interface{}
}

func NewPool(size int) *Pool {
	pool := &Pool{
		Tasks: make(chan Task, 128),
		Kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

func (p *Pool) worker() {
	defer p.Wg.Done()
	for {
		select {
		case task, ok := <-p.Tasks:
			if !ok {
				return
			}
			task.Execute()
		case <-p.Kill:
			return
		}
	}
}

func (p *Pool) Resize(n int) {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	for p.Size < n {
		p.Size++
		p.Wg.Add(1)
		go p.worker()
	}
	for p.Size > n {
		p.Size--
		p.Kill <- struct{}{}
	}
}

func (p *Pool) Close() {
	close(p.Tasks)
}

func (p *Pool) Wait() {
	p.Wg.Wait()
}

func (p *Pool) Exec(task Task) {
	p.Tasks <- task
}

func (e TarantoolTask) Execute() {
	//name_spaceT := convertToNameInTarantool(e.name_space, e.user_id)
	name_spaceT := "examples"
	tarantoolConn := InitTarantool()
	_, err := tarantoolConn.Insert(name_spaceT, []interface{}{e.Tuple_id, e.Data})
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
