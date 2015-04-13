package Send_to_floor
import (
	. "../driver"
	"time"
)

var speed int

func open_doors() {
	Heis_set_door_open_lamp(1)
	time.Sleep(1000 * time.Millisecond)
	Heis_set_door_open_lamp(0)
	return
}

func Send_to_floor(ordered_floor int){
outerloop:
	for i := 0; i < 4; i++ {
		speed = 150
		current_order := ordered_floor
		current_floor := Heis_get_floor()
		if current_order != -1 && current_order != current_floor{
			//Println("Order received: Floor ", current_order+1)
		}
		/*if current_floor == -1 {
			stop_all()
		}
		
		Print("Current floor: ")
		Println(current_floor + 1)
		Print("Going to: Floor ")
		Println(current_order + 1)
		Println(" ")
		*/

		if current_order == current_floor && current_order != -1 {
			open_doors()
		}

		if current_order > current_floor {
			Heis_set_speed(speed)
			for {
				if Heis_get_floor() == current_order {
					Heis_set_speed(0)
					//Print("Arrived at floor: ")
					//Println(current_order + 1)
					//Println(" ")
					open_doors()
					break outerloop
				}
			}
		}
		if current_order < current_floor && current_order != -1 {
			Heis_set_speed(-speed)
			for {
				if Heis_get_floor() == current_order {
					Heis_set_speed(0)
					//Print("Arrived at floor: ")
					//Println(current_order + 1)
					//Println(" ")
					open_doors()
					break outerloop
				}
			}
		}
	}
}