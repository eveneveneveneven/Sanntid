package main

import (
    "fmt"
    "net"
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
        Port: 8080,
        IP: net.ParseIP("localhost"),
    }
    ln, err := net.ListenUDP("udp", &laddr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    for {
        n, raddr, err := ln.ReadFromUDP(p)
        fmt.Printf("Read a message from %v %s\n", raddr, p[:n])
        if err != nil {
            fmt.Printf("Somer error %v\n", err)
            continue
        }
        go handleConnection(ln, raddr)
    }
}