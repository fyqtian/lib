package main

import (
	"fmt"
	"github.com/fyqtian/lib/server/tcp/example"
	"log"
	"net"
	"time"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	echoProtocol := &protocol.EchoProtocol{}

	// ping <--> pong
	for i := 0; i < 3; i++ {
		// write
		_, err := conn.Write(protocol.NewEchoPacket([]byte("你好"), false).Serialize())
		fmt.Println(err)

		// read
		p, err := echoProtocol.ReadPacket(conn)
		if err == nil {
			echoPacket := p.(*protocol.EchoPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
		} else {
			log.Fatal(err)
		}

		time.Sleep(2 * time.Second)
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
