package platform

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/streadway/amqp"
)

type MessageBusClient interface {
	PublishNewPlatform(PlatformPublishedDTO) error
}

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQClient(conn *amqp.Connection) *RabbitMQClient {
	channel, err := conn.Channel()
	if err != nil {
		log.Printf("cannot get channel err is: %v", err)
	}
	return &RabbitMQClient{
		conn:    conn,
		channel: channel,
	}
}

func (r *RabbitMQClient) PublishNewPlatform(dto PlatformPublishedDTO) error {
	if r.conn.IsClosed() {
		return errors.New("rabbit mq connection is closed cannot publish new platform")
	}
	if err := r.channel.ExchangeDeclare("trigger", amqp.ExchangeFanout, true, false, false, false, nil); err != nil {
		log.Printf("cannot get declare exchange err is: %v", err)
	}

	b, err := json.Marshal(&dto)
	if err != nil {
		log.Printf("cannot marshal dto err is: %v", err)
	}

	r.sendMessage(b)
	return nil
}

func (r *RabbitMQClient) sendMessage(message []byte) {
	r.channel.Publish("trigger", "", false, false, amqp.Publishing{
		Body:        message,
		ContentType: "text/plain",
	})
}
