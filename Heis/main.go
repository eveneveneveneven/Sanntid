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

	nethubToElevCh := make(chan *types.NetworkMessage)
	elevToNethubCh := make(chan *types.NetworkMessage)

	// Init of modules
	elevatorHub := elev.NewElevatorHub(becomeMaster, elevToNethubCh, nethubToElevCh)
	networkHub := network.NewNetworkHub(becomeMaster, nethubToElevCh, elevToNethubCh)

	go elevatorHub.Run()
	go networkHub.Run()

	select {}
}
