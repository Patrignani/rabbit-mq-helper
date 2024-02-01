package main

import (
	"context"

	rabbitmqhelper "github.com/Patrignani/rabbit-mq-helper"
)

func main() {

	ctx := context.Background()

	ex := rabbitmqhelper.ExchangeOptions{
		Name:       "logs",
		Kind:       "fanout",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}

	q := rabbitmqhelper.QueueOptions{
		Durable:    false,
		AutoDelete: false,
		Exclusive:  true,
		NoWait:     false,
		Args:       nil,
	}

	b := rabbitmqhelper.BindOptions{
		Key:      "",
		Exchange: "logs",
		NoWait:   false,
		Args:     nil,
	}

	s1 := &rabbitmqhelper.Subscribe{
		Exchange: ex,
		Queue:    q,
		Bind:     b,
		Consume: rabbitmqhelper.ConsumeOptions{
			Consumer:  "",
			AutoAck:   true,
			Exclusive: false,
			NoLocal:   false,
			NoWait:    false,
			Args:      nil,
			Action: func(ctx context.Context, body []byte) {
				println("test1", string(body))
			},
		},
	}

	s2 := &rabbitmqhelper.Subscribe{
		Exchange: ex,
		Queue:    q,
		Bind:     b,
		Consume: rabbitmqhelper.ConsumeOptions{
			Consumer:  "",
			AutoAck:   true,
			Exclusive: false,
			NoLocal:   false,
			NoWait:    false,
			Args:      nil,
			Action: func(ctx context.Context, body []byte) {
				println("test2", string(body))
			},
		},
	}

	s3 := &rabbitmqhelper.Subscribe{
		Exchange: rabbitmqhelper.ExchangeOptions{
			Name:       "aaaaa",
			Kind:       "fanout",
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Args:       nil,
		},
		Queue: q,
		Bind: rabbitmqhelper.BindOptions{
			Key:      "",
			Exchange: "aaaaa",
			NoWait:   false,
			Args:     nil,
		},
		Consume: rabbitmqhelper.ConsumeOptions{
			Consumer:  "",
			AutoAck:   true,
			Exclusive: false,
			NoLocal:   false,
			NoWait:    false,
			Args:      nil,
			Action: func(ctx context.Context, body []byte) {
				println("test3", string(body))
			},
		},
	}

	rabbitmqhelper.NewRabbitEventBuider("amqp://guest:guest@localhost:5672/").
		Subscribe("test1", s1).
		Subscribe("test2", s2).
		Subscribe("test3", s3).
		Run(ctx)

}
