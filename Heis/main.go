package main

import (
	"fmt"
	"./internal"
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
	int_button := make(chan int)
	ext_button := make(chan int)
	int_order := make(chan string)
	ext_order := make(chan string)
	direction := make(chan string)
	go internal.Internal(int_button, ext_button, int_order, ext_order, direction)
	neverQuit := make(chan string)
	<-neverQuit
	/*
	udp := network.NewUDPHub()
	if found, _ := udp.FindMaster(); found {
		fmt.Println("Found master!")
	} else {
		fmt.Println("Did not find master!")
	}*/
}
