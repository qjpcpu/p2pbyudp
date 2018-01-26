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
var switchTcp bool

const HAND_SHAKE_MSG = "我是打洞消息"

func main() {
	// 当前进程标记字符串,便于显示
	tag = os.Args[1]
	switchTcp = len(os.Args) == 3 && os.Args[2] == "tcp"

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
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, anotherPeer.String())

	// 开始打洞
	bidirectionHole(srcAddr, &anotherPeer)

}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, anotherAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, anotherAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	// 向另一个peer发送一条udp消息(对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
		log.Println("send handshake:", err)
	}
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from [" + tag + "]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	loop := 2
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("收到数据:%s\n", data[:n])
		}
		if loop >= 0 {
			loop -= 1
		}
		if loop < 0 && switchTcp {
			break
		}
	}
	log.Println("TCP.....")
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", srcAddr.Port))
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	go sendTcp(anotherAddr.IP.String(), fmt.Sprint(anotherAddr.Port))
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go handleConn(c)
	}
}

func sendTcp(ip, port string) {
	//打开连接:
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		//由于目标计算机积极拒绝而无法创建连接
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}

	for {
		if _, err = conn.Write([]byte(fmt.Sprintf("FROM %s %v", tag, time.Now().Nanosecond()))); err != nil {
			log.Println("write tcp fail:", err)
		}
		time.Sleep(5 * time.Second)
	}
}
func handleConn(c net.Conn) {
	log.Println("start to read from conn")
	for {
		var buf = make([]byte, 10)
		n, err := c.Read(buf)
		if err != nil {
			log.Println("conn read error:", err)
			return
		}
		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}
