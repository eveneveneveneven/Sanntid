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
var noBackup = flag.Bool("nb", true, "start the program without backup")

func cleanupFunc(cleanup, createBackup chan bool) {
	<-createBackup
	fmt.Println("Starting backup process")
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

func backupProcess(startBackup chan bool) {
	fmt.Println("Starting Backup listener")
	go network.StartBackupListener(startBackup)
	<-startBackup
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	cleanUp := make(chan bool)
	startBackup := make(chan bool)
	createBackup := make(chan bool)

	if !*noBackup {
		if *child {
			backupProcess(startBackup)
		}
		go cleanupFunc(cleanUp, createBackup)
	}
	fmt.Println("Start program!")

	nethubToElevCh := make(chan *types.NetworkMessage)
	elevToNethubCh := make(chan *types.NetworkMessage)

	// Init of modules
	elevatorHub := elev.NewElevatorHub(cleanUp, elevToNethubCh, nethubToElevCh)
	networkHub := network.NewNetworkHub(*noBackup, startBackup, createBackup,
		nethubToElevCh, elevToNethubCh)

	go elevatorHub.Run()
	go networkHub.Run()

	select {}
}
