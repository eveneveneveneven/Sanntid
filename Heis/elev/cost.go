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

func validateNetstat(netstat *NetworkMessage) bool {
	numElevs := len(netstat.Statuses)
	for id := 0; id < numElevs; id++ {
		if _, ok := netstat.Statuses[id]; !ok {
			return false
		}
		if len(netstat.Statuses[id].InternalOrders) != 4 {
			return false
		}
	}
	return true
}

func costFunction(network_msg *NetworkMessage) *Order {
	if !validateNetstat(network_msg) {
		fmt.Println("\t\x1b[31;1mError\x1b[0m |costFunction| [Received network_msg is not valid]")
		return nil
	}
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
	cost_id_order := make([][]int, number_of_elevs+2)
	internal_orders := make(map[int][]int, number_of_elevs)
	for i := 0; i < number_of_elevs+2; i++ {
		cost_id_order[i] = make([]int, number_of_orders)
	}
	for id, elevStat := range network_msg.Statuses {
		internal_orders[id] = make([]int, 5)
		if len(elevStat.InternalOrders) < 4 {
			elevStat.InternalOrders = []int{-1, -1, -1, -1}
			network_msg.Statuses[id] = elevStat
		}
		//fmt.Println("id ::", id, ":: elevstat ::", elevStat)
		internal_orders[id][0] = elevStat.Floor
		for i := 1; i < 5; i++ {
			internal_orders[id][i] = elevStat.InternalOrders[i-1]
		}
	}

	dirs := make([]int, number_of_elevs)

	j := 0
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

		for id, elevStat := range network_msg.Statuses {
			diff_cost := 0
			dir_cost := 0
			order_dir_cost := 0
			internal_order := elevStat.InternalOrders[0]

			dir_order := 0
			dir_elev := elevStat.Dir
			floor_elev := elevStat.Floor
			internal_orders[id][0] = floor_elev
			dirs[id] = dir_elev
			dir_internal_order := -1
			if internal_order-floor_elev < 0 && internal_order != -1 {
				dir_internal_order = 1
			} else if internal_order-floor_elev > 0 && internal_order != -1 {
				dir_internal_order = 0
			} else {
				dir_internal_order = 2
			}
			diff_cost = abs(order_button_floor - floor_elev)
			if order_button_floor > floor_elev {
				dir_order = UP
			} else if order_button_floor < floor_elev {
				dir_order = DOWN
			} else {
				dir_order = STOP
			}

			if dir_order == dir_elev || dir_elev == STOP {
				dir_cost = 0
			} else {
				dir_cost = 55
			}
			if internal_order != -1 {
				if dir_order == dir_internal_order &&
					dir_internal_order == 1 && order_button_floor > internal_order {
					if (floor_elev == order_button_floor && dir_elev != STOP) ||
						order_button_value != dir_internal_order {
						order_dir_cost = 12
					} else {
						order_dir_cost = -8
					}
				} else if dir_order == dir_internal_order &&
					dir_internal_order == 0 && order_button_floor < internal_order {
					if (floor_elev == order_button_floor && dir_elev != STOP) ||
						order_button_value != dir_internal_order {
						order_dir_cost = 12
					} else {
						order_dir_cost = -8

					}
				} else {
					order_dir_cost = 18
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
				internal_dir_change_cost = 60
			} else {
				internal_dir_change_cost = 0
			}
			cost_id_order[0][j] = order_button_floor
			cost_id_order[1][j] = order_button_value
			cost_id_order[id+2][j] = total_cost
			internal_dir_change[id] = internal_dir_change_cost
		}
		j++
	}

	for id, elevstat := range network_msg.Statuses {
		cost_id_order[0][id+number_of_orders-number_of_elevs] = elevstat.InternalOrders[0]
		cost_id_order[1][id+number_of_orders-number_of_elevs] = BUTTON_INTERNAL
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
	cost_id_order[0][j] = network_msg.Statuses[0].InternalOrders[0]
	cost_id_order[1][j] = BUTTON_INTERNAL
	return sort_according_to_ordered_values(network_msg, network_msg.Id, cost_id_order,
		number_of_elevs, number_of_orders, dirs, internal_orders)

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

func switch_rows(cost_matrix [][]int, switching_floor int, lowest_order int) [][]int {
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

func sort_according_to_ordered_values(network_msg *NetworkMessage, id int, cost_matrix [][]int,
	num_elevs int, num_orders int, dirs []int,
	internal_orders map[int][]int) *Order {
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
	cost_matrix = check_for_similar_buttons(network_msg, cost_matrix, num_elevs, num_orders, dirs)
	return smallest_total_cost(id, cost_matrix, num_elevs, num_orders, internal_orders)
}

func check_for_similar_buttons(network_msg *NetworkMessage, cost_matrix [][]int,
	num_elevs int, num_orders int, dirs []int) [][]int {
	for id := range network_msg.Statuses {
		number_of_occuring_down_values := 0
		number_of_occuring_up_values := 0
		for i := 0; i < num_orders-num_elevs; i++ {
			if cost_matrix[1][num_orders-num_elevs-i-1] == 1 && dirs[id] != 1 {
				for k := 2; k < num_elevs+2; k++ {
					cost_matrix[k][num_orders-num_elevs-i-1] += number_of_occuring_down_values * 4
				}
				number_of_occuring_down_values += 1
			}
			if cost_matrix[1][i] == 0 && dirs[id] != 0 {
				for k := 2; k < num_elevs+2; k++ {
					cost_matrix[k][i] += number_of_occuring_up_values * 4
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
		//fmt.Println(best_id, j)
		cost_matrix[best_id][j] += 100
	}
	for i := 0; i < num_elevs; i++ {
		cost_matrix[i+2][order_taken] += 100
	}
	return cost_matrix
}

func check_for_favorable_internals(internal_orders map[int][]int, id int) int {
	goal_floor := internal_orders[id][1]
	elev_floor := internal_orders[id][0]
	if elev_floor-goal_floor < 0 {
		for i := 2; i < 5; i++ {
			if internal_orders[id][i] != -1 && internal_orders[id][i] < goal_floor && elev_floor < internal_orders[id][i] {
				goal_floor = internal_orders[id][i]
			}
		}
		return goal_floor
	} else {
		for i := 2; i < 5; i++ {
			if internal_orders[id][i] != -1 && internal_orders[id][i] > goal_floor && internal_orders[id][i] < elev_floor {
				goal_floor = internal_orders[id][i]
			}
		}
		return goal_floor
	}
}

func smallest_total_cost(id int, cost_matrix [][]int, num_elevs int,
	num_orders int, internal_orders map[int][]int) *Order {
	num_active_orders := 0
	for i := 0; i < num_elevs; i++ {
		if cost_matrix[0][num_orders-num_elevs+i] != -1 {
			num_active_orders += 1
		}
	}
	num_active_orders += num_orders - num_elevs + 1
	for k := 0; k < num_active_orders; k++ {
		smallest_cost := 200
		best_id := -1
		order_taken := -1
		for k := 0; k < num_orders; k++ {
			smallest_cost = 200
			for i := 2; i < num_elevs+2; i++ {
				for j := 0; j < num_orders; j++ {
					if cost_matrix[i][j] < smallest_cost {
						smallest_cost = cost_matrix[i][j]
						best_id = i
						order_taken = j

					}
				}
			}
		}
		//print_matrix(cost_matrix, num_elevs, num_orders)
		if smallest_cost < 50 {
			if id == best_id-2 {
				if order_taken >= num_orders-num_elevs {
					goal_floor := check_for_favorable_internals(internal_orders, id)
					o := &Order{
						ButtonPress: BUTTON_INTERNAL,
						Floor:       goal_floor,
					}
					//fmt.Println("Sending internal order to id:", id, ", Order: ", goal_floor, BUTTON_INTERNAL)
					//print_matrix(cost_matrix, num_elevs, num_orders)
					return o
				} else {
					o := &Order{
						ButtonPress: cost_matrix[1][order_taken],
						Floor:       cost_matrix[0][order_taken],
					}
					//fmt.Println("Sending order to id:", id, ", Order: ", cost_matrix[0][order_taken], cost_matrix[1][order_taken])
					//print_matrix(cost_matrix, num_elevs, num_orders)
					return o
				}
			}
		} else {
			//print_matrix(cost_matrix, num_elevs, num_orders)
			//fmt.Println("Returning nil")
			return nil
		}
		cost_matrix = add_multiple_orders_penalty(cost_matrix, num_elevs, num_orders, best_id, order_taken)
	}
	return nil
}
