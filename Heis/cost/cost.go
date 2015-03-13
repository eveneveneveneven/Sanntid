package cost

import(
	. "../driver"
)

func abs (value int) (int) {
	if value<0 {
		return -value
	}
	return value
}

func Cost_function(orders [][]int, dir int, last_floor int, current_order int) (int){
	min_cost:=100
	floor_diff:=0
	cost:=0
	for i:=0; i<4; i++{
		if orders[i][0]==1{
			//Up
			if (dir == -1 && last_floor<i) || (dir == 1 && last_floor>i) {
				cost+= 10
			}else if (dir == -1 && last_floor>i) || (dir == 1 && last_floor <i) {
				cost+= 3
			}else if dir == 0{
				cost += 1
			}
			floor_diff = abs(last_floor-i)
			cost += floor_diff
			if current_order > i && i>last_floor {
				cost = 1
			}else if i == Heis_get_floor(){
				cost = 0
			}
		}
		if orders[i][1]==1{
			//Down
			if (dir == -1 && last_floor<i) || (dir == 1 && last_floor>i) {
				cost+= 10
			}else if (dir == -1 && last_floor>i) || (dir == 1 && last_floor <i) {
				cost+= 3
			}else if dir == 0{
				cost += 1
			}
			floor_diff = abs(last_floor-i)
			cost += floor_diff
			if current_order < i && i<last_floor && current_order != -1{
				cost = 1
			}else if i == Heis_get_floor(){
				cost = 0
			}
		}
		if orders[i][2]==1{
			//Internal
		}
		if cost<min_cost{
			min_cost = cost
		}
	}

	return cost

}