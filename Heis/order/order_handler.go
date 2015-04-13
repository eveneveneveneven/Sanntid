package order

import (
	"../types"
)

type OrderHandler struct {
	networkStatus  *types.NetworkMessage
	elevatorStatus *types.ElevStat

	localConn chan *types.NetworkMessage // connection from network module

	elevGoToFloor chan int // connection to elevator module
	orderDone     chan int // connection from elevator module

	newOrder        chan int // connection from button module
	clearOrderLight chan int // connection to light module
}

func NewOrderHandler(lc chan *types.NetworkMessage) *OrderHandler {
	return &OrderHandler{
		networkStatus:  new(types.NetworkMessage),
		elevatorStatus: new(types.ElevStat),

		localConn: lc,

		elevGoToFloor: nil,
		orderDone:     nil,

		newOrder:        nil,
		clearOrderLight: nil,
	}
}

func (oh *OrderHandler) Run() {
	for {
		select {
		case networkUpdate := <-oh.localConn:
			types.Clone(oh.networkStatus, networkUpdate)
			response := &types.NetworkMessage{
				Id:        -1,
				Statuses:  make([]types.ElevStat, 1),
				Orders:    nil,
				NewOrders: nil,
			}
			response.Statuses[0] = *oh.elevatorStatus
			oh.localConn <- response
		}
	}
}
