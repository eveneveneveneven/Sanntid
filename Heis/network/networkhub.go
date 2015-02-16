package network

import (
	"errors"
	//"fmt"
)

type Hub struct {
	master bool
	id int // 0 equals master, else slave

	numConns int

	udp *UDPHub
	tcp *TCPHub

}

func NewHub() *Hub {
	var h Hub

	h.id       = -1
	h.master   = false
	h.numConns = 0

	h.udp = newUDPHub()
	h.tcp = newTCPHub()

	return &h
}

func (h *Hub) ResolveMasterNetwork(stop chan bool) (bool, error) {
	found, masterIP, err := h.udp.findMaster();
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

		h.master = false
		h.id = id

		return false, nil
	} else {
		h.master = true
		h.id = 0

		go h.udp.broadcastMaster(stop)
		go h.tcp.startMasterServer(stop)

		return true, nil
	}
}