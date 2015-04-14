package order

import (
	"../types"
)

type OrderHandler struct {
	currNetwork *types.NetworkMessage

	becomeMaster chan bool

	netstatCurrentNetwork <-chan *types.NetworkMessage

	elevGiveNewObj  chan<- *types.Order
	elevObjComplete <-chan *types.Order
}

func NewOrderHandler(becomeMaster chan bool, netstatCurrNet chan *types.NetworkMessage,
	elevNewObj, elevObjComp chan *types.Order) *OrderHandler {
	return &OrderHandler{
		becomeMaster: becomeMaster,

		netstatCurrentNetwork: netstatCurrNet,

		elevGiveNewObj:  elevNewObj,
		elevObjComplete: elevObjComp,
	}
}

func (oh *OrderHandler) Run() {
	for {
		select {}
	}
}
