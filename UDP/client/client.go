package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"strings"
)

// var udpAddr *net.UDPAddr

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s :port", os.Args[0])
		os.Exit(1)
	}

	handleClient(os.Args[1])
	os.Exit(0)
}

func handleClient(service string) {
	wg := sync.WaitGroup{}

	// Setting up
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	// Write requests
	wg.Add(1)
	go func() {
		handleClientSendMessage(conn)
		wg.Done()
	}()

	// Read responses
	wg.Add(1)
	go func() {
		handleReceivedMessage(conn)
		wg.Done()
	}()

	wg.Wait()
}

func handleClientSendMessage(conn *net.UDPConn){
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadBytes('\n')
		checkError(err)
		messSTR := []byte(string(message) + "\n")

		// Write request
		_, err = conn.Write(messSTR)
		checkError(err)
	}
}

func handleReceivedMessage(conn *net.UDPConn) {
	for {
		var buf [512]byte
		n, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			continue
		}
		message := strings.TrimSpace(string(buf[0:n]))
		fmt.Println(message)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}