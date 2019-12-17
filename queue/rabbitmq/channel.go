package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"time"
)

type Channel struct {
	channel *amqp.Channel
}

func (s *Channel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (*amqp.Queue, error) {
	q, err := s.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
	if err != nil {
		return nil, errors.WithMessage(err, "declare queue")
	}
	//todo
	//point
	return &q, nil
}

func (s *Channel) QueueSample(name string, durable bool) (*amqp.Queue, error) {
	return s.QueueDeclare(name, durable, false, false, false, nil)
}

func (s *Channel) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	if err := s.channel.Publish(exchange, key, mandatory, immediate, msg); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("publish exchange:%s key:%s", exchange, key))
	}
	return nil
}

func (s *Channel) PushlishSample(exchange, route string, payload []byte) error {
	t := amqp.Publishing{
		ContentType: "text/plain",
		Body:        payload,
	}
	return s.Publish(exchange, route, false, false, t)
}

func (s *Channel) Consumer(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := s.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(3 * time.Second)

		}
	}()

	return deliveries, nil
}

func (s *Channel) consumerErrorStr(queue string) string {
	return fmt.Sprintf("consume queue:%s", queue)
}

func (s *Channel) ConusmerSample(queue string, autoAck bool) (<-chan amqp.Delivery, error) {
	return s.Consumer(queue, "", autoAck, false, false, false, nil)
}

func (s *Channel) exchangeErrorStr(name, exchangeType string) string {
	return fmt.Sprintf("exchange name:%s, type:%s", name, exchangeType)
}

func (s *Channel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	if err := s.channel.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args); err != nil {
		return errors.WithMessage(err, s.exchangeErrorStr(name, kind))
	}
	return nil
}

func (s *Channel) ExchangeSample(name, kind string, durable bool) error {
	return s.ExchangeDeclare(name, kind, durable, false, false, false, nil)
}

func (s *Channel) queueBindErrorStr(name, key, exchange string) string {
	return fmt.Sprintf("queue:%s key:%s exchange:%s", name, key, exchange)
}

func (s *Channel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	if err := s.channel.QueueBind(name, key, exchange, noWait, args); err != nil {
		return errors.WithMessage(err, s.queueBindErrorStr(name, key, exchange))
	}
	return nil
}

func (s *Channel) QueueBindSample(name, key, exchange string) error {
	return s.QueueBind(name, key, exchange, false, nil)
}
