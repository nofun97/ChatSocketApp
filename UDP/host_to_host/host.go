package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	// "time"
)

var currAddress *net.UDPAddr

var endConnection = false


func handleConnection(service string) {
	wg := sync.WaitGroup{}
	// Binding process
	udpAddr, err := net.ResolveUDPAddr("udp", service)
	checkError(err)
	currAddress = udpAddr
	// Creating a connection, essentially connect()
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)


	wg.Add(1)
	go func() {
		for {
			handleReceivedMessage(conn)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		handleClientSendMessage(conn)
		wg.Done()
	}()

	// Print result
	wg.Wait()
}

func handleServer(port string) {
	wg := sync.WaitGroup{}
	udpAddr, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)

	// Listen()
	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	wg.Add(1)
	go func() {
		for {
			handleReceivedMessage(conn)
		}
		wg.Done()
		fmt.Println("Done")
	}()

	wg.Add(1)
	go func() {
		handleSendMessage(conn)
		wg.Done()
	}()

	wg.Wait()
}

func main() {

	if len(os.Args) == 3 {
		service := os.Args[1]
		if service == "server" {
			handleServer(os.Args[2])
		} else if service == "client" {
			handleConnection(os.Args[2])
		} else {
			log.Println("Either be a server or a client, gotta choose one my dude")
		}
	} else {
		log.Println("Wrong format my dude")
	}

	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleReceivedMessage(conn *net.UDPConn) {
	var buf [512]byte

	n, addr, err := conn.ReadFromUDP(buf[0:])
	currAddress = addr
	if err != nil {
		return
	}
	fmt.Println(n)
	message := string(buf[0:n])
	fmt.Println("From ", addr, ": ", message)
}

func handleSendMessage(conn *net.UDPConn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadBytes('\n')
		checkError(err)
		messSTR := string(message) + "\n"
		// Write request
		conn.WriteToUDP([]byte(messSTR), currAddress)
	}
}

func handleClientSendMessage(conn *net.UDPConn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadBytes('\n')
		checkError(err)
		messSTR := string(message) + "\n"
		// Write request
		conn.Write([]byte(messSTR))
	}
}