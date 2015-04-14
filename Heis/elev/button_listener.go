package elev

import (
	. "../driver"
	"../lights"
)

type Button struct {
	Button_type Elev_button_type_t
	Floor       int
}

func reset_orders(orders [][]int) [][]int {
	for i := 0; i < 5; i++ {
		orders[i] = make([]int, 4)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			orders[i][j] = 0
		}
	}
	return orders
}

func Get_orders() {
	orders := make([][]int, 5)
	for i := 0; i < 5; i++ {
		orders[i] = make([]int, 4)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 4; j++ {
			orders[i][j] = 0
		}
	}

	for {
		for floor := 0; floor < 4; floor++ {
			if floor != 3 {
				if Heis_get_button(BUTTON_CALL_UP, floor) == 1 && orders[floor][BUTTON_CALL_UP] != 1 {
					reset_orders(orders)
					orders[floor][BUTTON_CALL_UP] = 1
					/*b := &Button{
						Button_type: BUTTON_CALL_UP,
						Floor: floor,
					}
					*/
					break
				}
			}
			if Heis_get_button(BUTTON_COMMAND, floor) == 1 && orders[floor][BUTTON_COMMAND] != 1 {
				reset_orders(orders)
				orders[floor][BUTTON_COMMAND] = 1
				/*b := &Button{
						Button_type: 	BUTTON_COMMAND,
						Floor: 			floor,
				}*/
				break
			}
			if floor != 0 {
				if Heis_get_button(BUTTON_CALL_DOWN, floor) == 1 && orders[floor][BUTTON_CALL_DOWN] != 1 {
					reset_orders(orders)
					orders[floor][BUTTON_CALL_DOWN] = 1
					/*b := &Button{
						Button_type: BUTTON_CALL_DOWN,
						Floor: floor,
					}*/
					break
				}
			}
		}
		lights.Set_internal_lights(orders)
		//ch <- 0
	}
}
