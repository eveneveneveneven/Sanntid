package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"./elev"
	"./network"
	"./types"
)

const (
	NUM_BACKUPS = 1
)

var child = flag.Bool("c", false, "decides if the program is a child")

func init() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func createBackupsAndListen(sigc chan os.Signal) {
	go func() {
		iter := 0
		for iter < NUM_BACKUPS {
			s := <-sigc
			iter++
			fmt.Println("\n", s)
			fmt.Println("Number of backups:", iter)
		}
		os.Exit(0)
	}()
	fmt.Println("Start backup process!")
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	for {
		process, err := os.StartProcess(dir+os.Args[0][1:], []string{os.Args[0], "-c"}, &procAttr)
		if err != nil {
			fmt.Printf("\x1b[31;1mError\x1b[0m |createBackupAndListen| [%v], continue\n", err)
		}
		_, err = process.Wait()
		if err != nil {
			fmt.Printf("\x1b[31;1mError\x1b[0m |createBackupAndListen| [%v], continue\n", err)
		}
		fmt.Println("Child process has exited!")
		time.Sleep(500 * time.Millisecond)
	}
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
		createBackupsAndListen(sigc)
	} else {
		cleanup := make(chan bool)
		go cleanupWhenExiting(cleanup, sigc)
		fmt.Println("Start program!")

		nethubToElevCh := make(chan *types.NetworkMessage, 1)
		elevToNethubCh := make(chan *types.NetworkMessage, 1)

		// Init of modules
		elevatorHub := elev.NewElevatorHub(cleanup, elevToNethubCh, nethubToElevCh)
		networkHub := network.NewNetworkHub(nethubToElevCh, elevToNethubCh)

		go elevatorHub.Run()
		go networkHub.Run()

		select {}
	}
	fmt.Println("The End!")
}
