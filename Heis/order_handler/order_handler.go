package order_handler

import (
	"../types"
)

type OrderHandler struct {
	networkStatus *types.NetworkMessage

	networkRecieve chan *types.NetworkMessage // connection from network module
	networkSend    chan *types.NetworkMessage // connection to network module

	elevGoToFloor chan int // connection to elevator module
	orderDone     chan int // connection from elevator module

	newOrder        chan int // connection from button module
	clearOrderLight chan int // connection to light module
}

func NewOrderHandler(statRec, statSend chan *types.NetworkMessage) *OrderHandler {
	return &OrderHandler{
		networkStatus: new(types.NetworkMessage),

		networkRecieve: statRec,
		networkSend:    statSend,

		elevGoToFloor: nil,
		orderDone:     nil,

		newOrder:        nil,
		clearOrderLight: nil,
	}
}

func (oh *OrderHandler) Run() {
	for {
		select {}
	}
}
