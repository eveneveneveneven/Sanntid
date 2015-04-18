package main

import (
	"fmt"
	"runtime"

	"./elev"
	"./network"
	"./types"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Start program!")

	// master notification channel
	becomeMaster := make(chan bool)

	// Init of modules
	elevatorHub := elev.NewElevator()
	networkHub := network.NewNetworkHub()

	go elevator.Run(becomeMaster)
	go networkHub.Run(becomeMaster)

	select {}
}
