package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Effector func(context.Context) (string, error)

func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	var tokens = max
	var once sync.Once
	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		once.Do(func() {
			ticker := time.NewTicker(d)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						fmt.Println("REFILL")
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})
		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}
		tokens--
		return e(ctx)
	}
}

func main() {
	e := Effector(func(ctx context.Context) (string, error) {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("call")
		return "", nil
	})

	t := Throttle(e, 5, 5, time.Second)

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*100)
			s, err := t(ctx)
			fmt.Println("results ", s, err)
		}()
	}
	wg.Wait()

	time.Sleep(2 * time.Second)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*100)
			s, err := t(ctx)
			fmt.Println("results ", s, err)
		}()
	}
	wg.Wait()

}
