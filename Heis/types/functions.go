package types

func Clone(dst, src *NetworkMessage) {
	dst.Id = src.Id
	dst.Statuses = make(map[int]ElevStat)
	for id, elev := range src.Statuses {
		dst.Statuses[id] = elev
	}
	dst.Orders = make(map[Order]struct{})
	for order := range src.Orders {
		dst.Orders[order] = struct{}{}
	}
}

func NewElevStat() *ElevStat {
	return &ElevStat{
		Dir:            0,
		Floor:          -1,
		InternalOrders: []int{-1, -1, -1, -1},
	}
}

func NewNetworkMessage() *NetworkMessage {
	return &NetworkMessage{
		Id:       -1,
		Statuses: make(map[int]ElevStat),
		Orders:   make(map[Order]struct{}),
	}
}
