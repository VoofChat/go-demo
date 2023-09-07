package demo4

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func doSomething(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("playing")
			return
		default:
			fmt.Println("I am working!")
			time.Sleep(time.Second)
		}
	}
}

func TestContext(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(5 * time.Second) // 5秒后取消
		cancelFunc()
	}()
	doSomething(ctx)
}
