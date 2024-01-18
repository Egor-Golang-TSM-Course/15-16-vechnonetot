package main

import (
	"context"
	"fmt"
	"time"
)

func worked(ctx context.Context, data <-chan int, next chan<- int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("WORKED STOPPED")
			return

		case val := <-data:
			val += 1
			if val > 5 {
				fmt.Println("ERROR IN DATA PROCESSING")
				return
			}
			next <- val
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	data1 := make(chan int)
	data2 := make(chan int)

	go worked(ctx, data1, data2)
	go worked(ctx, data2, data1)

	data1 <- 1

	time.Sleep(time.Second * 1)
	cancel()

	time.Sleep(time.Second * 1)
}
