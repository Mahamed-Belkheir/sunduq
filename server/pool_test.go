package server

import (
	"sync"
	"testing"
)

func TestPoolCreation(t *testing.T) {
	pool := NewPool(10, 5)
	pool.Init()
	pool.Close()
}

func TestPoolWork(t *testing.T) {
	pool := NewPool(10, 5)
	pool.Init()
	mock := newHandler(1000)
	for i := 0; i < 1000; i++ {
		pool.queue <- mock
	}
	mock.wg.Wait()
	pool.Close()
}

type mockHandler struct {
	wg *sync.WaitGroup
}

func newHandler(wait int) mockHandler {
	m := mockHandler{
		&sync.WaitGroup{},
	}
	m.wg.Add(wait)
	return m
}

func (m mockHandler) Run() {
	m.wg.Done()
}
