package main

import (
	"fmt"
	"time"
)

var i int

func thread1(ch chan<- int) {
	for j := 0; j < 100000; j++ {
		i++
	}
	ch <- 0
}

func thread2(ch chan<- int) {
	for j := 0; j < 100000; j++ {
		i--
	}
	ch <- 0
}

func main() {
	t0 := time.Now()
	ch1 := make(chan int)
	ch2 := make(chan int)
	go thread1(ch1)
	go thread2(ch2)
	<-ch1
	<-ch2
	t1 := time.Now()
	fmt.Printf("The call took %v to run, with value i = %v.\n", t1.Sub(t0), i)
}
