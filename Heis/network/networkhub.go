package network

import (
	"fmt"
	"os"
	"time"
)

type Hub struct {
	master        bool
	id            int
	networkStatus *networkMessage

	foundMaster   chan string
	missingMaster chan bool

	messageRecieved chan *networkMessage
	messageSend     chan *networkMessage
}

func NewHub() *Hub {
	var h Hub

	h.master = false
	h.id = -1
	h.networkStatus = &networkMessage{
		Id:     -1,
		Status: "",
		Orders: "",
	}

	h.foundMaster = make(chan string)
	h.missingMaster = make(chan bool)

	h.messageRecieved = make(chan *networkMessage)
	h.messageSend = make(chan *networkMessage)

	return &h
}

func (h *Hub) Run() {
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
				fmt.Printf("Some error %v, exit program\n", err)
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
		case msgRecieve := <-h.messageRecieved:
			h.parseMessage(msgRecieve)
			h.messageSend <- h.getResponse()
		}
	}

	timer := time.NewTimer(250 * time.Millisecond)
	go startUDPBroadcast()
	// Master loop
	for {
		select {
		case msgRecieve := <-h.messageRecieved:
			h.parseMessage(msgRecieve)

		case <-timer.C:
			timer.Reset(250 * time.Millisecond)
			h.messageSend <- h.getNextMessage()
		}
	}
}

func (h *Hub) parseMessage(msg *networkMessage) {
	h.id = msg.Id
	h.networkStatus = msg
}

func (h *Hub) getResponse() *networkMessage {
	return &networkMessage{
		Id:     1,
		Status: "yo",
		Orders: "none",
	}
}

func (h *Hub) getNextMessage() *networkMessage {
	return &networkMessage{
		Id:     1,
		Status: "yoman",
		Orders: "everything",
	}
}
