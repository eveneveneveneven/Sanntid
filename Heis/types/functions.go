package types

func Clone(dst, src *NetworkMessage) {
	dst.Id = src.Id
	length := len(src.Statuses)
	if length == 0 {
		dst.Statuses = nil
	} else {
		dst.Statuses = make([]ElevStat, length)
		for i, status := range src.Statuses {
			dst.Statuses[i] = status
		}
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
		Statuses: nil,
		Orders:   make(map[Order]struct{}),
	}
}
