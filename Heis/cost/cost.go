package cost

import (
	."../types"
	"fmt"
)

func abs (value int) (int) {
	if value<0 {
		return -value
	}
	return value
}

var optimal_id int

func Cost_function(network_msg *NetworkMessage){
	number_of_orders :=  len(network_msg.Orders)
	number_of_elevs :=  cap(network_msg.Statuses)

	var cost_id_order [10][10]int
	j:=0
	//taken_order_number :=0
	for order := range network_msg.Orders{
		i:=2
		order_button_type := order.ButtonPress
		order_button_floor := order.Floor
		order_button_value := 0

		if order_button_type == BUTTON_CALL_UP {
			order_button_value = UP
		}else{
			order_button_value = DOWN
		}
		//lowest_total_cost := 100
		optimal_id = -1

		fmt.Printf("Order: %+v\n", order)

		for id, elevStat := range network_msg.Statuses {
			diff_cost :=0
			dir_cost :=0
			order_dir_cost :=0
			//multiple_orders_cost :=0
			dir_order := 0
			dir_elev := elevStat.Dir
			floor_elev := elevStat.Floor
			//internal_orders_elev := elevStat.InternalOrders

			diff_cost = abs(order_button_floor - floor_elev)
			if order_button_floor>floor_elev{
				dir_order = UP
			}else {
				dir_order = DOWN
			}

			if dir_order == dir_elev || dir_elev==STOP {
				dir_cost = 0
			}else{
				dir_cost += 5
			}

			if order_button_value == dir_elev && order_button_floor!=floor_elev{
				order_dir_cost =-1
			}else if dir_elev == STOP{
				order_dir_cost =0
			}else if dir_elev != STOP && order_button_floor == floor_elev{
				order_dir_cost = 7
			}else{
				order_dir_cost =3
			}

			/*for i:=0; i<10; i++ {
				if id==active_orders[i]{
					multiple_orders_cost +=8
				}
			}*/

			total_cost := dir_cost + diff_cost + order_dir_cost
			fmt.Println("Id: ", id, ", Total cost: ", total_cost, "\n")
			/*if total_cost< lowest_total_cost {
				lowest_total_cost = total_cost
				optimal_id = id
			}*/
			cost_id_order[0][j] = order_button_floor
			cost_id_order[1][j] = order_button_value
			cost_id_order[i][j] = total_cost
			i++

		}
		j++
	}
	//print_matrix(cost_id_order, number_of_elevs, number_of_orders)
	sort_according_to_ordered_values(cost_id_order, number_of_elevs, number_of_orders)

}
/*func print_matrix(cost_matrix [10][10]int, num_elevs int, num_orders int){
	for i:=0; i<num_elevs+2; i++{
		for j:=0; j<num_orders; j++{
			fmt.Print(cost_matrix[i][j], ", ")
		}
		fmt.Print("\n")
	}
}*/

func switch_rows(cost_matrix [10][10] int, switching_floor int, lowest_order int) ([10][10]int){
	var temp_array [10]int
	for i:=0; i<10; i++{
		temp_array[i] = cost_matrix[i][switching_floor]
		cost_matrix[i][switching_floor] = cost_matrix[i][lowest_order]
	}
	for i:=0; i<10; i++{
		cost_matrix[i][lowest_order] = temp_array[i]
	}
	return cost_matrix
}

func sort_according_to_ordered_values(cost_matrix [10][10]int, num_elevs int, num_orders int){
	lowest_order :=10
	for i:=0; i<num_orders; i++{
		lowest_floor:=10
		lowest_dir:=3
		for j:=i; j<num_orders; j++{
			if cost_matrix[0][j]<lowest_floor{
				lowest_floor = cost_matrix[0][j]
				lowest_order = j
				lowest_dir = cost_matrix[1][j]
			}
			if cost_matrix[0][j] == lowest_floor{
				if cost_matrix[1][j] < lowest_dir{
					lowest_order = j
				}
			}
		}
		cost_matrix = switch_rows(cost_matrix, i, lowest_order)
	}
	fmt.Println("\n")
	//print_matrix(cost_matrix, num_elevs, num_orders)
	smallest_total_cost(cost_matrix, num_elevs, num_orders)
}

func add_multiple_orders_penalty(cost_matrix [10][10]int, num_elevs int, num_orders int, best_id int, order_taken int) ([10][10]int){
	for j:=0; j<num_orders; j++{
		cost_matrix[best_id][j] +=5
		cost_matrix[j+2][order_taken]+=100
	}
	return cost_matrix
}

func smallest_total_cost(cost_matrix [10][10]int, num_elevs int, num_orders int){
	for k:=0; k<num_orders; k++{
		smallest_cost:=100
		best_id:=-1
		order_taken:=-1
		for k:=0; k<num_orders; k++{
			smallest_cost = 100
			for i:=2; i<num_elevs+2; i++{
				for j:=0; j<num_orders; j++{
					if cost_matrix[i][j]<smallest_cost{
						smallest_cost = cost_matrix[i][j]
						best_id = i
						order_taken = j
					}
				}
			}
		}
		fmt.Println("Elevator ", best_id-1, " takes the order to floor ", cost_matrix[0][order_taken], ", in direction ", cost_matrix[1][order_taken])

		cost_matrix = add_multiple_orders_penalty(cost_matrix, num_elevs, num_orders, best_id, order_taken)
		fmt.Println("\n")
		//print_matrix(cost_matrix, num_elevs, num_orders)
	}
	
}