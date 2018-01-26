package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var tag string

const HAND_SHAKE_MSG = "我是打洞消息"

func main() {
	// 当前进程标记字符串,便于显示
	tag = os.Args[1]
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 9982} // 注意端口必须固定
	dstAddr := &net.UDPAddr{IP: net.ParseIP("207.148.70.129"), Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	if _, err = conn.Write([]byte("hello, I'm new peer:" + tag)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer, direction := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, anotherPeer.String())

	// 开始打洞
	switch direction {
	case "<->":
		// 双向打洞
		bidirectionHole(srcAddr, &anotherPeer)
	case "->":
		// 连接到另一节点
		toPeerHole(srcAddr, &anotherPeer)
	case "<-":
		// 允许其他节点连接到我
		fromPeerHole(srcAddr, &anotherPeer)
	}

}

func parseAddr(addr string) (net.UDPAddr, string) {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[2])
	return net.UDPAddr{
		IP:   net.ParseIP(t[1]),
		Port: port,
	}, t[0]
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = conn.Write([]byte("handshake from:" + tag)); err != nil {
		log.Println("send handshake:", err)
	} else {
		log.Println("send handshake ok")
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		}
		log.Printf("<%s> %s\n", remoteAddr, data[:n])
	}
}

func toPeerHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("wait hole ok,err:%v", err)
		}
		msg := string(data[:n])
		if msg == HAND_SHAKE_MSG {
			log.Println("握手OK...")
			break
		}
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		}
		log.Printf("Read data:%s\n", data[:n])
	}
}

func fromPeerHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	for {
		if _, err = conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
			log.Printf("发送握手消息失败:%v", err)
		} else {
			log.Println("发送握手消息成功")
			break
		}
	}
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		}
		log.Printf("Read data:%s\n", data[:n])
		conn.Write([]byte("from " + tag))
	}
}
