package elev

import (
	//"../driver"
	"../types"
)

func buttonListener(orderLn chan *types.Order) {
	pressed := make([]bool, 4)
	for i := 0; i < 4; i++ {
		pressed[i] = make([]bool, 3)
		for j := 0; j < 3; j++ {
			pressed[i][j] = false
		}
	}

	for {
		for f := 0; f < 4; f++ {
			for b := 0; b < 3; b++ {
				if (f == 0 && b == 1) || (f == 3 && b == 0) {
					continue
				}
				if driver.Heis_get_button(b, f) {
					if !pressed[f][b] {
						pressed[f][b] = true
						order := &types.Order{
							ButtonPress: b,
							Floor:       f,
							Completed:   false,
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
