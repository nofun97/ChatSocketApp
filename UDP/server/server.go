package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	MaxClients = 5
)

var clients = struct {
	sync.RWMutex
	m map[*net.UDPAddr]int
}{m: make(map[*net.UDPAddr]int)}

type messages struct {
	conn    *net.UDPConn
	addr    *net.UDPAddr
	message [512]byte
}

func main() {
	wg := sync.WaitGroup{}
	service := ":1337"
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for i := 0; i < MaxClients; i++ {
		wg.Add(1)
		go func() {
			for {
				handleClient(conn)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte

	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	conn.WriteToUDP([]byte("You're connected"), addr)

	clients.Lock()
	if _, ok := clients.m[addr]; !ok {
		clients.m[addr] = 1
	}
	clients.Unlock()

	fmt.Println(string(buf[:n]))
	newMessage := &messages{conn: conn, addr: addr, message: buf}
	sendMessages(newMessage)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func sendMessages(message *messages) {
	clients.Lock()
	defer clients.Unlock()
	for k, v := range clients.m {
		if k != message.addr && v == 1 {
			message.conn.WriteToUDP(message.message[0:], message.addr)
		}
	}
}
