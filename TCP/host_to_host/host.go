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

const (
	mode = "tcp"
)

var endConnection = false

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

func handleServer(port string) {
	wg := sync.WaitGroup{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	checkError(err)

	// Listen()
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	conn, err := listener.Accept()
	checkError(err)

	wg.Add(1)
	go func(){
		handleReceivedMessage(conn)
		wg.Done()
	}()

	wg.Add(1)
	go func(){
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

func handleReceivedMessage(conn net.Conn){
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if scanner.Text() == "" || scanner.Text() == "\n"{
			continue
		}
		fmt.Println("From ", conn.RemoteAddr().String(),": ", scanner.Text())
	}
}

func handleSendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadBytes('\n')
		checkError(err)
		messSTR := string(message) + "\n"
		conn.Write([]byte(messSTR))
	}
}