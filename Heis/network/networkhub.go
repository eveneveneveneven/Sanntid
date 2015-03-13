package network

import (
	"fmt"
	"os"
	"time"

	"../types"
)

type Hub struct {
	master        bool
	id            int
	networkStatus *types.NetworkMessage

	foundMaster   chan string
	missingMaster chan bool

	messageRecieved chan *types.NetworkMessage
	messageSend     chan *types.NetworkMessage

	statusRecieve chan *types.NetworkMessage
	statusSend    chan *types.NetworkMessage
}

func NewHub(statRec, statSend chan *types.NetworkMessage) *Hub {
	return &Hub{
		master: false,
		id:     -1,
		networkStatus: &types.NetworkMessage{
			Id:     -1,
			Status: "",
			Orders: "",
		},

		foundMaster:   make(chan string),
		missingMaster: make(chan bool),

		messageRecieved: make(chan *types.NetworkMessage),
		messageSend:     make(chan *types.NetworkMessage),

		statusRecieve: statRec,
		statusSend:    statSend,
	}
}

func (h *Hub) Run() {
	fmt.Println("Start network Hub!")
	cm := NewConnManager(h.messageRecieved, h.messageSend)
	go cm.run()
	go startUDPListener(h.foundMaster, h.missingMaster)

	// Slave loop
	connected := false
	for !h.master {
		select {
		case masterIp := <-h.foundMaster:
			if connected {
				continue
			}
			if err := cm.connectToNetwork(masterIp); err != nil {
				fmt.Printf("\x1b[31;1mError\x1b[0m |Hub.Run| [%v], exit program\n", err)
				os.Exit(1)
			}
			connected = true
		case <-h.missingMaster:
			switch h.id {
			case 1:
				fmt.Println("Master is dead, I am Master!")
				h.master = true
				h.id = 0
			case -1:
				fmt.Println("There is no Master, claim Master!")
				h.master = true
				h.id = 0
			default:
				fmt.Println("Master is dead, continue as slave...")
				h.id--
			}
			connected = false
		case msgRecieve := <-h.messageRecieved:
			fmt.Printf("Recieved: %+v\n", msgRecieve)
			h.parseMessage(msgRecieve)
			h.messageSend <- h.getResponse()
		}
	}

	// Master loop
	go startUDPBroadcast()
	timer := time.NewTimer(types.SEND_INTERVAL * time.Millisecond)
	for {
		select {
		case msgRecieve := <-h.messageRecieved:
			h.parseMessage(msgRecieve)

		case <-timer.C:
			timer.Reset(types.SEND_INTERVAL * time.Millisecond)
			h.messageSend <- h.getNextMessage()
		}
	}
}

func (h *Hub) parseMessage(msg *types.NetworkMessage) {
	h.id = msg.Id
	h.networkStatus = msg
}

func (h *Hub) getResponse() *types.NetworkMessage {
	return &types.NetworkMessage{
		Id:     1,
		Status: "yo",
		Orders: "none",
	}
}

func (h *Hub) getNextMessage() *types.NetworkMessage {
	return &types.NetworkMessage{
		Id:     1,
		Status: "yoman",
		Orders: "everything",
	}
}
