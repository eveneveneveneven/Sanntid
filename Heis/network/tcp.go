package network

import (
	"net"
	"fmt"
	"encoding/gob"
	"os"
	"time"
	"sync"
)

type TCPHub struct {
	id int // 0 equals master, else slave
	networkStatus string

	wg sync.WaitGroup

	masterIP string
	masterConn *net.TCPConn

	numSlaves int
	slaves []*slave
}

type slave struct {
	id int
	conn *net.TCPConn
	enc *gob.Encoder
	dec *gob.Decoder
}

// Init of new TCPHub variable
func newTCPHub() *TCPHub {
	var t TCPHub

	t.id = -1
	t.networkStatus = "xx"

	t.masterIP = ""
	t.masterConn = nil


	t.numSlaves = 0
	t.slaves = make([]*slave, MAX_ELEVATORS)

	return &t
}

func (t *TCPHub) removeDeadSlave(sl *slave) {
	fmt.Printf("%v\n%+v\n%v\n", t.networkStatus, t.slaves, t.numSlaves)
	t.networkStatus = t.networkStatus[:sl.id*2] + t.networkStatus[sl.id*2+2:]
	t.slaves = append(append(t.slaves[:sl.id-1], t.slaves[sl.id:]...), nil)
	t.numSlaves--
	fmt.Printf("%v\n%+v\n%v\n", t.networkStatus, t.slaves, t.numSlaves)
	for i := sl.id-1; i < t.numSlaves; i++ {
		t.slaves[i].id--
	}
	sl.conn.Close()
}

func (t *TCPHub) updateSlave(sl *slave, removeCh chan<- bool) {
	// Send the current network status
	sendMsg := &networkStatusMsg{
		Id: sl.id,
		Data: t.networkStatus,
	}
	if err := sl.enc.Encode(sendMsg); err != nil {
		fmt.Printf("1>Some error %v, removing slave from network\n", err)
		removeCh <- true
    	return
	}
	// Wait until all slaves have recieved the update
	t.wg.Done()
	t.wg.Wait()
	
	// Recieve elevator status from slave
	recMsg := &networkStatusMsg{}
	if err := sl.dec.Decode(recMsg); err != nil {
		fmt.Printf("2>Some error %v, removing slave from network\n", err)
		removeCh <- true
    	return
	}

	// Update network status
	index := recMsg.Id * 2
	t.networkStatus = t.networkStatus[:index] + recMsg.Data + t.networkStatus[index+2:]
}

func newSlave(id int, conn *net.TCPConn) *slave {
	s := &slave{
		id: id,
		conn: conn,
		enc: gob.NewEncoder(conn),
		dec: gob.NewDecoder(conn),
	}
	return s
}

func (t *TCPHub) handleSlaveConnection(conn *net.TCPConn) {
	fmt.Println("Got new slave connection")
	if t.numSlaves < MAX_ELEVATORS {
		t.numSlaves++
		slaveId := t.numSlaves
		sl := newSlave(slaveId, conn)
		t.slaves[slaveId-1] = sl
		t.networkStatus += "xx"
		timer := time.NewTimer(250 * time.Millisecond)
		removeCh := make(chan bool)

		for {
			select {
			case <-timer.C:
				go t.updateSlave(sl, removeCh)
			case <-removeCh:
				t.removeDeadSlave(sl)
				return
			}
			timer.Reset(250 * time.Millisecond)
		}
	} else {
		fmt.Println("Refuse connection, have max number of slaves.")
		conn.Close()
	}
}

func (t *TCPHub) startMasterServer(stop <-chan bool) {
	laddr := &net.TCPAddr{
		Port: TCP_PORT,
		IP: net.ParseIP("localhost"),
	}

	ln, err := net.ListenTCP("tcp", laddr)
    if err != nil {
        fmt.Printf("Some error %v, quitting program\n", err)
        os.Exit(1)
    }
    defer ln.Close()

    go func() {
		for {
			t.wg.Wait()
			t.wg.Add(t.numSlaves)
		}
	}()

    for {
    	conn, err := ln.AcceptTCP()
	    if err != nil {
	        fmt.Printf("Some error %v, continue listening\n", err)
	        continue
	    }
	    go t.handleSlaveConnection(conn)
    }
}

func (t *TCPHub) startSlaveClient() {
	decoder := gob.NewDecoder(t.masterConn)
	encoder := gob.NewEncoder(t.masterConn)
	for {
		recMsg := &networkStatusMsg{}
		if err := decoder.Decode(recMsg); err != nil {
			fmt.Printf("Some error %v\n", err)
	    	return
		}
		t.id = recMsg.Id
		t.networkStatus = recMsg.Data
		fmt.Println(recMsg.Data)

		sendMsg := &networkStatusMsg{
			Id: t.id,
			Data: "OK",
		}
		if err := encoder.Encode(sendMsg); err != nil {
			fmt.Printf("Some error %v\n", err)
	    	return
		}
	}
}

// Asks found Master if it can connect to the network.
// Connects itself to the network if approved,
// else shuts program off (not needed/allowed).
// Returns (isAllowed, error)
func (t *TCPHub) requestConnToNetwork(masterIP string) (bool, error) {
	t.masterIP = masterIP
	raddr := &net.TCPAddr{
		Port: TCP_PORT,
		IP: net.ParseIP(t.masterIP),
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
	    return false, err
	}

	t.masterConn = conn

	fmt.Println("Accepted connection to the network, begin transmition")
	return true, nil
}