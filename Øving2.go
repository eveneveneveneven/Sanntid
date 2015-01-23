package main

import(
"fmt"
"runtime"
"time"
)

var sum int = 0
var buffer = make(chan int,1)
var adder_done = make(chan int,1)
var subtracter_done = make(chan int,1)

func adder(){
	for i := 0; i <= 1000000; i++ {
		<- buffer
		sum++
		buffer <- 1	
	}
	adder_done <- 1
}

func subtracter(){
	for j := 0; j <= 1000000; j++{
		<- buffer
		sum--
		buffer <- 1
	}
	subtracter_done <- 1
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	t0 := time.Now()
	buffer <- 1
	go adder()
	go subtracter()
	<- adder_done
	<- subtracter_done
	t1 := time.Now()
	fmt.Printf("The call took %v to run, with value i = %v.\n", t1.Sub(t0), sum)
	fmt.Println("Yolo")
}
