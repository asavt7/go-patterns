package singleton

import "testing"

func TestSingletonCh(t *testing.T) {

	increment := 1000
	startCount := GetServiceInstanceCh().Get()
	t.Run("one instance", func(t *testing.T) {
		for i := 0; i < increment; i++ {
			instance := GetServiceInstanceCh()
			instance.Inc()
		}
		if got := GetServiceInstanceCh().Get() - startCount; got != increment {
			t.Errorf("expected %d, but got %d", increment, got)
		}
	})
}

func BenchmarkGetServiceInstanceCh(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetServiceInstanceCh().Inc()
	}
}
