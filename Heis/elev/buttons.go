package elev

import (
	"../driver"
	"../types"
)

func buttonListener(orderLn chan types.Order) {
	pressed := make([][]bool, M_FLOORS)
	for i := 0; i < M_FLOORS; i++ {
		pressed[i] = make([]bool, 3)
		for j := 0; j < 3; j++ {
			pressed[i][j] = false
		}
	}

	for {
		for f := 0; f < M_FLOORS; f++ {
			for b := 0; b < 3; b++ {
				if (f == 0 && b == 1) || (f == M_FLOORS-1 && b == 0) {
					continue
				}
				if driver.Heis_get_button(b, f) {
					if !pressed[f][b] {
						pressed[f][b] = true
						order := types.Order{
							ButtonPress: b,
							Floor:       f,
						}
						orderLn <- order
					}
				} else {
					if pressed[f][b] {
						pressed[f][b] = false
					}
				}
			}
		}
	}
}
