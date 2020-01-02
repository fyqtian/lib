package main

import (
	"fmt"
	"github.com/fyqtian/lib/mqtt"
)

func main() {
	op := mqtt.NewOptions()
	op.AddBroker("tcp://ubuntuVM:1883")
	c, _ := mqtt.NewMqtt(op)
	fmt.Println(<-c.ListenLostConnection())
}
