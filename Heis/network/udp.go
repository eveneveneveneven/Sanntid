package network

import (
	"time"
	"net"
    "fmt"
    "os"
    "strings"
)

const PORT = 20011

type UDPHub struct {
	master bool
	p []byte

    localAddr string

	laddr *net.UDPAddr
	raddr *net.UDPAddr
}

func NewUDPHub() *UDPHub {
	var u UDPHub

	u.master = false
	u.p = make([]byte, 1024)

    baddr := &net.UDPAddr{
        Port: PORT,
        IP: net.IPv4bcast,
    }
	tempConn, _ := net.DialUDP("udp4", nil, baddr)
    defer tempConn.Close()
    tempAddr := tempConn.LocalAddr()
    u.localAddr = strings.Split(tempAddr.String(), ":")[0] // only want ip

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

    searching := true

    timeout := make(chan bool, 1)
    go func() {
    	time.Sleep(1 * time.Second)
    	timeout <- true
        searching = false
    }()

    listener := make(chan bool, 1)
    go func() {
        defer ln.Close()
        for searching {
            fmt.Println("Searching for master")
            n, raddr, err := ln.ReadFromUDP(u.p)
            if err != nil {
                fmt.Printf("Somer error %v, continuing listening\n", err)
                continue
            }
            fmt.Printf("Got %s\n", u.p[:n])
            listener <- true
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

func (u *UDPHub) BroadcastMaster(stop chan bool) {
	baddr := &net.UDPAddr{
        Port: PORT,
        IP: net.IPv4bcast,
    }

    conn, err := net.DialUDP("udp", nil, baddr)
    if err != nil {
        fmt.Println("BroadcastMaster could not dial up. Exiting")
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Broadcasting")
    for {
        select {
        case <-stop:
            break
        default:
        	fmt.Fprintf(conn, "Master=" + u.localAddr) // trick!
            if err != nil {
                fmt.Printf("Some error writing! %v\n", err)
            }

            time.Sleep(100 * time.Millisecond)
        }
    }
}