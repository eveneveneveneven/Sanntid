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
	netToNetStat := make(chan *types.NetworkMessage)
	netStatToNet := make(chan *types.NetworkMessage)

	orderHandler := order.NewOrderHandler(c)
	netStatHandler := netstat.NewNetStatHandler(becomeMaster, netToNetStat, netStatToNet)
	networkHub := network.NewHub(becomeMaster, netToNetStat, netStatToNet)

	go orderHandler.Run()
	go netStatHandler.Run()
	go networkHub.Run()

	select {}
}
