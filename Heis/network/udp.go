package network

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Finds the local IP address of the machine
func getLocalAddress() string {
	baddr := &net.UDPAddr{
		Port: UDP_PORT,
		IP:   net.IPv4bcast,
	}
	tempConn, _ := net.DialUDP("udp4", nil, baddr)
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr()
	return strings.Split(tempAddr.String(), ":")[0] // only want ip
}

func startUDPListener(foundMaster chan string, masterMissing chan bool) {
	laddr := &net.UDPAddr{
		Port: UDP_PORT,
		IP:   net.ParseIP("localhost"),
	}
	ln, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Println("Some error %v, exiting program\n")
		os.Exit(1)
	}
	defer ln.Close()

	p := make([]byte, 1024)
	for {
		ln.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, _, err := ln.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error %v, continuing listening\n", err)
			masterMissing <- true
			continue
		}
		foundMaster <- string(p[:n])
	}
}

func startUDPBroadcast() {
	// Broadcast address
	baddr := &net.UDPAddr{
		Port: UDP_PORT,
		IP:   net.IPv4bcast,
	}
	conn, err := net.DialUDP("udp", nil, baddr)
	if err != nil {
		fmt.Printf("Some error %v, exiting program\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	ip := getLocalAddress()
	for {
		fmt.Fprintf(conn, ip) // trick to send a message on the UDP network!
		if err != nil {
			fmt.Printf("Some error writing! %v\n", err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
