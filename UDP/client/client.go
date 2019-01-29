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
	wg := sync.WaitGroup{}
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	_, err = conn.Write([]byte("anything"))
	checkError(err)

	// Write requests
	wg.Add(1)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			message, err := reader.ReadBytes('\n')
			checkError(err)
			messSTR := string(message) + "\n"
			// Write request
			handleWrite(conn, []byte(messSTR))
		}
		wg.Done()
	}()

	// Read responses
	wg.Add(1)
	go func() {
		for {
			var buf [512]byte
			n, _, err := conn.ReadFromUDP(buf[0:])
			if err != nil {
				continue
			}
			message := strings.TrimSpace(string(buf[0:n]))
			fmt.Println(message)
		}
		wg.Done()
	}()

	wg.Wait()
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func handleWrite(conn *net.UDPConn, message []byte) {
	_, err := conn.Write(message)
	checkError(err)
}
