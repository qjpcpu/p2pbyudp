package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	sip := net.ParseIP("207.148.70.129")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: sip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	if _, err = conn.Write([]byte("hello, I'm new peer " + srcAddr.String())); err != nil {
		log.Panic(err)
	}
	fmt.Printf("<%s>\n", conn.RemoteAddr())
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
	anotherPeer := parseAddr(string(data[:n]))
	log.Printf("get another peer:%s", anotherPeer.String())

	conn, err = net.DialUDP("udp", srcAddr, &anotherPeer)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte("handshake from:" + srcAddr.String())); err != nil {
		log.Println("send handshake:", err)
	} else {
		log.Println("send ok")
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)
			if _, err = conn.Write([]byte("from " + srcAddr.String())); err != nil {
				log.Println("send msg ", err)
			} else {
				log.Println("send ok")
			}
		}
	}()
	for {
		data = make([]byte, 1024)
		n, remoteAddr, err = conn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("<%s> %s\n", remoteAddr, data[:n])
	}
}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}
