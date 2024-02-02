package main

import (
	"context"
	"log"
	"time"

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

func send() {
	ctx := context.Background()

	publish := rabbitmqhelper.NewRabbitPublish("amqp://guest:guest@localhost:5672/", 5).
		RegisterSendOptions(&rabbitmqhelper.SendOptions{
			Exchange: rabbitmqhelper.ExchangeOptions{
				Name:       "cash-flow",
				Kind:       "fanout",
				Durable:    true,
				AutoDelete: false,
				Internal:   false,
				NoWait:     false,
			},
		})

	body := Auditory{
		EntityId: "22315",
		DatedIn:  time.Now(),
		Type:     "test",
		Entity:   "test",
		ClientId: "56898",
		UserId:   "465463",
		Data: map[string]any{
			"test":  "aaaaa",
			"valor": 1000.08,
		},
	}

	publish.Send(ctx, body, "cash-flow", "", false, false)

	log.Printf(" [x] Sent %s", body)
}

type Auditory struct {
	Id       string    `json:"id,omitempty" bson:"_id,omitempty"`
	EntityId string    `json:"entityId" bson:"entity_id"`
	Type     string    `json:"type" bson:"type"`
	Entity   string    `json:"entity" bson:"entity"`
	DatedIn  time.Time `json:"datedIn" bson:"dated_in"`
	ClientId string    `json:"clientId" bson:"client_id"`
	UserId   string    `json:"userId" bson:"user_id"`
	Data     any       `json:"data" bson:"data"`
}
