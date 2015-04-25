package network

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	UDP_PORT = 20011
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
	fmt.Println("\x1b[34;1m::: Start UDP Listener :::\x1b[0m")

	laddr := &net.UDPAddr{
		Port: UDP_PORT,
		IP:   net.ParseIP("localhost"),
	}
	ln, err := net.ListenUDP("udp", laddr)
	if err != nil {
		fmt.Printf("\t\x1b[31;1mError\x1b[0m |startUDPListener| [%v], exiting program\n", err)
		os.Exit(1)
	}
	defer ln.Close()
	myIp := getLocalAddress()
	p := make([]byte, 1024)
	for {
		ln.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		n, _, err := ln.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("\t\x1b[31;1mError\x1b[0m |startUDPListener| [%v], continuing listening\n", err)
			masterMissing <- true
			continue
		}
		ipFound := string(p[:n])
		if ipFound != myIp {
			foundMaster <- ipFound
		}
	}
}

func startUDPBroadcast(resetCh chan bool) {
	fmt.Println("\x1b[34;1m::: Start UDP Broadcaster :::\x1b[0m")

	// Broadcast address
	baddr := &net.UDPAddr{
		Port: UDP_PORT,
		IP:   net.IPv4bcast,
	}
	conn, err := net.DialUDP("udp", nil, baddr)
	if err != nil {
		fmt.Printf("\t\x1b[31;1mError\x1b[0m |startUDPBroadcast| [%v], exiting program\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	ip := getLocalAddress()
	tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		case _, ok := <-resetCh:
			if !ok {
				return
			}
		case <-tick:
			fmt.Fprintf(conn, ip) // trick to send a message on the UDP network!
			if err != nil {
				fmt.Printf("\t\x1b[31;1mError\x1b[0m |startUDPBroadcast| [%v]\n", err)
			}
		}
	}
}
