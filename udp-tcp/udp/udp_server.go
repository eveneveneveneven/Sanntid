package main

import (
    "fmt"
    "net"
	"time"
)

type P struct {
    X, Y int64
}

func handleConnection(conn *net.UDPConn, addr *net.UDPAddr) {
    _, err := conn.WriteToUDP([]byte("From server: Hello I got your msg"), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v\n", err)
    }
}

func main() {
    fmt.Println("start udp_server");
    p := make([]byte, 2048)
    laddr := net.UDPAddr{
        Port: 30000,
        IP: net.ParseIP("localhost"),
    }
    ln, err := net.ListenUDP("udp", &laddr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    for {
        n, raddr, err := ln.ReadFromUDP(p)
        if err != nil {
            fmt.Printf("Somer error %v\n", err)
            continue
        }
        fmt.Printf("%v :: Message from %v = %s\n", time.Now(), raddr, p[:n])
        go handleConnection(ln, raddr)
    }
}
