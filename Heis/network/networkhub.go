package network

import (
	"errors"
	"fmt"
	"os"
)

type Hub struct {
	master bool

	udp *UDPHub
	tcp *TCPHub

	missingMaster chan bool
	stop chan bool
}

func NewHub() *Hub {
	var h Hub

	h.master = false

	h.udp = newUDPHub()
	h.tcp = newTCPHub()

	h.missingMaster = make(chan bool)
	h.stop = make(chan bool, 2)

	return &h
}

func (h *Hub) Run() {
	err := h.resolveMasterNetwork()
	if err != nil {
		fmt.Printf("Some error %v, exit program\n", err)
		os.Exit(1)
	}

	for {
		select {
		case <-h.missingMaster:
			fmt.Println("Master is dead")
			h.tcp.id -= 1
			if h.tcp.id == 0 {
				h.becomeMaster()
			} else {
				go h.udp.alertWhenMaster(h.missingMaster)
			}
		case <-h.stop:
		}
	}
}

func (h *Hub) becomeMaster() {
	fmt.Println("Becoming Master")
	h.master = true
	h.tcp.id = 0
	go h.udp.broadcastMaster(h.stop)
	go h.tcp.startMasterServer(h.stop)
}

func (h *Hub) resolveMasterNetwork() error {
	found, masterIP, err := h.udp.findMaster(true);
	if err != nil {
		return err
	}

	if found {
		ok, err := h.tcp.requestConnToNetwork(masterIP)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("Refused connection to network")
		}
		fmt.Println("I am a slave...")
		go h.udp.alertWhenMaster(h.missingMaster)
		go h.tcp.startSlaveClient()

		return nil
	} else {
		h.becomeMaster()

		return nil
	}
}