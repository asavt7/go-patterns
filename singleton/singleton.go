package singleton

import (
	"fmt"
	"sync"
)

type Service interface {
	Inc()
	Get() int
}

type singleton struct {
	sync.RWMutex
	count int
}

var instance *singleton
var once = sync.Once{}

func GetServiceInstance() Service {
	once.Do(func() {
		fmt.Println("Create instance")
		instance = &singleton{}
	})
	return instance
}

func (s *singleton) Inc() {
	s.Lock()
	defer s.Unlock()
	s.count += 1
}

func (s *singleton) Get() int {
	s.RLock()
	defer s.RUnlock()
	return s.count
}
