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

		q, err := ch.QueueSample("test", true)
		if err != nil {
			log.Fatal(q)
		}
		if err := ch.QueueBindSample(q.Name, "van", exchangeName); err != nil {
			log.Fatal(err)
		}
		if msgs, err := ch.ConusmerSample(q.Name, true); err != nil {
			log.Fatal(err)
		} else {
			for d := range msgs {
				fmt.Println(string(d.Body))
			}
		}
	}
}
