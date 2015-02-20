package main

import (
	"fmt"

	"./driver"
	//"./network"
)


func main() {
	fmt.Println("Start main!")
	if driver.Heis_init() {
		fmt.Println("init success")
	} else {
		fmt.Println("init failed")
	}
	
	/*
	udp := network.NewUDPHub()
	if found, _ := udp.FindMaster(); found {
		fmt.Println("Found master!")
	} else {
		fmt.Println("Did not find master!")
	}*/
}
