package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"log"
	"net"
	"os"
	"time"
)

const messageId = 0

type Message string

var server *p2p.Server

func MyProtocol() p2p.Protocol {
	return p2p.Protocol{
		Name:    "MyProtocol",
		Version: 1,
		Length:  1,
		Run:     msgHandler,
	}
}

func main() {
	nodekey, _ := crypto.GenerateKey()
	srv := p2p.Server{
		Config: p2p.Config{
			MaxPeers:   10,
			PrivateKey: nodekey,
			Name:       "my node name",
			ListenAddr: ":30300",
			Protocols:  []p2p.Protocol{MyProtocol()},
			BootstrapNodes: []*discover.Node{
				&discover.Node{
					IP:  net.ParseIP("120.27.209.161"),
					UDP: 30300,
					TCP: 30300,
					ID:  discover.MustHexID("6dc090ed5c0472ca71feff2ae7242d5cc6aaf16de8a24bef3827666c3e6ba5c39665d51824b06554a6ef4d4a13e3bd37423d92b24e432785d4bb35568733cff4"),
				},
			},
		},
	}
	server = &srv
	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("started..", srv.NodeInfo())
	select {}
}

func msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
	fmt.Println(peer, " in,total peers:", server.PeerCount())
	for {
		msg, err := ws.ReadMsg()
		if err != nil {
			return err
		}

		var myMessage Message
		err = msg.Decode(&myMessage)
		if err != nil {
			// handle decode error
			continue
		}

		fmt.Println("code:", msg.Code, "receiver at:", msg.ReceivedAt, "msg:", myMessage)
		switch myMessage {
		case "foo":
			err := p2p.SendItems(ws, messageId, "bar")
			if err != nil {
				return err
			}
		default:
			fmt.Println("recv:", myMessage)
		}
	}
}

func startPeer(serverAddr string) {
	conn, err := net.DialTimeout("tcp", serverAddr, 5*time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	log.Println("peer ok")
	time.Sleep(5 * time.Minute)
}
