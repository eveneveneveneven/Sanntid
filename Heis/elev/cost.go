package elev

import (
	. "../types"
	"fmt"
)

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

var optimal_id int

func costFunction(network_msg *NetworkMessage) *Order {
	noInternal := true
	for _, order := range network_msg.Statuses[network_msg.Id].InternalOrders {
		if order != -1 {
			noInternal = false
			break
		}
	}
	if len(network_msg.Orders) == 0 && noInternal {
		return nil
	}

	number_of_elevs := len(network_msg.Statuses)
	number_of_orders := len(network_msg.Orders) + number_of_elevs

	var cost_id_order [10][14]int
	j := 0
	//taken_order_number :=0
	for order := range network_msg.Orders {
		i := 2
		order_button_type := order.ButtonPress
		order_button_floor := order.Floor
		order_button_value := 0

		if order_button_type == BUTTON_CALL_UP {
			order_button_value = UP
		} else {
			order_button_value = DOWN
		}
		//lowest_total_cost := 100
		optimal_id = -1

		//fmt.Printf("Order: %+v\n", order)

		for id, elevStat := range network_msg.Statuses {
			diff_cost := 0
			dir_cost := 0
			order_dir_cost := 0
			//internal_orders_cost :=0
			//fmt.Println(elevStat.InternalOrders)
			internal_order := elevStat.InternalOrders[0]
			/*for n := range elevStat.InternalOrders{
				if elevStat.InternalOrders[n] != -1{
					internal_orders_cost+=5
				}
			}*/
			//multiple_orders_cost :=0
			dir_order := 0
			dir_elev := elevStat.Dir
			floor_elev := elevStat.Floor
			dir_internal_order := -1
			if internal_order-floor_elev < 0 && internal_order != -1 {
				dir_internal_order = 1
			} else if internal_order-floor_elev > 0 && internal_order != -1 {
				dir_internal_order = 0
			} else {
				dir_internal_order = -1
			}
			//internal_orders_elev := elevStat.InternalOrders
			diff_cost = abs(order_button_floor - floor_elev)
			if order_button_floor > floor_elev {
				dir_order = UP
			} else {
				dir_order = DOWN
			}

			if dir_order == dir_elev || dir_elev == STOP {
				dir_cost = 0
			} else {
				dir_cost += 5
			}
			if internal_order != -1 {
				if order_button_value == dir_internal_order &&
					dir_internal_order == 1 && order_button_floor >= internal_order {
					if floor_elev == order_button_floor && dir_elev != STOP {
						order_dir_cost = 7
					} else {
						order_dir_cost = -7
					}
				} else if order_button_value == dir_internal_order &&
					dir_internal_order == 0 && order_button_floor <= internal_order {
					if floor_elev == order_button_floor && dir_elev != STOP {
						order_dir_cost = 7
					} else {
						order_dir_cost = -7
					}
				}
			}
			if internal_order == -1 {
				if dir_elev == STOP {
					order_dir_cost = 0
				} else if dir_elev != STOP && order_button_floor == floor_elev {
					order_dir_cost = 7
				} else {
					order_dir_cost = 3
				}
			}

			/*for i:=0; i<10; i++ {
				if id==active_orders[i]{
					multiple_orders_cost +=8
				}
			}*/

			total_cost := dir_cost + diff_cost + order_dir_cost
			//fmt.Println("Id: ", id, "direction: ", dir_elev,  ", dir int: ",
			//	dir_internal_order, ", Total cost: ", total_cost, "\n")
			/*if total_cost< lowest_total_cost {
				lowest_total_cost = total_cost
				optimal_id = id
			}*/
			cost_id_order[0][j] = order_button_floor
			cost_id_order[1][j] = order_button_value
			cost_id_order[i][j] = total_cost
			if internal_order != -1 {
				cost_id_order[id+2][j+number_of_orders-number_of_elevs] = -2
			} else {
				cost_id_order[id+2][j+number_of_orders-number_of_elevs] = 100
			}
			i++

		}
		j++
	}
	cost_id_order[0][j] = network_msg.Statuses[0].InternalOrders[0]
	cost_id_order[1][j] = BUTTON_INTERNAL
	//print_matrix(cost_id_order, number_of_elevs, number_of_orders)
	return sort_according_to_ordered_values(network_msg.Id, cost_id_order,
		number_of_elevs, number_of_orders)

}
func print_matrix(cost_matrix [10][14]int, num_elevs int, num_orders int) {
	for i := 0; i < num_elevs+2; i++ {
		for j := 0; j < num_orders; j++ {
			fmt.Print(cost_matrix[i][j], ", ")
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

func switch_rows(cost_matrix [10][14]int, switching_floor int,
	lowest_order int) [10][14]int {
	var temp_array [10]int
	for i := 0; i < 10; i++ {
		temp_array[i] = cost_matrix[i][switching_floor]
		cost_matrix[i][switching_floor] = cost_matrix[i][lowest_order]
	}
	for i := 0; i < 10; i++ {
		cost_matrix[i][lowest_order] = temp_array[i]
	}
	return cost_matrix
}

func sort_according_to_ordered_values(id int, cost_matrix [10][14]int,
	num_elevs int, num_orders int) *Order {
	lowest_order := 10
	for i := 0; i < num_orders-4; i++ {
		lowest_floor := 10
		lowest_dir := 3
		for j := i; j < num_orders-4; j++ {
			if cost_matrix[0][j] < lowest_floor {
				lowest_floor = cost_matrix[0][j]
				lowest_order = j
				lowest_dir = cost_matrix[1][j]
			}
			if cost_matrix[0][j] == lowest_floor {
				if cost_matrix[1][j] < lowest_dir {
					lowest_order = j
				}
			}
		}
		cost_matrix = switch_rows(cost_matrix, i, lowest_order)
	}
	//fmt.Println("\n")
	//print_matrix(cost_matrix, num_elevs, num_orders)
	return smallest_total_cost(id, cost_matrix, num_elevs, num_orders)
}

func add_multiple_orders_penalty(cost_matrix [10][14]int, num_elevs int,
	num_orders int, best_id int, order_taken int) [10][14]int {
	for j := 0; j < num_orders; j++ {
		cost_matrix[best_id][j] += 100
	}
	for i := 0; i < num_elevs; i++ {
		//print_matrix(cost_matrix, num_elevs, num_orders)
		cost_matrix[i+2][order_taken] += 100
	}
	return cost_matrix
}

func smallest_total_cost(id int, cost_matrix [10][14]int, num_elevs int,
	num_orders int) *Order {
	for k := 0; k < num_elevs; k++ {
		smallest_cost := 100
		best_id := -1
		order_taken := -1
		for k := 0; k < num_orders; k++ {
			smallest_cost = 100
			for i := 2; i < num_elevs+2; i++ {
				for j := 0; j < num_orders; j++ {
					if cost_matrix[i][j] < smallest_cost {
						smallest_cost = cost_matrix[i][j]
						best_id = i
						if j > num_orders-num_elevs {
							order_taken = j + best_id - 2
						} else {
							order_taken = j
						}

					}
				}
			}
		}
		//fmt.Println("Elevator ", best_id-2, " takes the order to floor ",
		//	cost_matrix[0][order_taken], ", in direction ",
		//	cost_matrix[1][order_taken],", cost: ", smallest_cost)
		//print_matrix(cost_matrix, num_elevs, num_orders)
		//fmt.Println("Order taken", order_taken)
		if id == best_id-2 {
			if order_taken >= num_orders-num_elevs {
				o := &Order{
					ButtonPress: BUTTON_INTERNAL,
					Floor:       cost_matrix[0][order_taken],
					Completed:   false,
				}
				return o
			} else {
				o := &Order{
					ButtonPress: cost_matrix[1][order_taken],
					Floor:       cost_matrix[0][order_taken],
					Completed:   false,
				}
				return o
			}
		}
		cost_matrix = add_multiple_orders_penalty(cost_matrix, num_elevs, num_orders, best_id, order_taken)
	}
	return nil
}
