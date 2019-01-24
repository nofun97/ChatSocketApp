package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

type messages struct {
	conn    net.Conn
	message []byte
}

var clients = struct {
	sync.RWMutex
	m map[net.Conn]int
}{m: make(map[net.Conn]int)}

func main() {
	service := ":1337"
	// Binding address
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	// Listen()
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// messageQueue := make(chan messages, 20)
	connectionQueue := make(chan net.Conn, 5)
	var conn net.Conn

	go func() {
		for {
			select {
			case currentConn := <-connectionQueue:
				go handleConnection(currentConn)
			}
		}
	}()

	for {
		conn, err = listener.Accept()

		checkError(err)
		clients.Lock()
		if _, ok := clients.m[conn]; !ok {
			clients.m[conn] = 1
		}
		clients.Unlock()
		connectionQueue <- conn
		fmt.Println(len(connectionQueue))
	}

	os.Exit(0)
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Bytes()
		fmt.Println(string(text))
		message := &messages{conn: conn, message: text}
		handleMessage(message)
	}
}

func handleMessage(message *messages) {
	if string(message.message) == "/quit" {
		clients.Lock()
		if _, ok := clients.m[message.conn]; ok {
			message.conn.Close()
			clients.m[message.conn] = 0
		}
		clients.Unlock()
	}

	sendMessages(message)
}

func sendMessages(message *messages) {
	clients.Lock()
	defer clients.Unlock()
	for k, v := range clients.m {
		if k != message.conn && v == 1 {
			text := string(message.message) + "\n"
			k.Write([]byte(text))
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
