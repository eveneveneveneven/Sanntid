package network

import (
	"errors"
	"fmt"
	"os"
)

type Hub struct {
	master bool

	udp *UDPHub
	cm  *connManager

	missingMaster chan bool
	stop          chan bool
}

func NewHub() *Hub {
	var h Hub

	h.master = false

	h.udp = newUDPHub()
	h.cm = newConnManager()

	h.missingMaster = make(chan bool)

	return &h
}

func (h *Hub) Run() {
	go h.cm.run()
	newMaster, err := h.resolveMasterNetwork()
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

func (h *Hub) resolveMasterNetwork() (bool, error) {
	found, masterIP, err := h.udp.findMaster(true)
	if err != nil {
		return err
	}

	if found {
		ok, id, err := h.tcp.connectToNetwork(masterIP)
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
