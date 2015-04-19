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
	cm            *connManager

	foundMaster   chan string
	missingMaster chan bool

	msgRecieve    chan *types.NetworkMessage
	msgSendGlobal chan *types.NetworkMessage
	msgSendLocal  chan *types.NetworkMessage

	netstatNewMsg chan *types.NetworkMessage
	netstatUpdate chan *types.NetworkMessage
	netstatTick   chan bool
}

func NewNetworkHub(becomeMaster chan bool,
	sendLocalCh, recieveCh chan *types.NetworkMessage) *NetworkHub {
	nh := &NetworkHub{
		id: -1,

		becomeMaster: becomeMaster,

		networkStatus: types.NewNetworkMessage(),
		cm:            nil,

		foundMaster:   make(chan string),
		missingMaster: make(chan bool),

		msgRecieve:    recieveCh,
		msgSendGlobal: make(chan *types.NetworkMessage),
		msgSendLocal:  sendLocalCh,

		netstatNewMsg: make(chan *types.NetworkMessage),
		netstatUpdate: make(chan *types.NetworkMessage),
		netstatTick:   make(chan bool),
	}
	nh.cm = newConnManager(nh.msgRecieve, nh.msgSendGlobal)
	go startUDPListener(nh.foundMaster, nh.missingMaster)
	go nh.cm.run()
	return nh
}

func (nh *NetworkHub) Run() {
	fmt.Println("Start NetworkHub!")

	connected := false
	// Slave loop
slaveloop:
	for {
		select {
		case masterIp := <-nh.foundMaster:
			if connected {
				continue
			}
			if err := nh.cm.connectToNetwork(masterIp); err != nil {
				fmt.Printf("\x1b[31;1mError\x1b[0m |NetworkHub.Run| [%v], exit program\n", err)
				os.Exit(1)
			}
			connected = true
		case <-nh.missingMaster:
			switch nh.id {
			case 1:
				fmt.Println("Master is dead, I am Master!")
				nh.id = 0
				break slaveloop
			case -1:
				fmt.Println("There is no Master, claim Master!")
				nh.id = 0
				break slaveloop
			default:
				fmt.Println("Master is dead, continue as slave...")
				nh.id--
			}
			connected = false
		case msgRecieve := <-nh.msgRecieve:
			nh.parseMessage(msgRecieve)
		}
	}

	close(nh.becomeMaster)

	go startUDPBroadcast()
	go newNetStatManager(nh.netstatNewMsg, nh.netstatUpdate, nh.netstatTick).run()

	tick := time.Tick(types.SEND_INTERVAL * time.Millisecond)
	// Master loop
	for {
		select {
		case msgRec := <-nh.msgRecieve:
			nh.netstatNewMsg <- msgRec
		case <-tick:
			nh.netstatTick <- true
		case newNetstat := <-nh.netstatUpdate:
			nh.networkStatus = newNetstat
			nh.msgSendGlobal <- newNetstat
			nh.msgSendLocal <- newNetstat
		}
	}
}

func (nh *NetworkHub) parseMessage(msg *types.NetworkMessage) {
	nh.id = msg.Id
	nh.msgSendGlobal <- nh.networkStatus
	nh.networkStatus = msg
	nh.msgSendLocal <- msg
}
