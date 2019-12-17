package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fyqtian/lib/mqtt"
	"log"
	"time"
)

func main() {
	op := mqtt.NewOptions()
	op.AddBroker("tcp://ubuntuVM:1883")
	op.SetClientID("tttttt")
	op.SetUsername("pushCore")
	op.SetPassword("pushCore")
	op.OnConnect = nil
	if c, err := mqtt.NewMqtt("abc", op); err != nil {
		log.Fatal(err)
	} else {
		//mqtt.Debug()
		c.Sub("sss", 0, func(client MQTT.Client, message MQTT.Message) {
			fmt.Println(string(message.Payload()))
		})
		time.Sleep(2e9)
		if token := c.GetClient().Publish("sss", 0, false, "first"); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}

		log.Println("connect success")
		c.GetClient().Disconnect(1)
		if c, err = mqtt.NewMqtt("da", op); err != nil {
			log.Fatal("dddd")
		}
		c.Sub("sssa", 0, func(client MQTT.Client, message MQTT.Message) {
			fmt.Println(string(message.Payload()))
		})
		time.Sleep(5e9)
		if token := c.GetClient().Publish("sssa", 0, false, "555"); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}
		log.Println()
		time.Sleep(200e9)
	}
}
