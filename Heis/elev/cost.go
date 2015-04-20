package elev

import (
	"fmt"

	. "../types"
)

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

var optimal_id int

func costFunction(network_msg *NetworkMessage) *Order {
	fmt.Println("inn netmsg ::", network_msg)
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
	fmt.Println("numelev ::", number_of_elevs)
	fmt.Println("numorder ::", number_of_orders)
	cost_id_order := make([][]int, number_of_elevs+2)
	for i := 0; i < number_of_elevs+2; i++ {
		cost_id_order[i] = make([]int, number_of_orders)
	}

	dirs := make([]int, number_of_elevs)

	j := 0
	//taken_order_number :=0
	internal_dir_change := make([]int, number_of_elevs)
	for order := range network_msg.Orders {
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
			dirs[id] = dir_elev
			dir_internal_order := -1
			if internal_order-floor_elev < 0 && internal_order != -1 {
				dir_internal_order = 1
			} else if internal_order-floor_elev > 0 && internal_order != -1 {
				dir_internal_order = 0
			} else {
				dir_internal_order = 2
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
					dir_internal_order == 1 && order_button_floor > internal_order {
					if floor_elev == order_button_floor && dir_elev != STOP {
						order_dir_cost = 12
					} else {
						order_dir_cost = -8
					}
				} else if order_button_value == dir_internal_order &&
					dir_internal_order == 0 && order_button_floor < internal_order {
					if floor_elev == order_button_floor && dir_elev != STOP {
						order_dir_cost = 12
					} else {
						order_dir_cost = -8

					}
				} else {
					order_dir_cost = 15
				}
			}
			if internal_order == -1 {
				if floor_elev == order_button_floor && dir_elev != STOP {
					order_dir_cost = 10
				} else {
					order_dir_cost = 8
				}
			}

			total_cost := dir_cost + diff_cost + order_dir_cost
			internal_dir_change_cost := 0
			if dir_internal_order != dir_elev && dir_elev != STOP && internal_order != -1 {
				internal_dir_change_cost = 40
			} else {
				internal_dir_change_cost = 0
			}
			fmt.Println("j ::", j)
			cost_id_order[0][j] = order_button_floor
			cost_id_order[1][j] = order_button_value
			cost_id_order[id+2][j] = total_cost
			internal_dir_change[id] = internal_dir_change_cost

		}
		j++
	}

	fmt.Println(cost_id_order)
	for id, elevstat := range network_msg.Statuses {
		if elevstat.InternalOrders[0] != -1 {
			cost_id_order[id+2][id+number_of_orders-number_of_elevs] = -3 + internal_dir_change[id]
		} else {
			cost_id_order[id+2][id+number_of_orders-number_of_elevs] = 100
		}
		for i := 0; i < number_of_elevs; i++ {
			if i != id {
				cost_id_order[id+2][i+number_of_orders-len(network_msg.Statuses)] = 100
			}
		}

	}
	fmt.Println(cost_id_order)
	cost_id_order[0][j] = network_msg.Statuses[0].InternalOrders[0]
	cost_id_order[1][j] = BUTTON_INTERNAL
	//print_matrix(cost_id_order, number_of_elevs, number_of_orders)
	return sort_according_to_ordered_values(network_msg.Id, cost_id_order,
		number_of_elevs, number_of_orders, dirs)

}
func print_matrix(cost_matrix [][]int, num_elevs int, num_orders int) {
	for i := 0; i < num_elevs+2; i++ {
		for j := 0; j < num_orders; j++ {
			fmt.Print(cost_matrix[i][j], ", ")
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

func switch_rows(cost_matrix [][]int, switching_floor int,
	lowest_order int) [][]int {
	temp_array := make([]int, len(cost_matrix))
	for i := 0; i < len(temp_array); i++ {
		temp_array[i] = cost_matrix[i][switching_floor]
		cost_matrix[i][switching_floor] = cost_matrix[i][lowest_order]
	}
	for i := 0; i < len(temp_array); i++ {
		cost_matrix[i][lowest_order] = temp_array[i]
	}
	return cost_matrix
}

func sort_according_to_ordered_values(id int, cost_matrix [][]int,
	num_elevs int, num_orders int, dirs []int) *Order {
	lowest_order := 10
	for i := 0; i < num_orders-num_elevs; i++ {
		lowest_floor := 10
		lowest_dir := 3
		for j := i; j < num_orders-num_elevs; j++ {
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
	cost_matrix = check_for_similar_buttons(cost_matrix, num_elevs, num_orders, dirs)
	return smallest_total_cost(id, cost_matrix, num_elevs, num_orders)
}

func check_for_similar_buttons(cost_matrix [][]int,
	num_elevs int, num_orders int, dirs []int) [][]int {
	for j := 0; j < num_elevs; j++ {
		number_of_occuring_down_values := 0
		number_of_occuring_up_values := 0
		for i := 0; i < num_orders-num_elevs; i++ {
			fmt.Println(cost_matrix[1][num_orders-num_elevs-i-1], dirs[j])
			if cost_matrix[1][num_orders-num_elevs-i-1] == 1 && dirs[j] != 1 {
				for k := 2; k < num_elevs+2; k++ {
					fmt.Println("sub'ed some ", num_orders-num_elevs-i-1)
					cost_matrix[k][num_orders-num_elevs-i-1] += number_of_occuring_down_values * 6
				}
				number_of_occuring_down_values += 1
			}
			if cost_matrix[1][i] == 0 && dirs[j] != 0 {
				fmt.Println("Added some ", i)
				for k := 2; k < num_elevs+2; k++ {
					cost_matrix[k][i] += number_of_occuring_up_values * 6
					fmt.Println("Added some ", i)
				}
				number_of_occuring_up_values += 1
			}
		}
	}
	return cost_matrix
}

func add_multiple_orders_penalty(cost_matrix [][]int, num_elevs int,
	num_orders int, best_id int, order_taken int) [][]int {
	for j := 0; j < num_orders; j++ {
		fmt.Println(best_id, j)
		cost_matrix[best_id][j] += 100
	}
	for i := 0; i < num_elevs; i++ {
		//print_matrix(cost_matrix, num_elevs, num_orders)
		cost_matrix[i+2][order_taken] += 100
	}
	return cost_matrix
}

func smallest_total_cost(id int, cost_matrix [][]int, num_elevs int,
	num_orders int) *Order {
	for k := 0; k < num_elevs; k++ {
		smallest_cost := 200
		best_id := -1
		order_taken := -1
		for k := 0; k < num_orders; k++ {
			smallest_cost = 200
			for i := 2; i < num_elevs+2; i++ {
				for j := 0; j < num_orders; j++ {
					if cost_matrix[i][j] < smallest_cost {
						fmt.Println("lakjslksjdfl")
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
		fmt.Println("Best id: ", best_id)
		//fmt.Println("Elevator ", best_id-2, " takes the order to floor ",
		//	cost_matrix[0][order_taken], ", in direction ",
		//	cost_matrix[1][order_taken],", cost: ", smallest_cost)
		print_matrix(cost_matrix, num_elevs, num_orders)
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
				fmt.Println("OrderCost: ", smallest_cost)
				return o
			}
		}
		cost_matrix = add_multiple_orders_penalty(cost_matrix, num_elevs, num_orders, best_id, order_taken)
	}
	return nil
}
