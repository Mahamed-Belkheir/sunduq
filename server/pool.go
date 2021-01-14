package server

import (
	"github.com/Mahamed-Belkheir/sunduq"
)

//Pool struct handles pooling the connections, and holds all the queues
type Pool struct {
	workers chan chan sunduq.Handler
	queue   chan sunduq.Handler
	close   chan bool
}

//NewPool creates a new pool instance, initalizating its channels
func NewPool(workerLimit, queueLimit int) Pool {
	return Pool{
		make(chan chan sunduq.Handler, workerLimit),
		make(chan sunduq.Handler, queueLimit),
		make(chan bool),
	}
}

//Init starts the pool workers
func (p Pool) Init() {
	for i := 0; i < cap(p.workers); i++ {
		NewWorker(p.workers)
	}
	go func() {
		for {
			select {
			case <-p.close:
				return

			case job := <-p.queue:
				worker := <-p.workers
				worker <- job
			}
		}
	}()
}

//NewWorker creates a new worker goroutine and adds it to the worker queue
func NewWorker(workerQueue chan chan sunduq.Handler) {
	recieveChannel := make(chan sunduq.Handler)
	workerQueue <- recieveChannel
	go func() {
		for {
			handler, active := <-recieveChannel
			if !active {
				return
			}
			handler.Run()
			workerQueue <- recieveChannel
		}
	}()
}

//Close closes the pool and shuts down all workers
func (p Pool) Close() {
	close(p.queue)
	p.close <- true

	for i := 0; i < cap(p.workers); i++ {
		worker := <-p.workers
		close(worker)
	}
}
