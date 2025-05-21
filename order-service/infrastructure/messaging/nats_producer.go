package messaging

import (
	"encoding/json"
	"order-service/domain"

	"github.com/nats-io/nats.go"
)

type NATSProducer struct {
	conn *nats.Conn
}

func NewNATSProducer(url string) (*NATSProducer, error) {
	nc, err := nats.Connect(url, nats.UserInfo("natsuser", "natspass"))
	if err != nil {
		return nil, err
	}
	return &NATSProducer{conn: nc}, nil
}

func (p *NATSProducer) PublishOrderCreated(order domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	subject := "order.created"
	return p.conn.Publish(subject, data)
}
