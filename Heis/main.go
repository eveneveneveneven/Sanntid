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
	go hub.Run()

	select {}
	fmt.Println("Ending program")
}