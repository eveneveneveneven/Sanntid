package order

import (
	"../types"
)

type OrderHandler struct {
	becomeMaster chan bool

	netstatCurrentNetwork <-chan *types.NetworkMessage
	netstatNotify         chan<- *types.NetworkMessage

	elevGiveNewObj  chan<- int
	elevObjComplete <-chan int
}

func NewOrderHandler(becomeMaster chan bool,
	netstatCurrNet, netstatNotify chan *types.NetworkMessage,
	elevNewObj, elevObjComp chan int) *OrderHandler {
	return &OrderHandler{
		becomeMaster: becomeMaster,

		netstatCurrentNetwork: netstatCurrNet,
		netstatNotify:         netstatNotify,

		elevGiveNewObj:  elevNewObj,
		elevObjComplete: elevObjComp,
	}
}

func (oh *OrderHandler) Run() {
	for {
		select {}
	}
}
