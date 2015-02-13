package network

import (
	"time"
	"net"
    "fmt"
    "os"
    "strings"
)

type UDPHub struct {
	master bool
	p []byte

    localAddr string

	laddr *net.UDPAddr
	raddr *net.UDPAddr
}

func getLocalAddres() string {
    baddr := &net.UDPAddr{
        Port: PORT,
        IP: net.IPv4bcast,
    }
    tempConn, _ := net.DialUDP("udp4", nil, baddr)
    defer tempConn.Close()
    tempAddr := tempConn.LocalAddr()
    return strings.Split(tempAddr.String(), ":")[0] // only want ip
}

func newUDPHub() *UDPHub {
	var u UDPHub

	u.master = false
	u.p = make([]byte, 1024)

    u.localAddr = getLocalAddres()

    u.laddr = &net.UDPAddr{
        Port: PORT,
        IP: net.ParseIP("localhost"),
    }
    u.raddr = nil

    return &u
}

func (u *UDPHub) findMaster() (bool, string, error) {
	ln, err := net.ListenUDP("udp", u.laddr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return false, "", err
    }
    defer ln.Close()

    searching := true
    var masterIP string

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
            masterIP = string(u.p[:n])
            listener <- true
            u.raddr = raddr
            break
        }
    }()

    select {
    case <-listener:
        return true, masterIP, nil

    case <-timeout:
    	u.master = true
    	return false, "", nil
    }
}

func (u *UDPHub) broadcastMaster(stop chan bool) {
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
        	fmt.Fprintf(conn, u.localAddr) // trick! broadcasting master ip
            if err != nil {
                fmt.Printf("Some error writing! %v\n", err)
            }

            time.Sleep(100 * time.Millisecond)
        }
    }
}