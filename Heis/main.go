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
	fmt.Println("Start program!")

	// master notification channel
	becomeMaster := make(chan bool)

	// NETSTAT <--> NETHUB
	nethubToNetstat := make(chan *types.NetworkMessage)
	netstatToNethub := make(chan *types.NetworkMessage)

	// NETSTAT --> ORDER
	netstatToOrder := make(chan *types.NetworkMessage)

	// ELEV --> NETSTAT
	elevNewElevStat := make(chan *types.ElevStat)
	elevOrder := make(chan *types.Order)

	// ORDER --> ELEV
	orderToElev := make(chan *types.Order)

	// Init of modules
	orderHandler := order.NewOrderHandler(netstatToOrder, orderToElev)
	netstatHandler := netstat.NewNetstatHandler(becomeMaster, nethubToNetstat, netstatToNethub,
		netstatToOrder, elevNewElevStat, elevOrder)
	networkHub := network.NewHub(becomeMaster, nethubToNetstat, netstatToNethub)

	go orderHandler.Run()
	go netstatHandler.Run()
	go networkHub.Run()

	select {}
}
