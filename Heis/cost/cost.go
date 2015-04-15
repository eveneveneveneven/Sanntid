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
	
	

	for order := range network_msg.Orders{
		order_button_type := order.ButtonPress
		order_button_floor := order.Floor
		order_button_value := 0

		if order_button_type == BUTTON_CALL_UP {
			order_button_value = UP
		}else{
			order_button_value = DOWN
		}
		lowest_total_cost := 100
		optimal_id = -1

		fmt.Printf("Order: %+v\n", order)

		for id, elevStat := range network_msg.Statuses {
			diff_cost :=0
			dir_cost :=0
			order_dir_cost :=0
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

			total_cost := dir_cost + diff_cost + order_dir_cost
			fmt.Println("Id: ", id, ", Total cost: ", total_cost, "\n")
			if total_cost< lowest_total_cost {
				lowest_total_cost = total_cost
				optimal_id = id
			}
		}
		fmt.Println("Best elevator: ", optimal_id, ", cost: ", + lowest_total_cost, "\n\n")

	}

}