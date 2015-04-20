package elev

import (
	"fmt"

	"../types"
)

type orderHandler struct {
	currObj *types.Order
	lastObj *types.Order

	newNetwork chan *types.NetworkMessage
	sendNewObj chan *types.Order
}

func newOrderHandler(newNetworkCh chan *types.NetworkMessage,
	sendNewObjCh chan *types.Order) *orderHandler {
	return &orderHandler{
		currObj: nil,
		lastObj: nil,

		newNetwork: newNetworkCh,
		sendNewObj: sendNewObjCh,
	}
}

func (oh *orderHandler) run() {
	fmt.Println("Start OrderHandler!")
	for {
		newNetwork := <-oh.newNetwork
		oh.parseNewNetwork(newNetwork)
	}
}

func (oh *orderHandler) parseNewNetwork(netstat *types.NetworkMessage) {
	setActiveLights(netstat)
	oh.currObj = costFunction(netstat)
	if oh.currObj == nil {
		oh.lastObj = nil
		return
	}
	fmt.Printf("@@@> CurrObj : %v\n@@@> LastObj : %v\n", oh.currObj, oh.lastObj)
	if oh.lastObj == nil || *oh.currObj != *oh.lastObj {
		fmt.Println("new obj ::", oh.currObj)
		oh.sendNewObj <- oh.currObj
		oh.lastObj = oh.currObj
	}
}