package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var clients = make(map[string]*net.UDPAddr)

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s :port ", os.Args[0])
		os.Exit(1)
	}
	handleServer(os.Args[1])
}

func handleServer(service string) {
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for{
		handleClient(conn)
	}
}


func handleClient(conn *net.UDPConn) {
	var buf [512]byte

	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	// Updates the resolved UDP address
	clients[addr.String()] = addr

	// Creating the message and including from where the message came from
	message := []byte("From " + addr.String() + ": " + string(buf[0:n]) + "\n")
	for currAddr, realAddr := range clients {
		// To avoid sending the message to the original sender
		if currAddr != addr.String() {
			conn.WriteToUDP(message, realAddr)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
