package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
	mode = "tcp"
)

var endConnection = false

func main() {
	if len(os.Args) != 2 {
		log.Printf("Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}

	handleConnection(os.Args[1])
}

func handleConnection(service string) {
	wg := sync.WaitGroup{}
	// Binding process
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	// Creating a connection, essentially connect()
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	wg.Add(1)
	go func() {
		handleReceivedMessage(conn)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		handleSendMessage(conn)
		wg.Done()
	}()

	// Print result
	wg.Wait()
}

func handleReceivedMessage(conn net.Conn){
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if scanner.Text() == "" || scanner.Text() == "\n"{
			continue
		}
		fmt.Println("From ", conn.RemoteAddr().String(),": ", scanner.Text())
	}
}

func handleSendMessage(conn *net.TCPConn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadBytes('\n')
		checkError(err)
		messSTR := string(message) + "\n"
		// Write request
		conn.Write([]byte(messSTR))
	}
}

func checkError(err error) {
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
