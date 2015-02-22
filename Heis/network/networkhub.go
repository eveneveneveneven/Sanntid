package network

import (
	"errors"
	"fmt"
	"os"
)

type Hub struct {
	master bool
	id int // 0 equals master, else slave

	numConns int

	udp *UDPHub
	tcp *TCPHub

	missingMaster chan bool
	stop chan bool
}

func NewHub() *Hub {
	var h Hub

	h.id       = -1
	h.master   = false
	h.numConns = 0

	h.udp = newUDPHub()
	h.tcp = newTCPHub()

	h.missingMaster = make(chan bool)
	h.stop = make(chan bool, 2)

	return &h
}

func (h *Hub) Run() {
	newMaster, err := h.resolveMasterNetwork()
	if err != nil {
		fmt.Printf("Some error %v, exit program\n", err)
		os.Exit(1)
	}

	if newMaster {
		fmt.Println("I am Master!")
		h.becomeMaster()
	} else {
		fmt.Println("I am a slave...")
		go h.udp.alertWhenMaster(h.missingMaster)
	}

	for {
		select {
		case <-h.missingMaster:
			if h.id == 1 {
				h.becomeMaster()
			} else {
				h.id -= 1
				h.udp.alertWhenMaster(h.missingMaster)
			}
		case <-h.stop:
		}
	}
}

func (h *Hub) becomeMaster() {
	fmt.Println("Becoming Master")
	h.master = true
	go h.udp.broadcastMaster(h.stop)
	go h.tcp.startMasterServer(h.stop)
}

func (h *Hub) resolveMasterNetwork() (bool, error) {
	found, masterIP, err := h.udp.findMaster(true);
	if err != nil {
		return false, err
	}

	if found {
		ok, id, err := h.tcp.requestConnToNetwork(masterIP)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, errors.New("Refused connection to network")
		}

		fmt.Printf("Got ID %v\n", id)
		h.id = id
		return false, nil
	} else {
		h.id = 0
		return true, nil
	}
}