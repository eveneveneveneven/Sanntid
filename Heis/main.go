package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"./elev"
	"./network"
	"./types"
)

func cleanupFunc(cleanup chan bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Got a crash with err :: %+v\n", r)
		}
		close(cleanup)
		elev.CleanExit()
		os.Exit(0)
	}()
	sigCh := make(chan os.Signal, 1)
	startCleanup := make(chan bool)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		for _ = range sigCh {
			fmt.Println("\nReceived an interrupt, stopping program...\n")
			startCleanup <- true
		}
	}()
	<-startCleanup
	panic(nil)
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cleanup := make(chan bool)
	go cleanupFunc(cleanup)

	fmt.Println("Start program!")

	nethubToElevCh := make(chan *types.NetworkMessage)
	elevToNethubCh := make(chan *types.NetworkMessage)

	// Init of modules
	elevatorHub := elev.NewElevatorHub(cleanup, elevToNethubCh, nethubToElevCh)
	networkHub := network.NewNetworkHub(nethubToElevCh, elevToNethubCh)

	go elevatorHub.Run()
	go networkHub.Run()

	select {}
}
