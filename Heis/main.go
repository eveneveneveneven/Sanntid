package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"./backup"
	"./elev"
	"./network"
	"./types"
)

var child = flag.Bool("c", false, "Notifies if the program is a child process")

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func cleanupWhenExiting(cleanup chan bool, sigc chan os.Signal) {
	s := <-sigc
	fmt.Println("\n", s)
	cleanup <- true
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)
}

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt)
	if !*child {
		backup.CreateBackupsAndListen(sigc)
	} else {
		cleanup := make(chan bool)
		go cleanupWhenExiting(cleanup, sigc)
		fmt.Println("Start program!")

		resetCh := make(chan bool)
		nethubToElevCh := make(chan *types.NetworkMessage, 1)
		elevToNethubCh := make(chan *types.NetworkMessage, 1)

		// Init of modules
		elevatorHub := elev.NewElevatorHub(cleanup, resetCh, elevToNethubCh, nethubToElevCh)
		networkHub := network.NewNetworkHub(resetCh, nethubToElevCh, elevToNethubCh)

		go elevatorHub.Run()
		go networkHub.Run()

		select {}
	}
	fmt.Println("The End!")
}
