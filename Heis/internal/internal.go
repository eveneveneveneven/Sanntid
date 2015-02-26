package internal
import (
	. ".././driver"
	. "fmt"
	"time"
)


func Init_orders() ([][]int){
	orders := make([][]int, 5)
	for i := 0; i < 5; i++ {
		orders[i] = make([]int, 4)
	}
	for i :=0; i<5; i++ {
		for j := 0; j<4; j++{
			orders[i][j]=0
		}
	}
	Print("Orders:")
	Println("")
	Print_orders(orders)
	return orders
}

func Print_orders(orders [][]int){
	for i :=0; i<4; i++ {
		for j := 0; j<3; j++{
			Print(orders[i][j])
			Print("     ")
		}
		println("")
	}
	println("")
}
//Sends elevator to specified floor
func remove_from_queue(queue []int) ([]int){

for i:=0; i<3; i++{
		queue[i] = queue[i+1]
	}
	queue[3] = -1
	return queue
}
func remove_from_orders(orders [][]int, current_order int) ([][]int){
	for i:=0; i<3; i++{
		orders[current_order][i] = 0
	}
	return orders
}
func Send_to_floor(queue []int, orders[][]int) ([]int, [][]int) {
	for i:=0; i<4; i++ {
		current_order := queue[i]
		current_floor := Heis_get_floor()
		if current_order==-1{
			Heis_set_speed(0)
		}
		if Heis_get_floor() != -1 && current_order !=-1{
			Print("Current floor: ")
			Println(current_floor + 1)			
			Print("Going to: Floor ")
			Println(current_order + 1)
		}
		if current_order==current_floor {
			Heis_set_speed(0)
			set_lights(current_floor)
		}
		if current_order>current_floor && current_order != -1{
			Heis_set_speed(150)
			for {
				if Heis_get_floor() == current_order {
					Heis_set_speed(0)
					set_lights(current_order)
					Print("Arrived at floor: ")
					Println(current_order + 1)
					queue = remove_from_queue(queue)
					orders = remove_from_orders(orders, current_order)
					return queue, orders
				}
			}
		}
		if current_order<current_floor && current_order != -1{
			Heis_set_speed(-150)		
			for {
				if Heis_get_floor() == current_order {
				Heis_set_speed(0)
				set_lights(current_order)
				Print("Arrived at floor: ")
				Println(current_order + 1)
				queue = remove_from_queue(queue)
				orders = remove_from_orders(orders, current_order)
				return queue, orders
				}
			}
		}
	}
	return queue, orders
}
func get_queue(orders [][]int) ([]int){
	queue := make([]int, 4)
	for i :=0; i<4; i++{
		queue[i]=-1
	}
	k:=0
	for i :=0; i<4; i++ {
		
		if orders[i][0]==1 || orders[i][1]==1 || orders[i][2]==1{
			queue[k]=i
			Println("Order received: Floor ")
			Println(i + 1)
			k++
		}
	
	}
	return queue
}
//Handles external button presses
func get_orders(orders [][]int) {
	for{
		for floor:=0; floor<4; floor++ {
		

			if floor != 3 {
				if Heis_get_button(BUTTON_CALL_UP, floor)==1{
					//Println("External call up button nr: " + Itoa(floor) + " has been pressed!")
					Heis_set_button_lamp(BUTTON_CALL_UP, floor, 1)
					orders[floor][BUTTON_CALL_UP]=1
			
				}
			}
			if Heis_get_button(BUTTON_COMMAND, floor) == 1 {
				//Println("Internal button nr: " + Itoa(floor) + " has been pressed!")
				Heis_set_button_lamp(BUTTON_COMMAND, floor, 1)
				orders[floor][BUTTON_COMMAND]=1;
			}
			if floor != 0 {
				if Heis_get_button(BUTTON_CALL_DOWN, floor) == 1 {
					//Println("External call down button nr: " + Itoa(floor) + " has been pressed!")
					Heis_set_button_lamp(BUTTON_CALL_DOWN, floor, 1)
					orders[floor][BUTTON_CALL_DOWN]=1
				
				}
			}
		}
	}
}

//Checks which floor the elevator is on and sets the floor-light
func Floor_indicator() {
	Println("executing floor indicator!")
	var floor int
	for {
		floor = Heis_get_floor()
		//Println(floor)
		if floor != -1 {
			Heis_set_floor_indicator(floor)
			//time.Sleep(50 * time.Millisecond)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

/*func Clear_orders(orders [][]int){
		
	for{
		current_floor := Heis_get_floor()
		if current_floor!=-1{
			for k:=0; k<3; k++{
				orders[current_floor][k]=0
			}
		}
	time.Sleep(1000*time.Millisecond)
	}
*/
func To_nearest_floor() {
	for {
		
		Heis_set_speed(-150)
		if Heis_get_floor() == 0 {
			Heis_set_speed(0)
			return
		}
		if Heis_get_floor() != -1 {
			Heis_set_speed(0)
			return
		}
	}
}
/*func get_stop(){
	for{
		if Heis_stop(){
			Heis_set_stop_lamp(1)
			stop_all()
		}
		time.Sleep(time.Millisecond*10)
	}
}

func stop_all(){
	Heis_set_speed(0)
	time.Sleep(time.Millisecond*1000)
	Heis_set_stop_lamp(0)
	To_nearest_floor()
}
*/

func set_lights(current_floor int){
	if current_floor == 0 {
		Heis_set_button_lamp(BUTTON_CALL_UP, current_floor, 0)
		Heis_set_button_lamp(BUTTON_COMMAND, current_floor, 0)
	}
	if current_floor == 3 {
		Heis_set_button_lamp(BUTTON_CALL_DOWN, current_floor, 0)
		Heis_set_button_lamp(BUTTON_COMMAND, current_floor, 0)
	}
	if current_floor == 2 || current_floor == 1{
		Heis_set_button_lamp(BUTTON_CALL_UP, current_floor, 0)
		Heis_set_button_lamp(BUTTON_CALL_DOWN, current_floor, 0)
		Heis_set_button_lamp(BUTTON_COMMAND, current_floor, 0)
	}
}
func Internal() {
// Initialize
	orders := Init_orders()
	queue := get_queue(orders)
	Heis_init()
	Heis_set_speed(0)
	To_nearest_floor()
	Heis_set_stop_lamp(0)
	//go get_stop()
	go Floor_indicator()
	go get_orders(orders)
	//go Interpret_orders(orders)
	for{
		//orders = Ext_order(orders)
		//Print_orders(orders)
		//time.Sleep(20 * time.Millisecond)
		queue = get_queue(orders)
		queue, orders = Send_to_floor(queue, orders)
	}

	neverQuit := make(chan string)
	<-neverQuit
}
