package delivery

import (
	"github.com/nats-io/nats.go"
	"log"
)

const (
	channel = "orders"
)

type Nats struct {
	conn *nats.Conn
}

func NewNats() (*Nats, error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	return &Nats{conn: conn}, nil
}

func (n *Nats) GetConsumer() (msg chan *nats.Msg, err error) {
	msg = make(chan *nats.Msg, 10)
	_, err = n.conn.ChanSubscribe(channel, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (n *Nats) PutMessage(data []byte) error {
	return n.conn.Publish(channel, data)
}
