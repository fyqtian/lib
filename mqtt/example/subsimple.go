package main

import (
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"github.com/fyqtian/lib/mqtt"
)

func main() {
	op := mqtt.NewOptions()
	op.AddBroker("tcp://ubuntuVM:1883")
	c, _ := mqtt.NewMqtt(op)
	ch, _ := c.SubSimple("sub/#")
	c.SubMultiple(map[string]byte{"submultiple/a": 0, "submultiple/b": 0}, func(client mqtt2.Client, message mqtt2.Message) {
		fmt.Println(5555, message.Topic(), string(message.Payload()))
	})

	go func() {
		for v := range c.ListenLostConnection() {
			fmt.Println(v)
		}
	}()
	for val := range ch {
		fmt.Println(val.Topic(), string(val.Payload()))
	}
}
