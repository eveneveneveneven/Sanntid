package main

import (
	"fmt"

	"./driver"
	"./network"
)


func main() {
	fmt.Println("Start main!")
	driver.Io_init()
	
	hub := network.NewHub()
	stop := make(chan bool)
	becameMaster, _ := hub.ResolveMasterNetwork(stop)
	if becameMaster {
		fmt.Println("I am Master!")
	} else {
		fmt.Println("I am a slave...")
	}
	select {}
	fmt.Println("Ending program")
}