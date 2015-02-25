package main
import (
. "./driver"
. "fmt"
. "strconv"
"time"
)
//Waits for internal and external inputs from elevator control
func Wait_for_input(int_button, ext_button chan int, int_order, ext_order, last_order, direction chan string) {
	var floor int
	for {
		select {
		case floor = <-int_button:
			int_order <- Itoa(floor)
		case floor = <-ext_button:
			dir := <-direction
			ext_order <- Itoa(floor) + ":" + dir
		case temp := <-last_order:
			_ = temp
			time.Sleep(50 * time.Millisecond)
		}
	}
}
//Sends elevator to specified floor
func Send_to_floor(floor int, button string) {
	current_floor := Heis_get_floor()
	Heis_set_door_open_lamp(0)
	Heis_set_stop_lamp(0)
	if current_floor < floor {
		Println("Going up")
		for {
			Heis_set_speed(150)
			if Heis_get_floor() == floor {
				Println("I am now at floor: " + Itoa(Heis_get_floor()))
				Heis_set_stop_lamp(1)
				Heis_set_door_open_lamp(1)
				time.Sleep(25 * time.Millisecond)
				Heis_set_speed(0)
				if button == "int" {
					Heis_set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					if button == "up" {
						Heis_set_button_lamp(BUTTON_CALL_UP, floor, 0)
					} else {
						Heis_set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
					}
				}
				return
			}
			time.Sleep(25 * time.Millisecond)
		}
	} else {
		Println("Going down")
		for {
			Heis_set_speed(-150)
			if Heis_get_floor() == floor {
				Println("I am now at floor: " + Itoa(Heis_get_floor()))
				Heis_set_stop_lamp(1)
				Heis_set_door_open_lamp(1)
				time.Sleep(25 * time.Millisecond)
				Heis_set_speed(0)
				if button == "int" {
					Heis_set_button_lamp(BUTTON_COMMAND, floor, 0)
				} else {
					if button == "up" {
						Heis_set_button_lamp(BUTTON_CALL_UP, floor, 0)
					} else {
						Heis_set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
					}
				}
				return
			}
			time.Sleep(25 * time.Millisecond)
		}
	}
}
//Keyboard terminal input (For testing)
func KeyboardInput(ch chan int) {
	var a int
	for {
		Scan(&a)
		ch <- a
	}
}
//Handles external button presses
func Ext_order(int_button chan int, direction chan string) {
	i := 0
	for {
		if i < 3 {
			if Heis_get_button(BUTTON_CALL_UP, i) == 1 {
				Println("External call up button nr: " + Itoa(i) + " has been pressed!")
				Heis_set_button_lamp(BUTTON_CALL_UP, i, 1)
				int_button <- i
				direction <- "up"
				time.Sleep(300 * time.Millisecond)
			}
		}
		if i > 0 {
			if Heis_get_button(BUTTON_CALL_DOWN, i) == 1 {
				Println("External call down button nr: " + Itoa(i) + " has been pressed!")
				Heis_set_button_lamp(BUTTON_CALL_DOWN, i, 1)
				int_button <- i
				direction <- "down"
				time.Sleep(300 * time.Millisecond)
			}
		}
		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)
	}
}
//Handles internal button presses
func Int_order(int_button chan int) {
	i := 0
	for {
		if Heis_get_button(BUTTON_COMMAND, i) == 1 {
			Println("Internal button nr: " + Itoa(i) + " has been pressed!")
			int_button <- i
			Heis_set_button_lamp(BUTTON_COMMAND, i, 1)
			time.Sleep(300 * time.Millisecond)
		}
		i++
		i = i % 4
		time.Sleep(25 * time.Millisecond)
	}
}
//Checks which floor the elevator is on and sets the floor-light
func Floor_indicator(last_order chan string) {
	Println("executing floor indicator!")
//_ = last_order
	var floor int
	for {
		floor = Heis_get_floor()
		if floor != -1 {
			Heis_set_floor_indicator(floor)
			last_order <- Itoa(floor)
			time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(25 * time.Millisecond)
	}
}
func To_nearest_floor() {
	for {
		Heis_set_speed(150)
		if Heis_get_floor() != -1 {
			time.Sleep(25 * time.Millisecond)
			Heis_set_speed(0)
		}
	}
}
func Internal(int_button, ext_button chan int, int_order, ext_order, last_order, direction chan string) {
// Initialize
	Heis_init()
	Heis_set_speed(0)
	Heis_set_stop_lamp(1)
	go Floor_indicator(last_order)
	go Int_order(int_button)
	go Ext_order(ext_button, direction)
	go Wait_for_input(int_button, ext_button, int_order, 		ext_order, last_order, direction)
	neverQuit := make(chan string)
	<-neverQuit
}

