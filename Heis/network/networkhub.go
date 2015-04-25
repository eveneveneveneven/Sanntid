package network

import (
	"fmt"
	"os"
	"time"

	"../types"
)

const (
	SEND_INTERVAL = 200 // milliseconds
)

type NetworkHub struct {
	reset chan bool

	id int

	networkStatus *types.NetworkMessage
	cm            *connManager

	foundMaster   chan string
	missingMaster chan bool

	msgRecieveGlobal chan *types.NetworkMessage
	msgRecieveLocal  chan *types.NetworkMessage
	msgSendGlobal    chan *types.NetworkMessage
	msgSendLocal     chan *types.NetworkMessage

	netstatNewMsg chan *types.NetworkMessage
	netstatUpdate chan *types.NetworkMessage
	netstatTick   chan bool
}

func NewNetworkHub(resetCh chan bool,
	sendLocalCh, recieveLocalCh chan *types.NetworkMessage) *NetworkHub {
	nh := &NetworkHub{
		reset: resetCh,

		id: -1,

		networkStatus: types.NewNetworkMessage(),
		cm:            nil,

		foundMaster:   make(chan string),
		missingMaster: make(chan bool),

		msgRecieveGlobal: make(chan *types.NetworkMessage, 10),
		msgRecieveLocal:  recieveLocalCh,
		msgSendGlobal:    make(chan *types.NetworkMessage, 1),
		msgSendLocal:     sendLocalCh,

		netstatNewMsg: make(chan *types.NetworkMessage, 10),
		netstatUpdate: make(chan *types.NetworkMessage, 1),
		netstatTick:   make(chan bool),
	}
	nh.cm = newConnManager(nh.msgRecieveGlobal, nh.msgSendGlobal)
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
		case msgRecieve := <-nh.msgRecieveGlobal:
			nh.id = msgRecieve.Id
			types.DeepCopy(nh.networkStatus, msgRecieve)
			nh.msgSendLocal <- msgRecieve
		case msgUpdate := <-nh.msgRecieveLocal:
			nh.networkStatus = msgUpdate
			nh.msgSendGlobal <- msgUpdate
		}
	}

	go startUDPBroadcast(nh.reset)
	go newNetStatManager(nh.netstatNewMsg, nh.netstatUpdate, nh.netstatTick).
		run(nh.networkStatus, nh.reset)

	tick := time.Tick(SEND_INTERVAL * time.Millisecond)
	// Master loop
	for {
		select {
		case msgRecGlobal := <-nh.msgRecieveGlobal:
			nh.netstatNewMsg <- msgRecGlobal
		case msgRecLocal := <-nh.msgRecieveLocal:
			nh.netstatNewMsg <- msgRecLocal
		case <-tick:
			nh.netstatTick <- true
		case newNetstat := <-nh.netstatUpdate:
			nh.networkStatus = newNetstat
			nh.msgSendGlobal <- newNetstat
			nh.msgSendLocal <- newNetstat
		case <-nh.foundMaster:
			fmt.Println("\x1b[31;1m::: MULTIPLE MASTERS FOUND :::\x1b[0m")
			if len(nh.cm.conns) == 0 {
				fmt.Println("\x1b[31;1m::: I QUIT :::\x1b[0m")
				close(nh.reset)
				return
			} else {
				fmt.Println("\x1b[31;1m::: numConns", len(nh.cm.conns), ":::\x1b[0m")
			}
		}
	}
}
