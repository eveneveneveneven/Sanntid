package lights

import(
	. "../driver"
)

func Set_external_lights(orders [][]int) {
	for i := 0; i < 4; i++ {
		if i != 3 {
			Heis_set_button_lamp(BUTTON_CALL_UP, i, orders[i][0])
		}
		if i != 0 {
			Heis_set_button_lamp(BUTTON_CALL_DOWN, i, orders[i][1])
		}
	}
}
func Set_internal_lights(orders [][]int) {
	for i := 0; i < 4; i++ {
		Heis_set_button_lamp(BUTTON_COMMAND, i, orders[i][2])
	}
}