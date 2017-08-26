package workerpool

import (
	"sync"
	"time"
	"github.com/tarantool/go-tarantool"
)

type Task interface {
	Execute()
}

type Result struct {

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
	//tarantoolConn := InitTarantool()
	switch e.Command {
	case "CreateSpace":
		{

		}
	case "DeleteSpace":
		{

		}
	case "InsertTuple":
		{
			//tarantool2.InsertTuple(tarantoolConn, e.Tuple_id, e.Name_space, e.User_id, e.Data)
		}
	case "SelectTuple":
		{
			//tuple, _ := tarantool2.SelectTuple(tarantoolConn, e.Tuple_id, e.Name_space, e.User_id)
		}
	case "DeleteTuple":
		{

		}
	case "UpdateTuple":
		{

		}
	case "SelectAllTuples":
		{

		}
	}
}

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
