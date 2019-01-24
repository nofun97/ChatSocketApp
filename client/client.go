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
	wg := sync.WaitGroup{}
	reader := bufio.NewReader(os.Stdin)
	if len(os.Args) != 2 {
		log.Printf("Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	// Binding process
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	// Creating a connection, essentially connect()
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	scanner := bufio.NewScanner(conn)
	wg.Add(1)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			message, err := reader.ReadBytes('\n')
			checkError(err)
			messSTR := string(message) + "\n"
			// Write request
			handleWrite(conn, []byte(messSTR))
		}
		wg.Done()
	}()

	// Print result
	wg.Wait()
}

func handleWrite(conn *net.TCPConn, message []byte) {
	if string(message) == "/quit" {
		endConnection = true
	}
	_, err := conn.Write(message)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
