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
	clearAllLights()
	for _, etg := range netstat.Statuses[netstat.Id].InternalOrders {
		if etg == -1 {
			break
		}
		setOrderLight(&types.Order{
			ButtonPress: types.BUTTON_INTERNAL,
			Floor:       etg,
			Completed:   false,
		})
	}
	for order := range netstat.Orders {
		setOrderLight(&order)
	}
}

func setOrderLight(order *types.Order) {
	if (order.ButtonPress >= 0 && order.ButtonPress <= 2) &&
		(order.Floor >= 0 && order.Floor <= 3) {
		if order.Completed {
			driver.Heis_set_button_lamp(order.ButtonPress, order.Floor, 0)
		} else {
			driver.Heis_set_button_lamp(order.ButtonPress, order.Floor, 1)
		}
	} else {
		fmt.Printf(`\t\x1b[31;1mError\x1b[0m |SetOrderLight| [Order recieved is not valid,
			 got the following %+v], exit program\n`, order)
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
