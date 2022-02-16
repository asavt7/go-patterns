package breaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Circuit Полезная работа (вызов сервиса)
type Circuit func(context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex
	return func(ctx context.Context) (string, error) {
		m.RLock()
		// Установить "блокировку чтения"
		d := consecutiveFailures - int(failureThreshold)
		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}
		m.RUnlock()                   // Освободить блокировку чтения
		response, err := circuit(ctx) // Послать запрос, как обычно
		m.Lock()
		defer m.Unlock()         // Заблокировать общие ресурсы
		lastAttempt = time.Now() // Зафиксировать время попытки
		if err != nil {          // Если Circuit вернула ошибку,Шаблоны стабильности
			consecutiveFailures++
			return response, err
			// увеличить счетчик ошибок
			// и вернуть ошибку
		}
		consecutiveFailures = 0
		// Сбросить счетчик ошибок
		return response, nil
	}
}
