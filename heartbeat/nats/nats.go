package nats

import (
	"context"
	nats "github.com/nats-io/nats.go"
)

type Nats struct {
	address string
	topic   string
	nc      *nats.Conn
}

func New(address, topic string) (*Nats, error) {
	nc, err := nats.Connect(address)
	if err != nil {
		return nil, err
	}

	return &Nats{
		address: address,
		topic:   topic,
		nc:      nc,
	}, nil
}

func (ns *Nats) Watch(ctx context.Context, f func(data string)) {
	ns.nc.Subscribe(ns.topic, func(msg *nats.Msg) {
		f(string(msg.Data))
	})
	ns.nc.Flush()
}
