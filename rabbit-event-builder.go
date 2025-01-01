package rabbitmqhelper

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitEventBuider struct {
	con              Connection
	SubscribeOptions map[string]*Subscribe
	Channels         map[string]*amqp.Channel
}

func NewRabbitEventBuider(url string) *RabbitEventBuider {
	return &RabbitEventBuider{
		con: Connection{
			Url: url,
		},
		SubscribeOptions: make(map[string]*Subscribe),
		Channels:         make(map[string]*amqp.Channel),
	}
}

func (r *RabbitEventBuider) Subscribe(name string, sub *Subscribe) *RabbitEventBuider {
	r.SubscribeOptions[name] = sub

	return r
}

func (r *RabbitEventBuider) Run(ctx context.Context) {
	conn, err := amqp.Dial(r.con.Url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	var forever chan struct{}

	for key, s := range r.SubscribeOptions {
		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		r.Channels[key] = ch

		err = ch.ExchangeDeclare(
			s.Exchange.Name,
			s.Exchange.Kind,
			s.Exchange.Durable,
			s.Exchange.AutoDelete,
			s.Exchange.Internal,
			s.Exchange.NoWait,
			s.Exchange.Args,
		)
		failOnError(err, "Failed to declare an exchange")

		q, err := ch.QueueDeclare(
			key,
			s.Queue.Durable,
			s.Queue.AutoDelete,
			s.Queue.Exclusive,
			s.Queue.NoWait,
			s.Queue.Args,
		)
		failOnError(err, "Failed to declare a queue")

		err = ch.QueueBind(
			q.Name,
			s.Bind.Key,
			s.Bind.Exchange,
			s.Bind.NoWait,
			s.Bind.Args,
		)
		failOnError(err, "Failed to bind a queue")

		msgs, err := ch.Consume(
			q.Name,
			s.Consume.Consumer,
			s.Consume.AutoAck,
			s.Consume.Exclusive,
			s.Consume.NoLocal,
			s.Consume.NoWait,
			s.Consume.Args,
		)
		failOnError(err, "Failed to register a consumer")

		if s.Consume.WorkConsume == 0 {
			s.Consume.WorkConsume = 1
		}

		go func(action func(ctx context.Context, body []byte)) {
			semaphore := make(chan struct{}, s.Consume.WorkConsume)

			for d := range msgs {
				semaphore <- struct{}{}
				go func(body []byte) {
					defer func() { <-semaphore }()
					action(ctx, body)
				}(d.Body)
			}
		}(s.Consume.Action)

		log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	}

	<-forever

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
