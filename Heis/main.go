package main

import (
	"fmt"
	"time"

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
	
	for i := 0; i < 3; i++ {
		driver.Heis_set_speed(300)
		for driver.Heis_get_floor() != 1 {}
		
		driver.Heis_set_speed(-300)
		for driver.Heis_get_floor() != 0 {}
	}
	
	driver.Heis_set_speed(0)
	select{}
	/*
	udp := network.NewUDPHub()
	if found, _ := udp.FindMaster(); found {
		fmt.Println("Found master!")
	} else {
		fmt.Println("Did not find master!")
	}*/
}