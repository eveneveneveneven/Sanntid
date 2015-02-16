package network

import (
	"net"
	"fmt"
	"encoding/gob"
	"os"
)

type TCPHub struct {
	masterIP string
	masterConn *net.TCPConn

	numConns int
	conns []*net.TCPConn
}

// Init of new TCPHub variable
func newTCPHub() *TCPHub {
	var t TCPHub

	t.masterIP = ""
	t.masterConn = nil

	t.numConns = 0
	t.conns = make([]*net.TCPConn, MAX_ELEVATORS)

	return &t
}

func (t *TCPHub) handleConneciton(conn *net.TCPConn) {
	decoder := gob.NewDecoder(conn)
	recMsg := &networkMessage{}
	if err := decoder.Decode(recMsg); err != nil {
		fmt.Printf("Some error %v\n", err)
	    return 
	}

	encoder := gob.NewEncoder(conn)
	sendMsg := NM_REQ_ACCE
	if err := encoder.Encode(sendMsg); err != nil {
		fmt.Printf("Some error %v\n", err)
	    return 
	}

	t.conns[t.numConns] = conn
	t.numConns += 1
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

    listening := true

    go func(){
    	<-stop
    	listening = false
    	ln.Close()
    }()

    for listening {
		fmt.Println("Listening for connection")
    	conn, err := ln.AcceptTCP()
	    if err != nil {
	        fmt.Printf("Some error %v, continue listening\n", err)
	        continue
	    }
	    fmt.Println("Got connection!")
	    go t.handleConneciton(conn)
    }
}

// Asks found Master if it can connect to the network.
// Connects itself to the network if approved,
// else shuts program off (not needed/allowed).
// Returns (isAllowed, ID, error)
func (t *TCPHub) requestConnToNetwork(masterIP string) (bool, int, error) {
	t.masterIP = masterIP
	raddr := &net.TCPAddr{
		Port: TCP_PORT,
		IP: net.ParseIP(t.masterIP),
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	    return false, -1, err
	}

	encoder := gob.NewEncoder(conn)
	sendMsg := NM_REQ_CONN
	if err := encoder.Encode(sendMsg); err != nil {
		fmt.Printf("Some error %v\n", err)
	    return false, -1, err
	}

	decoder := gob.NewDecoder(conn)
	recMsg  := &networkMessage{}
	if err := decoder.Decode(recMsg); err != nil {
		fmt.Printf("Some error %v\n", err)
	    return false, -1, err
	}

	if recMsg.Bool {
		fmt.Println("Accepted connection to the network, begin transmition")
		t.masterConn = conn
		return true, 1, nil
	} else {
		fmt.Println("Denied connection to the network, quit program")
		return false, -1, nil
	}
}