package main

import (
    "fmt"
    "net"
    "bufio"
)

type P struct {
    X, Y int64
}

func main() {
    p := make([]byte, 2048)
    conn, err := net.Dial("udp", "localhost:8080")
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    defer conn.Close()
    fmt.Fprintf(conn, "Hi UDP Server, how are you doing?")
    _, err = bufio.NewReader(conn).Read(p)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
        return
    }
    fmt.Printf("%s\n", p)
}