package network

import (
	"time"
	"net"
    "fmt"
    "os"
)

const PORT = 20011

type UDPHub struct {
	master bool
	p []byte

	laddr *net.UDPAddr
	raddr *net.UDPAddr
}

func NewUDPHub() *UDPHub {
	var u UDPHub

	u.master = false
	u.p = make([]byte, 1024)

	u.laddr = &net.UDPAddr{
        Port: PORT,
        IP: net.ParseIP("localhost"),
    }
    u.raddr = nil

    return &u
}

func (u *UDPHub) FindMaster() (bool, error) {
	ln, err := net.ListenUDP("udp", u.laddr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return false, err
    }
    defer ln.Close()

    quitSearching := false

    timeout := make(chan bool, 1)
    go func() {
    	time.Sleep(1 * time.Second)
    	timeout <- true
        quitSearching = true
    }()

    listener := make(chan bool, 1)
    go func() {
        for quitSearching {
            _, raddr, err := ln.ReadFromUDP(u.p)
            if err != nil {
                fmt.Printf("Somer error %v\n, continuing listening", err)
                continue
            }  
            u.raddr = raddr
            return
        }
    }()

    select {
    case <-listener:
        return true, nil

    case <-timeout:
    	u.master = true
    	return false, nil
    }
}

func (u *UDPHub) BroadcastMaster() {
	baddr := &net.UDPAddr{
        Port: PORT,
        IP: net.ParseIP("broadcast"),
    }

    socket, err := net.DialUDP("udp", nil, baddr)
    if err != nil {
        fmt.Println("BroadcastMaster could not dial up. Exiting")
        os.Exit(1)
    }

    for {
    	msg := "Master exists"
    	socket.WriteToUDP([]byte(msg), baddr)
        time.Sleep(100 * time.Millisecond)
    }
}