package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	direction := os.Args[1]
	// 双向打洞
	if os.Args[1] == "bidirection" {
		direction = "<->"
	} else {
		//单向打洞
		direction = "->"
	}
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 9981})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Printf("本地地址: <%s> \n", listener.LocalAddr().String())
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
			if direction == "->" {
				log.Printf("进行单向打洞,建立 %s --> %s 的连接\n", peers[0].String(), peers[1].String())
				listener.WriteToUDP([]byte("->:"+peers[1].String()), &peers[0])
				listener.WriteToUDP([]byte("<-:"+peers[0].String()), &peers[1])
			} else {
				log.Printf("进行双向向打洞,建立 %s <--> %s 的连接\n", peers[0].String(), peers[1].String())
				listener.WriteToUDP([]byte("<->:"+peers[1].String()), &peers[0])
				listener.WriteToUDP([]byte("<->:"+peers[0].String()), &peers[1])
			}
			time.Sleep(time.Second * 8)
			log.Println("中转服务器退出,仍不影响peers间通信")
			return
		}
	}
}
