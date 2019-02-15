package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var clients = make(map[net.Conn]int)

func main() {
	if len(os.Args) != 1 {
		log.Printf("Usage: %s :port ", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	handleServer(service)

	os.Exit(0)
}

func handleServer(port string) {
	// Binding address
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	checkError(err)

	// Listen()
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	connectionQueue := make(chan net.Conn, 5)

	go func() {
		for {
			select {
			case currentConn := <-connectionQueue:
				go handleConnection(currentConn)
			}
		}
	}()

	for {
		conn, err := listener.Accept()

		checkError(err)
		if _, ok := clients[conn]; !ok {
			clients[conn] = 1
		}
		connectionQueue <- conn
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Bytes()
		fmt.Println(string(text))
		if len(text) <= 0 {
			continue
		}

		for k := range clients {
			if k != conn {
				text := string(text) + "\n"
				k.Write([]byte(text))
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
