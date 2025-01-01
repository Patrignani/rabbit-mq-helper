package rabbitmqhelper

import "context"

type Connection struct {
	Url string
}

type ExchangeOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       map[string]interface{}
}

type QueueOptions struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]interface{}
}

type BindOptions struct {
	Key      string
	Exchange string
	NoWait   bool
	Args     map[string]interface{}
}

type ConsumeOptions struct {
	Consumer    string
	AutoAck     bool
	Exclusive   bool
	NoLocal     bool
	NoWait      bool
	Args        map[string]interface{}
	Action      func(ctx context.Context, body []byte)
	WorkConsume int
}

type Subscribe struct {
	Exchange ExchangeOptions
	Queue    QueueOptions
	Bind     BindOptions
	Consume  ConsumeOptions
}

type SendOptions struct {
	Exchange ExchangeOptions
}
