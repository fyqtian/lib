package main

import (
	"fmt"
	"github.com/fyqtian/lib/queue/rabbitmq"
	"log"
	"time"
)

func create() (*rabbitmq.Helper, error) {
	op := &rabbitmq.Options{
		Host:   "ubuntuVM",
		Port:   "5672",
		User:   "guest",
		Passwd: "guest",
		Vhost:  "test",
	}
	if client, err := rabbitmq.NewWithRetry("mq", op, 0, 10*time.Second); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func main() {
	var exchangeName = "test-exchange"
	c, err := create()
	if err != nil {
		log.Fatal(err)
	}
	if ch, err := c.Channel(); err != nil {
		log.Fatal(err)
	} else {
		if err := ch.ExchangeSample(exchangeName, "fanout", true); err != nil {
			log.Fatal(err)
		}

		i := 0
		for range time.Tick(5e9) {
			i++
			str := fmt.Sprintf("%d pushlish message", i)
			err := ch.PushlishSample(exchangeName, "", []byte(str))
			log.Println(err, str)
		}
	}
}
