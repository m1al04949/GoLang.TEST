package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const N = 20

var mu sync.Mutex

func main() {

	fn := func(x int) int {
		time.Sleep(time.Duration(rand.Int31n(N)) * time.Second)
		return x * 2
	}
	in1 := make(chan int, N)
	in2 := make(chan int, N)
	out := make(chan int, N)

	start := time.Now()
	merge2Channels(fn, in1, in2, out, N+1)
	for i := 0; i < N+1; i++ {
		in1 <- i
		in2 <- i
	}

	orderFail := false
	EvenFail := false
	for i, prev := 0, 0; i < N; i++ {
		c := <-out
		if c%2 != 0 {
			EvenFail = true
		}
		if prev >= c && i != 0 {
			orderFail = true
		}
		prev = c
		fmt.Println(c)
	}
	if orderFail {
		fmt.Println("порядок нарушен")
	}
	if EvenFail {
		fmt.Println("Есть не четные")
	}
	duration := time.Since(start)
	if duration.Seconds() > N {
		fmt.Println("Время превышено")
	}
	fmt.Println("Время выполнения: ", duration)
}

func merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go func(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {

		res1 := make([]*int, n)
		res2 := make([]*int, n)

		input := func(input <-chan int, results []*int) {
			for i := 0; i < n; i++ {
				x := <-input
				go func(i int, x int) {
					res := fn(x)
					results[i] = &res
				}(i, x)
			}
		}

		go input(in1, res1)
		go input(in2, res2)

		go func() {
			i := 0
			mu.Lock()
			for true {
				if res1[i] != nil && res2[i] != nil {
					res := *res1[i] + *res2[i]
					out <- res
					if i++; i == n {
						mu.Unlock()
						return
					}
				}
			}
		}()
	}(fn, in1, in2, out, n)
}
