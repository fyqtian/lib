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
	var exchangeName = ""
	var querueName = "test-q"

	c, err := create()
	if err != nil {
		log.Fatal(err)
	}
	if ch, err := c.Channel(); err != nil {
		log.Fatal(err)
	} else {
		q, err := ch.QueueSample(querueName, true)
		if err != nil {
			log.Fatal(q)
		}
		i := 0
		for range time.Tick(5e9) {
			i++
			str := fmt.Sprintf("%d pushlish message", i)
			err := ch.PushlishSample(exchangeName, q.Name, []byte(str))
			log.Println(err, str)
		}
	}
}
