package main

import (
	"fmt"

	//"./driver"
	"./netstat"
	"./network"
	"./order"
	"./types"
)

func main() {
	fmt.Println("Start main!")

	c := make(chan *types.NetworkMessage)

	becomeMaster := make(chan bool)

	netToNetstat := make(chan *types.NetworkMessage)
	netstatToNet := make(chan *types.NetworkMessage)

	netstatToOrder := make(chan *types.NetworkMessage)
	orderToNetstat := make(chan *types.NetworkMessage)

	orderHandler := order.NewOrderHandler(becomeMaster, netstatToOrder, orderToNetstat, nil, nil)
	netStatHandler := netstat.NewNetStatHandler(becomeMaster, netToNetstat, netstatToNet, netstatToOrder, orderToNetstat)
	networkHub := network.NewHub(becomeMaster, netToNetstat, netstatToNet)

	go orderHandler.Run()
	go netStatHandler.Run()
	go networkHub.Run()

	select {}
}
