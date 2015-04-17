package elev

import (
	"../driver"
	"../types"
	"fmt"
	"time"
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
}

func SetOrderLight(Order *types.Order) {
	if (Order.ButtonPress >= 0 && Order.ButtonPress <= 2) &&
		(Order.Floor >= 0 && Order.Floor <= 3) {
		if Order.Completed {
			driver.Heis_set_button_lamp(Order.ButtonPress, Order.Floor, 0)
		} else {
			driver.Heis_set_button_lamp(Order.ButtonPress, Order.Floor, 1)
		}
	} else {
		fmt.Printf(`\t\x1b[31;1mError\x1b[0m |SetOrderLight| [Order recieved is not valid,
			 got the following %+v], exit program\n`, Order)
	}
}

func floorIndicator() {
	var floor int
	for {
		floor = driver.Heis_get_floor()
		if floor != -1 {
			driver.Heis_set_floor_indicator(floor)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
