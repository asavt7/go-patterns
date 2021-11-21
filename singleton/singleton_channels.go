package singleton

import (
	"fmt"
	"sync"
)

type singletonCh struct {
	count int
	get   chan chan int
	inc   chan struct{}
	done  chan struct{}
}

func newSingletonCh() *singletonCh {
	get := make(chan chan int)
	inc := make(chan struct{})
	done := make(chan struct{})
	s := &singletonCh{
		count: 0,
		get:   get,
		inc:   inc,
		done:  done,
	}

	go func() {
		for {
			select {
			case getRsCh := <-s.get:
				getRsCh <- s.count
			case <-done:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-s.inc:
				s.count++
			case <-done:
				return
			}
		}
	}()

	return s
}

var instanceCh *singletonCh
var onceCh = sync.Once{}

func GetServiceInstanceCh() Service {
	onceCh.Do(func() {
		fmt.Println("Create instance")
		instanceCh = newSingletonCh()
	})
	return instanceCh
}

func (s *singletonCh) Inc() {
	s.inc <- struct{}{}
}

func (s *singletonCh) Get() int {
	ch := make(chan int)
	s.get <- ch
	return <-ch
}
