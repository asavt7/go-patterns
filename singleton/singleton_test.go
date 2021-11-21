package singleton

import (
	"sync"
	"testing"
)

func TestSingleton(t *testing.T) {

	increment := 1000
	startCount := GetServiceInstance().Get()
	t.Run("one instance", func(t *testing.T) {
		for i := 0; i < increment; i++ {
			instance := GetServiceInstance()
			instance.Inc()
		}
		if got := GetServiceInstance().Get() - startCount; got != increment {
			t.Errorf("expected %d, but got %d", increment, got)
		}
	})
}

func TestSingleton_ConcurrentAccess(t *testing.T) {

	increment := 1000
	startCount := GetServiceInstance().Get()
	t.Run("concurrent access", func(t *testing.T) {
		wg := sync.WaitGroup{}
		for i := 0; i < increment; i++ {
			wg.Add(1)
			go func() {
				instance := GetServiceInstance()
				instance.Inc()
				wg.Done()
			}()
		}
		wg.Wait()
		if got := GetServiceInstance().Get() - startCount; got != increment {
			t.Errorf("expected %d, but got %d", increment, got)
		}
	})
}

func BenchmarkSingleton_RWMutex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetServiceInstance().Inc()
	}
}