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

	statusRecieve := make(chan *types.NetworkMessage)
	statusSend := make(chan *types.NetworkMessage)
	
	oh := order_handler.NewOrderHandler(statusRecieve, statusSend)
	hub := network.NewHub(statusRecieve, statusSend)

	go oh.Run()
	go hub.Run()

	select {}
}
