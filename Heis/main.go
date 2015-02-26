package main

import (
	"fmt"

	//"./driver"
	"./network"
)


func main() {
	fmt.Println("Start main!")

	/*
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
	*/

	hub := network.NewHub()
	go hub.Run()

	select{}
}