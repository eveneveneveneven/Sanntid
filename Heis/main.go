package main

import (
	"fmt"
<<<<<<< HEAD
	"runtime"
=======
>>>>>>> 23d545ecd2d0f08b3f30f1b40de27a25060d8f4c

	"./elev"
	"./netstat"
	"./network"
	"./order"
	"./types"
)

func main() {
<<<<<<< HEAD
	runtime.GOMAXPROCS(runtime.NumCPU())
=======
>>>>>>> 23d545ecd2d0f08b3f30f1b40de27a25060d8f4c
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
	elevator := elev.NewElevator(orderToElev, elevOrder, elevNewElevStat)
	orderHandler := order.NewOrderHandler(netstatToOrder, orderToElev)
	networkHub := network.NewHub(becomeMaster, nethubToNetstat, netstatToNethub)
	netstatHandler := netstat.NewNetstatHandler(becomeMaster, nethubToNetstat, netstatToNethub,
		netstatToOrder, elevNewElevStat, elevOrder)

	go elevator.Run()
	go networkHub.Run()
	go orderHandler.Run()
	go netstatHandler.Run()

	select {}
}
