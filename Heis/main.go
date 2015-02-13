package main

import (
	"fmt"
	"time"

	"./driver"
	"./network"
)


func main() {
	fmt.Println("Start main!")
	driver.Io_init()
	udp := network.NewUDPHub()
	if found, _ := udp.FindMaster(); found {
		fmt.Println("Found master do nothing!")
	} else {
		fmt.Println("Did not find master, becoming master!")
		stop := make(chan bool)
		go udp.BroadcastMaster(stop)
		time.Sleep(5 * time.Second)
		stop <- true

	}

	fmt.Println("Ending program")
}