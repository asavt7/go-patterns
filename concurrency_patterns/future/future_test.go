package future

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDoFutureWithContext(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
	res := (DoFutureWithContext(ctx, "das")).Get()
	fmt.Println(res)
}
