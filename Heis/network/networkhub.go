package network

import (
	"fmt"
	"os"
	"time"

	"../types"
)

type NetworkHub struct {
	id int

	becomeMaster chan bool

	networkStatus *types.NetworkMessage

	foundMaster   chan string
	missingMaster chan bool

	messageRecieved chan *types.NetworkMessage
	messageSend     chan *types.NetworkMessage

	netstatNewMsg chan *types.NetworkMessage
	netstatUpdate chan *types.NetworkMessage
}

func NewNetworkHub(becomeMaster chan bool,
	elevSendNewNetstat, elevRecUpd chan *types.NetworkMessage) *NetworkHub {
	nh := &NetworkHub{
		id: -1,

		becomeMaster: becomeMaster,

		networkStatus: types.NewNetworkMessage(),

		foundMaster:   make(chan string),
		missingMaster: make(chan bool),

		messageRecieved: make(chan *types.NetworkMessage),
		messageSend:     make(chan *types.NetworkMessage),

		netstatNewMsg: make(chan *types.NetworkMessage),
		netstatUpdate: make(chan *types.NetworkMessage),
	}
	return nh
}

func (h *NetworkHub) Run() {
	fmt.Println("Start network NetworkHub!")

	cm := newConnManager(h.messageRecieved, h.messageSend)
	go cm.run()
	go startUDPListener(h.foundMaster, h.missingMaster)

	connected := false
	// Slave loop
slaveloop:
	for {
		select {
		case masterIp := <-h.foundMaster:
			if connected {
				continue
			}
			if err := cm.connectToNetwork(masterIp); err != nil {
				fmt.Printf("\x1b[31;1mError\x1b[0m |NetworkHub.Run| [%v], exit program\n", err)
				os.Exit(1)
			}
			connected = true
		case <-h.missingMaster:
			switch h.id {
			case 1:
				fmt.Println("Master is dead, I am Master!")
				h.id = 0
				break slaveloop
			case -1:
				fmt.Println("There is no Master, claim Master!")
				h.id = 0
				break slaveloop
			default:
				fmt.Println("Master is dead, continue as slave...")
				h.id--
			}
			connected = false
		case msgRecieve := <-h.messageRecieved:
			h.parseMessage(msgRecieve)
		}
	}

	close(h.becomeMaster)

	go startUDPBroadcast()
	go newNetStatManager(h.netstatNewMsg, h.netstatUpdate).run()

	tick := time.Tick(types.SEND_INTERVAL * time.Millisecond)
	// Master loop
	for {
		select {
		case msgRec := <-h.messageRecieved:
			h.netstatNewMsg <- msgRec
		case <-tick:
			h.netstatUpdate <- nil
		case newNetstat := <-h.netstatUpdate:
			h.networkStatus = newNetstat
			h.messageSend <- h.networkStatus
		}
	}
}

func (h *NetworkHub) parseMessage(msg *types.NetworkMessage) {
	h.id = msg.Id
	h.networkStatus = msg
	h.netstatNewMsg <- msg
	h.messageSend <- h.netMsgUpd
}
