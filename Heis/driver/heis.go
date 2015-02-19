
package heis
import(
"driver"
"time"
"math"
"fmt"
)
func HeisInitTest()(int, int, int){
	direction := 0;
	for driver.GetFloor() == -1 {
		driver.SetSpeed(-300)
	}
	driver.SetSpeed(0)
	current_floor := driver.GetFloor()
	destination := -1
	return direction, current_floor, destination
}
func HeisTest(order_list [8]int, command_list [8]int, cost [8]int)([8]int,[8]int,[8]int){
	direction, current_floor, destination := HeisInitTest()
	destination = getDestination(direction, current_floor, order_list, command_list)
// decides direction required to reach destination from current floor.
	direction = getDirection(destination, current_floor)
	fmt.Println("Destination: ", destination)
	if (direction != 0){
		fmt.Println("Current floor: ", current_floor, "Destination: ", destination, "Direction: ", direction, "\nCommand list: ", command_list)
	}
// order_list, command_list = RemoveOrders(order_list, command_list, direction)
	driver.SetSpeed(direction*300)
	for(destination != -1){
		if(driver.GetFloor() != -1){
			current_floor = driver.GetFloor()
		}
		if(direction==-1 && (order_list[2*current_floor]==1) || (direction==1 && order_list[2*current_floor+1]==1) || command_list[current_floor] == 1 || (destination == current_floor)){
			fmt.Println("dørene åpnes i etasje ", current_floor)
			driver.SetSpeed(0)
			driver.SetDoorLamp(1)
			order_list, command_list = removeOrders(order_list, command_list, direction, destination)
			time.Sleep(3*time.Second)
			driver.SetDoorLamp(0)
			fmt.Println("dørene lukkes")
			fmt.Println(order_list)
			if current_floor == destination{
				destination = -1
			}
			break
		}
		time.Sleep(time.Millisecond*5)
	}
// direction
// 2) else if there are orders in order list. Complete them until
	return order_list, command_list, cost
//}
}
func getDirection(destination int, current_floor int)(int){
	direction := 0
	if (destination == -1){
		return direction
	}else if(destination > current_floor){
		direction = 1
	}else if(destination < current_floor){
		direction = -1
	}
	return direction
}
func getDestination(direction int, current_floor int, order_list [8]int, command_list [8]int)(int){
	var i int
	if(direction == 1){
		i = 3
		for(i >= current_floor){
			if (order_list[i*2+1] == 1 || command_list[i] == 1){
				return i
			}
			i -= 1
		}
		return -1
	}else if (direction == -1){
		i = 0
		for(i <= current_floor){
			if (order_list[i*2] == 1 || command_list[i] == 1){
				return i
			}
			i += 1
		}
		return -1
//hvis ikke behold det som destination
	}else{
		i = 0
		for(i < 4){
			if (order_list[i*2] == 1 || order_list[i*2+1] == 1 || command_list[i] == 1){
				return i
			}
			i += 1
		}
		return -1
//sjekk, command lista, så order lista(?) sett første som finnes til destination.
//hvis ikke, sett destination til eller 0 eller noe.
	}
}
func costFunction(current_floor int,direction int, destination int)([8]int){
	i := 0
	var cost [8]int
	for i<8{
		if (direction == 0){
cost[i] = int(math.Abs(float64(i/2 - current_floor))) //ABSOLUTT VERDI
}else if(direction == 1){
	if(i%2 == 1 && i/2 > current_floor){
		cost[i] = i/2 - current_floor - 1
	}else if (i%2 == 1 && i/2 <= current_floor || i%2 == 0){
		cost[i] =int (math.Abs(float64(i/2 - destination)) + math.Abs(float64(destination - current_floor - 1)))
	}else{
		cost[i] = 6
	}
}else{
	if(i%2 == 0 && i/2 < current_floor){
		cost[i] = current_floor - i/2 - 1
	}else if (i%2 == 0 && i/2 >= current_floor || i%2 == 1){
		cost[i] = int(math.Abs(float64(i/2 - destination)) + math.Abs(float64(current_floor- destination - 1)))
	}else{
		cost[i] = 6
	}
}
i += 1
}
return cost
}
func removeOrders(order_list [8]int, command_list [8]int, direction int, destination int)([8]int,[8]int){
	i := 0
	for (i < 4){
		if (driver.GetFloor() == i){
//driver.SetButtonLamp("command", i, 0)
			command_list[i] = 0
			if (destination == i){
				command_list[i] = 0
				order_list[i*2] = 0
				order_list[i*2+1] = 0
				else if (direction == 1){
//driver.SetButtonLamp("up", i, 0)
					order_list[i*2+1] = 0
				} else if (direction == 0){
//driver.SetButtonLamp("down", i, 0)
					order_list[i*2] = 0
				}
			}
			i +=1
		}
		return order_list, command_list
	}
}
