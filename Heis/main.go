package main

import (
	"fmt"

	//"./driver"
	"./network"
	"./order_handler"
	"./types"
)

func main() {
	fmt.Println("Start main!")

	c := make(chan *types.NetworkMessage)

	orderHandler := order_handler.NewOrderHandler(c)
	networkHub := network.NewHub(c)

	go orderHandler.Run()
	go networkHub.Run()

	select {}
}
