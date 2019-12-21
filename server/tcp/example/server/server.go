package main

import (
	"fmt"
	"github.com/fyqtian/lib/server/tcp"
	"github.com/fyqtian/lib/server/tcp/example"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Handler struct{}

func (this *Handler) OnConnect(c *tcp.Conn) error {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	fmt.Println("OnConnect:", addr)
	return nil
}

func (this *Handler) OnMessage(c *tcp.Conn, p tcp.Packet) error {
	echoPacket := p.(*protocol.EchoPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	c.AsyncWritePacket(protocol.NewEchoPacket(echoPacket.Serialize(), true), time.Second)
	return nil
}

func (this *Handler) OnClose(c *tcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func main() {

	srv := tcp.NewServer(tcp.SampleOptions("", "8989"), &Handler{}, &protocol.EchoProtocol{})

	go srv.Start()

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
