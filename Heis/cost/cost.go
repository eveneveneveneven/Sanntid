package cost

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

func Cost_function(current_floor int, current_order int, dir int) int {
	diff_cost := abs(current_order - current_floor)
	dir_final := 0
	direction_cost := 0
	final_goal_floor := 0
	diff_cost += abs(final_goal_floor - current_order)

	if final_goal_floor-current_order < 0 {
		dir_final = -1
	} else {
		dir_final = 1
	}
	if dir == dir_final {
		direction_cost = 0
	} else {
		direction_cost = 5
	}

	cost := diff_cost + direction_cost

	return cost

}
