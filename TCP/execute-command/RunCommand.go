package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	// "time"
)

const (
	mode = "tcp"
)

var endConnection = false
var messageToSend = make(chan []byte)

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
		for {
			handleCommandRun(conn)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func(){
		handleWriteOutput(conn)
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

func handleCommandRun(conn net.Conn) {

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		command := []byte(scanner.Text())
		if len(command) <= 0 {
			continue
		}
		fmt.Println("Command received: ", string(command))
		commLen := len(strings.Split(string(command), " "))
		var out []byte
		var err error
		if commLen < 1 {
			log.Println("No command received")
			return
		}
		out, err = exec.Command("bash", "-c", string(command)).Output()

		if err != nil {
			log.Println(err)
			messageToSend <- []byte(err.Error())
		} else {
			// fmt.Println(string(out))
			messageToSend <- []byte(out)
		}
	}

}

func handleWriteOutput(conn net.Conn) {
	for {
		select {
		case out := <-messageToSend :
			handleWrite(conn, out)
			fmt.Println("Message sent")
		}
	}
}

func handleWrite(conn net.Conn, message []byte) {
	if string(message) == "/quit" {
		endConnection = true
	}
	_, err := conn.Write(message)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
	}
}

func handleReceivedMessage(conn net.Conn){
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func handleSendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadBytes('\n')
		checkError(err)
		messSTR := string(message) + "\n"
		// Write request
		handleWrite(conn, []byte(messSTR))
	}
}