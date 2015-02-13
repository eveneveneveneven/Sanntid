package main

import (
    "fmt"
    "net"
)

type P struct {
    X, Y int64
}

func main() {
    p := make([]byte, 2048)
	recieve := false
	fmt.Println("UDP client start")
	// 129.241.187.136 - server IP
	raddr, err := net.ResolveUDPAddr("udp", "129.241.187.255:30000")
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
    conn, err := net.DialUDP("udp", nil, raddr)
    if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
    }
    defer conn.Close()
    fmt.Fprintf(conn, "Hi UDP Server, how are you doing?")
	if recieve {
		laddr := net.UDPAddr{
		    Port: 20011,
		    IP: net.ParseIP("localhost"),
		}
		ln, err := net.ListenUDP("udp", &laddr)
		if err != nil {
			fmt.Printf("Some error %v\n", err)
			return
		}
		n, newRaddr, err := ln.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error %v\n", err)
			return
		}
		fmt.Println("Recieved something!")
		fmt.Printf("We got back from addr %v :: %s\n", newRaddr, p[:n])
	}
}
