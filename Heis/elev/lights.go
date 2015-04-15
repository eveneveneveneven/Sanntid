package elev

import (
	"../driver"
	"../types"
	"time"
)

func setOrderLight(Order *types.Order) {
	if (Order.ButtonPress == BUTTON_CALL_UP || Order.ButtonPress == BUTTON_CALL_UP ||
		Order.ButtonPress == BUTTON_CALL_UP) && (Order.Floor >= 0 && Order.Floor <= 3) {
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
			Heis_set_floor_indicator(floor)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
