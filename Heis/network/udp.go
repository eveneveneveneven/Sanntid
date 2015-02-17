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

    ln *net.UDPConn
}

// Finds the local IP address of the machine
func getLocalAddres() string {
    baddr := &net.UDPAddr{
        Port: UDP_PORT,
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
        Port: UDP_PORT,
        IP: net.ParseIP("localhost"),
    }
    u.raddr = nil
    u.ln = nil

    return &u
}

// Find a Master on the network if there are any.
// Listens for 0.3 second, quits after first reading or none after 0.3 second.
// returns (ifFound, masterIP, error)
func (u *UDPHub) findMaster(report bool) (bool, string, error) {
    if u.ln == nil {
        // Create a listener which listens after broadcasts form potential Master on UDP_PORT
    	ln, err := net.ListenUDP("udp", u.laddr)
        if err != nil {
            return false, "", err
        }
        u.ln = ln
    }
    // Variable for continuing listening
    searching := true
    // Variable for storing the potential Master IP
    var masterIP string

    // A timeout function which ends the search after 0.3 second
    timeout := make(chan bool)
    go func() {
    	time.Sleep(300 * time.Millisecond)
    	timeout <- true
        searching = false
    }()

    // Listener function which listens after Master broadcast, if any.
    // Stores the first read it gets, sets this as Master, and quits searching.
    // The message read from Master is the Master IP, which is stored for later use.
    listener := make(chan bool)
    go func() {
        for searching {
            n, raddr, err := u.ln.ReadFromUDP(u.p)
            if err != nil {
                fmt.Printf("Some error %v, continuing listening\n", err)
                continue
            }
            listener <- true
            if report {
                masterIP = string(u.p[:n])
                u.raddr = raddr
            }
            break
        }
    }()

    // Quits if it finds a broadcast or 1 second has passed.
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
func (u *UDPHub) broadcastMaster(stop <-chan bool) {
    // Broadcast address
	baddr := &net.UDPAddr{
        Port: UDP_PORT,
        IP: net.IPv4bcast,
    }

    // Create a broadcast-connection on the UDP network
    conn, err := net.DialUDP("udp", nil, baddr)
    if err != nil {
        // If the program cant make a UDP connection, no need to run the program.
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
            // Brodcasting Master IP
        	fmt.Fprintf(conn, u.localAddr) // trick to send a message on the UDP network!
            if err != nil {
                fmt.Printf("Some error writing! %v\n", err)
            }

            time.Sleep(100 * time.Millisecond)
        }
    }
}

func (u *UDPHub) alertWhenMaster(alert chan<- bool) {
    for {
        found, _, err := u.findMaster(false)
        if err != nil {
            fmt.Printf("Some error %v, trying again\n", err)
            continue
        }
        if !found {
            fmt.Println("Master is dead!")
            alert <- true
            return
        }
    }
}