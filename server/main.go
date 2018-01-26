package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 9981})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("Local: <%s> \n", listener.LocalAddr().String())
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {
			log.Printf("got 2 peers, start p2p test:%+v", peers)
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
			time.Sleep(time.Second * 1)
			log.Println("server exit.")
			return
		}
	}
}
