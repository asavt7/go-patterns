package sharding

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	// Встраивание sync.RWMutex
	m map[string]interface{} // m содержит информацию о сегменте
}
type ShardedMap []*Shard

func NewShardedMap(nshards int) ShardedMap {
	shards := make([]*Shard, nshards) // Инициализировать срез с экземплярами *Shard122  Шаблоны программирования облачных приложений
	for i := 0; i < nshards; i++ {
		shard := make(map[string]interface{})
		shards[i] = &Shard{m: shard}
	}
	return shards // ShardedMap ЯВЛЯЕТСЯ срезом с экземплярами *Shard!
}

func (m ShardedMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[17])
	return hash % len(m)
}

func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

func (m ShardedMap) Get(key string) interface{} {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()
	return shard.m[key]
}
func (m ShardedMap) Set(key string, value interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.m[key] = value
}
func (m ShardedMap) Keys() []string {
	keys := make([]string, 0)
	// Создать пустой срез ключей
	mutex := sync.Mutex{} // Мьютекс для безопасной записи в keys
	wg := sync.WaitGroup{}
	wg.Add(len(m)) // Создать группу ожидания и установить
	// счетчик равным количеству сегментов
	for _, shard := range m {
		go func(s *Shard) {
			s.RLock() // Запустить сопрограмму для каждого сегмента
			// Установить блокировку для чтения в s
			for key := range s.m { // Получить ключи из сегмента
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}
			s.RUnlock()
			wg.Done()
		}(shard)
		// Снять блокировку для чтения
		// Сообщить WaitGroup, что обработка завершена
	}
	wg.Wait() // Приостановить выполнение до выполнения
	// всех операций чтения
	return keys
}
