package main

import( 
	. "net"
	. "fmt"

	"time"
)


var err error
var chanCon chan *UDPConn

func ListenToNetwork(chanCon chan *UDPConn){

	Println("Start UDP server")
	
	
	udpAddr, err := ResolveUDPAddr("udp4", ":20011") //resolving
	if err != nil{
		println("ERROR: Resolving error")
	}

	conn, err := ListenUDP("udp", udpAddr) //initiating listening
	if err != nil{
		println("ERROR: Listening error")
	}
   chanCon <- conn
	data := make([]byte,1024)
	for{
		_, addr, err := conn.ReadFromUDP(data) //Kan bruke: if addr not [egen i.p]
		if err != nil{
			println("ERROR: while reading")
		}
		Println("Recieved from: ", addr,"\nMessage: ",string(data))
	}	
}

func SendToNetwork(chanCon chan *UDPConn){
   sendAddr, err := ResolveUDPAddr("udp4","129.241.187.255:20011") //Spesifiserer adresse
	//connection,err := DialUDP("udp",nil, sendAddr) //setter opp "socket" for sending
	if err != nil {
		println("ERROR while resolving UDP addr")
	}
	connection := <- chanCon
	testmsg := []byte("testing")

	if connection ==  nil{
		println("ERROR, connection = nil")
	}
	for{
		connection.WriteToUDP(testmsg, sendAddr)
		time.Sleep(1*time.Second)
	}	
}

func main(){

   chanCon := make(chan *UDPConn, 1)
	go ListenToNetwork(chanCon)
   go SendToNetwork(chanCon)
	for{
	   time.Sleep(5*time.Second)
	   Println("I am ALIVE")
	}
}

