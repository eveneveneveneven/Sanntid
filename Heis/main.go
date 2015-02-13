package main

import (
	"fmt"

	"./driver"
	"./network"
)


func main() {
	fmt.Println("Start main!")
	driver.Io_init()
	udp := network.NewUDPHub()
	if found, _ := udp.FindMaster(); found {
		fmt.Println("Found master!")
	} else {
		fmt.Println("Did not find master!")
	}
}