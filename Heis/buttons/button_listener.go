package button_listener

import (
	. "../driver"
	"../lights"
)

func Get_orders(orders [][]int, current_order int) ([][]int) {
	for {
		for floor := 0; floor < 4; floor++ {
			if floor != 3 {
				if Heis_get_button(BUTTON_CALL_UP, floor) == 1 && orders[floor][BUTTON_CALL_UP] != 1{
					orders[floor][BUTTON_CALL_UP] = 1
				}
			}
			if Heis_get_button(BUTTON_COMMAND, floor) == 1 && orders[floor][BUTTON_COMMAND] != 1{
				orders[floor][BUTTON_COMMAND] = 1
			}
			if floor != 0 {
				if Heis_get_button(BUTTON_CALL_DOWN, floor) == 1 && orders[floor][BUTTON_CALL_DOWN] != 1{
					orders[floor][BUTTON_CALL_DOWN] = 1
				}
			}
		}
		lights.Set_internal_lights(orders)
		lights.Set_external_lights(orders)
		return orders
	}
}