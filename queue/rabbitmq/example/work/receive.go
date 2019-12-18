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
	if client, err := rabbitmq.NewWithRetry(op, 0, 10*time.Second); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func main() {
	//var exchangeName = ""
	var querueName = "test-q"

	c, err := create()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if ch, err := c.Channel(); err != nil {
			log.Fatal(err)
		} else {
			log.Println("start work 1")
			q, err := ch.QueueSample(querueName, true)
			if err != nil {
				log.Fatal(q)
			}

			if msgs, err := ch.ConusmerSample(q.Name, true); err != nil {
				log.Fatal(err)
			} else {
				for d := range msgs {
					fmt.Println("[worke1 receive]", string(d.Body))
				}
			}
		}
	}()

	if ch, err := c.Channel(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("start work 2")

		q, err := ch.QueueSample(querueName, true)
		if err != nil {
			log.Fatal(q)
		}

		if msgs, err := ch.ConusmerSample(q.Name, true); err != nil {
			log.Fatal(err)
		} else {
			for d := range msgs {
				fmt.Println("[worke2 receive]", string(d.Body))
			}
		}
	}
}
