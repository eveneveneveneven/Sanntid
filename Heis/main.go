package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"./elev"
	"./network"
	"./types"
)

var child = flag.Bool("c", false, "decides if the program is a child")
var noBackup = flag.Bool("nb", false, "start the program without backup")

func cleanupFunc(cleanup chan bool) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt)

	cmd := exec.Command("gnome-terminal", "-e", "./main -c")
	cmd.Output()

	go func() {
		<-sigc
		cleanup <- true
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()
}

func backupProcess() {
	fmt.Println("Starting Backup Process")
	startBackup := make(chan bool)
	go network.StartBackupListener(startBackup)
	<-startBackup
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	cleanup := make(chan bool)
	if !*noBackup {
		if *child {
			backupProcess()
		} else {
			go cleanupFunc(cleanup)
		}
	}
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
