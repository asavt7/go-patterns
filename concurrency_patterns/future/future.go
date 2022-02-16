package future

import (
	"context"
)

func futureFunc(x string) <-chan string {
	out := make(chan string)
	go func() {
		// do some work and send result in channel
		out <- x + x + x
		close(out)
	}()
	return out
}

// Example with wrapper

type Result struct {
	Res string
	Err error
}

type Future struct {
	out chan Result
}

func DoFuture(in string) *Future {
	out := make(chan Result)
	go func() {
		//do some work
		out <- Result{
			Res: in,
			Err: nil,
		}
		close(out)
	}()
	return &Future{out: out}
}

func DoFutureWithContext(ctx context.Context, in string) *Future {
	out := make(chan Result)
	go func() {
		defer close(out)

		c := make(chan Result)

		go func() {
			defer close(c)
			//do some work
			c <- Result{
				Res: in,
				Err: nil,
			}
		}()

		select {
		case <-ctx.Done():
			out <- Result{
				Res: "",
				Err: ctx.Err(),
			}
		case r := <-c:
			out <- r
		}
	}()
	return &Future{out: out}
}

func (f *Future) Get() Result {
	return <-f.out
}
