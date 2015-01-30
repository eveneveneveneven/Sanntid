package main

import (
    "fmt"
    "log"
    "net"
)


func main() {
    fmt.Println("start client");
    conn, err := net.Dial("tcp", "129.241.187.136:33546")
    if err != nil {
        log.Fatal("Connection error", err)
    }
    defer conn.Close()
   	
   	p := make([]byte, 2048)
    conn.Write([]byte("Connect to: 129.241.187.159:20011\x00"))
    n, _ := conn.Read(p)
    fmt.Printf("We got back msg :: %s\n", p[:n])
    fmt.Println("done");
}
