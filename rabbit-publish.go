package rabbitmqhelper

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublish struct {
	con               Connection
	SendOptions       map[string]*SendOptions
	sendTimeoutSecond int
}

func NewRabbitPublish(url string, sendTimeoutSecond int) *RabbitPublish {
	return &RabbitPublish{
		con: Connection{
			Url: url,
		},
		sendTimeoutSecond: sendTimeoutSecond,
		SendOptions:       make(map[string]*SendOptions),
	}
}

func (p *RabbitPublish) RegisterSendOptions(options *SendOptions) *RabbitPublish {

	p.SendOptions[options.Exchange.Name] = options

	return p
}

func (p *RabbitPublish) Send(ctx context.Context, data any, name, key string, mandatory, immediate bool) error {
	conn, err := amqp.Dial(p.con.Url)

	if err != nil {
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		return err
	}

	defer ch.Close()

	ex := p.SendOptions[name]

	err = ch.ExchangeDeclare(
		ex.Exchange.Name,
		ex.Exchange.Kind,
		ex.Exchange.Durable,
		ex.Exchange.AutoDelete,
		ex.Exchange.Internal,
		ex.Exchange.NoWait,
		ex.Exchange.Args,
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(p.sendTimeoutSecond)*time.Second)
	defer cancel()

	body, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx,
		ex.Exchange.Name,
		key,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		return err
	}

	log.Printf(" [x] Sent %s", body)

	return err
}
