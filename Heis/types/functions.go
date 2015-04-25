package types

import (
	"strconv"
)

func DeepCopy(dst, src *NetworkMessage) {
	dst.Id = src.Id
	dst.Statuses = make(map[int]ElevStat)
	for id, elev := range src.Statuses {
		dst.Statuses[id] = elev
	}
	dst.Orders = make(map[Order]bool)
	for order, active := range src.Orders {
		dst.Orders[order] = active
	}
}

func NewElevStat() *ElevStat {
	return &ElevStat{
		Dir:            STOP,
		Floor:          -1,
		InternalOrders: []int{-1, -1, -1, -1},
	}
}

func NewNetworkMessage() *NetworkMessage {
	return &NetworkMessage{
		Id:       -1,
		Statuses: make(map[int]ElevStat),
		Orders:   make(map[Order]bool),
	}
}

func (o Order) String() string {
	str := "Button="
	switch o.ButtonPress {
	case BUTTON_CALL_UP:
		str += "UP"
	case BUTTON_CALL_DOWN:
		str += "DOWN"
	case BUTTON_INTERNAL:
		str += "INTERNAL"
	}
	str += " Floor=" + strconv.Itoa(o.Floor+1)
	return str
}
