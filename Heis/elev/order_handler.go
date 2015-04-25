package elev

import (
	"fmt"

	"../types"
)

type orderHandler struct {
	currObj *types.Order
	lastObj *types.Order

	newNetwork chan types.NetworkMessage
	sendNewObj chan types.Order
}

func newOrderHandler(newNetworkCh chan types.NetworkMessage,
	sendNewObjCh chan types.Order) *orderHandler {
	return &orderHandler{
		currObj: nil,
		lastObj: nil,

		newNetwork: newNetworkCh,
		sendNewObj: sendNewObjCh,
	}
}

func (oh *orderHandler) run() {
	fmt.Println("\x1b[34;1m::: Start Order Handler :::\x1b[0m")
	for {
		select {
		case newNetwork := <-oh.newNetwork:
			oh.parseNewNetwork(&newNetwork)
		}
	}
}

func (oh *orderHandler) parseNewNetwork(netstat *types.NetworkMessage) {
	oh.currObj = costFunction(netstat)
	if oh.currObj == nil {
		oh.lastObj = nil
		return
	}
	if oh.lastObj == nil || *oh.currObj != *oh.lastObj {
		oh.sendNewObj <- *oh.currObj
		oh.lastObj = oh.currObj
	}
}
