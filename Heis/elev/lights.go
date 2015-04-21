package elev

import (
	"fmt"

	"../driver"
	"../types"
)

func clearAllLights() {
	for f := 0; f < 4; f++ {
		for b := 0; b < 3; b++ {
			if (f == 0 && b == 1) || (f == 3 && b == 0) {
				continue
			}
			driver.Heis_set_button_lamp(b, f, 0)
		}
	}
	driver.Heis_set_stop_light(0)
	driver.Heis_set_door_open_lamp(0)
}

func setActiveLights(netstat *types.NetworkMessage) {
	for order, completed := range netstat.Orders {
		if completed {
			delete(netstat.Orders, order)
		}
		setOrderLight(&order, completed)
	}
	etgs := []bool{true, true, true, true}
	for _, etg := range netstat.Statuses[netstat.Id].InternalOrders {
		if etg != -1 {
			setOrderLight(&types.Order{types.BUTTON_INTERNAL, etg}, false)
			etgs[etg] = false
		}
	}
	for etg, v := range etgs {
		if v {
			setOrderLight(&types.Order{types.BUTTON_INTERNAL, etg}, true)
		}
	}
}

func setOrderLight(order *types.Order, completed bool) {
	f := order.Floor
	b := order.ButtonPress
	if (b >= 0 && b <= 2) && (f >= 0 && f <= 3) &&
		!((f == 0 && b == 1) || (f == 3 && b == 0)) {

		if completed {
			driver.Heis_set_button_lamp(b, f, 0)
		} else {
			driver.Heis_set_button_lamp(b, f, 1)
		}
	} else {
		fmt.Printf("\t\x1b[31;1mError\x1b[0m |SetOrderLight| [Order recieved is not valid, got the following %+v], exit program\n", order)
	}
}

func floorIndicator() {
	var floor int
	for {
		floor = driver.Heis_get_floor()
		if floor != -1 {
			driver.Heis_set_floor_indicator(floor)
		}
	}
}

func setStopLight(value int) {
	driver.Heis_set_stop_light(1)
}

func setDoorLight(value int) {
	driver.Heis_set_door_open_lamp(value)
}
