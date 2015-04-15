package main

import (
	//"fmt"
	//"./internal"
	//"./driver"
	//"runtime"
	."./cost"
	."./types"
	//"./network"
)


func main() {
	/*runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Start main!")
	if driver.Heis_init() {
		fmt.Println("init success")
	} else {
		fmt.Println("init failed")
	}
	go internal.Internal()
	select{}*/
	elev1 := ElevStat{
		Dir: STOP,
		Floor: 1,
		InternalOrders: nil,
	}
	elev2 := ElevStat{
		Dir: UP,
		Floor: 1,
		InternalOrders: nil,
	}
	elev3 := ElevStat{
		Dir: DOWN,
		Floor: 3,
		InternalOrders: nil,
	}
	elev4 := ElevStat{
		Dir: UP,
		Floor: 0,
		InternalOrders: nil,
	}
	orders := make(map[Order]struct{})
	o1 := Order{
		ButtonPress: BUTTON_CALL_UP,
		Floor: 2,
		Completed: false,
	}
	o2 := Order{
		ButtonPress: BUTTON_CALL_DOWN,
		Floor: 2,
		Completed: false,
	}
	o3 := Order{
		ButtonPress: BUTTON_CALL_DOWN,
		Floor: 3,
		Completed: false,
	}
	o4 := Order{
		ButtonPress: BUTTON_CALL_UP,
		Floor: 0,
		Completed: false,
	}
	orders[o1] = struct{}{}
	orders[o2] = struct{}{}
	orders[o3] = struct{}{}
	orders[o4] = struct{}{}

	nm := &NetworkMessage{
		Id: 0,
		Statuses: []ElevStat{elev1, elev2, elev3, elev4},
		Orders: orders,
	}
	Cost_function(nm)


	
}
