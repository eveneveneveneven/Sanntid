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

	statusConnector := make(chan *types.NetworkMessage)
	oh := order_handler.NewOrderHandler(statusConnector)
	hub := network.NewHub(statusConnector)

	go oh.Run()
	go hub.Run()

	select {}
}
