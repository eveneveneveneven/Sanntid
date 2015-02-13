package network

import (
	"time"
	"net"
    "fmt"
    "os"
    "strings"
)

type UDPHub struct {
	p []byte

    localAddr string

	laddr *net.UDPAddr
	raddr *net.UDPAddr
}

// Finds the local IP address of the machine
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

// Init of a new UDPHub variable
func newUDPHub() *UDPHub {
	var u UDPHub

	u.p = make([]byte, 1024)
    u.localAddr = getLocalAddres()

    u.laddr = &net.UDPAddr{
        Port: PORT,
        IP: net.ParseIP("localhost"),
    }
    u.raddr = nil

    return &u
}

// Find a Master on the network if there are any.
// Listens for 1 second, quits after first reading or none after 1 second.
// returns (ifFound, masterIP, error)
func (u *UDPHub) findMaster() (bool, string, error) {
    // Create a listener which listens after broadcasts form potential Master on PORT
	ln, err := net.ListenUDP("udp", u.laddr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return false, "", err
    }
    defer ln.Close()

    // Variable for continuing listening
    searching := true
    // Variable for storing the potential Master IP
    var masterIP string

    // A timeout function which ends the search after 1 second
    timeout := make(chan bool, 1)
    go func() {
    	time.Sleep(1 * time.Second)
    	timeout <- true
        searching = false
    }()

    // Listener function which listens after Master broadcast, if any.
    // Stores the first read it gets, sets this as Master, and quits searching.
    // The message read from Master is its Master IP, which is stored for later.
    listener := make(chan bool, 1)
    go func() {
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

    // Quits if it finds a broadcast, or 1 second has passed.
    select {
    case <-listener:
        return true, masterIP, nil

    case <-timeout:
    	return false, "", nil
    }
}

// Functions which serves as the broadcaster for the UDPHub.
// Is meant to be runned as a go-routine by the overall Hub struct
// if Master is not found on the network in the first place.
func (u *UDPHub) broadcastMaster(stop chan bool) {
	baddr := &net.UDPAddr{
        Port: PORT,
        IP: net.IPv4bcast,
    }

    // Create a broadcast-connection on the UDP network
    conn, err := net.DialUDP("udp", nil, baddr)
    if err != nil {
        fmt.Println("BroadcastMaster could not dial up. Exiting")
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Broadcasting as Master")
    // Keep broadcasting Masters IP 10 times a second on the UDP network
    // until told to stop by the stop-channel.
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