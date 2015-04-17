package network

import (
	"fmt"
	"os"
	"time"

	"../types"
)

type Hub struct {
	id int

	becomeMaster chan bool

	networkStatus *types.NetworkMessage
	netMsgUpd     *types.NetworkMessage

	foundMaster   chan string
	missingMaster chan bool

	messageRecieved chan *types.NetworkMessage
	messageSend     chan *types.NetworkMessage

	netstatNewMsg chan *types.NetworkMessage
	netstatUpdate chan *types.NetworkMessage
}

func NewHub(becomeMaster chan bool,
	netStatSend, netStatRec chan *types.NetworkMessage) *Hub {
	return &Hub{
		id: -1,

		becomeMaster: becomeMaster,

		networkStatus: &types.NetworkMessage{
			Id: -1,
			Statuses: []types.ElevStat{
				*types.NewElevStat(),
			},
			Orders: make(map[types.Order]struct{}),
		},
		netMsgUpd: &types.NetworkMessage{
			Id: -1,
			Statuses: []types.ElevStat{
				*types.NewElevStat(),
			},
			Orders: make(map[types.Order]struct{}),
		},

		foundMaster:   make(chan string),
		missingMaster: make(chan bool),

		messageRecieved: make(chan *types.NetworkMessage),
		messageSend:     make(chan *types.NetworkMessage),

		netstatNewMsg: netStatSend,
		netstatUpdate: netStatRec,
	}
}

func (h *Hub) Run() {
	fmt.Println("Start network Hub!")
	cm := NewConnManager(h.messageRecieved, h.messageSend)
	go cm.run()
	go startUDPListener(h.foundMaster, h.missingMaster)

	// Slave loop
	connected := false
slaveloop:
	for {
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
		case netStatUpd := <-h.netstatUpdate:
			h.netMsgUpd = netStatUpd
		}
	}

	close(h.becomeMaster)

	// Master loop
	go startUDPBroadcast()
	tick := time.Tick(types.SEND_INTERVAL * time.Millisecond)
	for {
		select {
		case netStatUpd := <-h.netstatUpdate:
			types.Clone(h.networkStatus, netStatUpd)
		case msgRec := <-h.messageRecieved:
			h.netstatNewMsg <- msgRec
		case <-tick:
			h.netstatNewMsg <- nil
			h.messageSend <- h.networkStatus
		}
	}
}

func (h *Hub) parseMessage(msg *types.NetworkMessage) {
<<<<<<< HEAD
=======
	fmt.Printf("Recieved: %+v\n", msg)
>>>>>>> 23d545ecd2d0f08b3f30f1b40de27a25060d8f4c
	h.id = msg.Id
	h.networkStatus = msg
	h.netstatNewMsg <- msg
	h.messageSend <- h.netMsgUpd
}
