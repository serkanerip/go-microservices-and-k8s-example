package infra

import (
	"fmt"
	"log"

	"github.com/serkanerip/commands-service/command"
	"github.com/serkanerip/commands-service/internal/messagebroker"
	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	commandRepo command.CommandRepo
	conn        *amqp.Connection
	channel     *amqp.Channel
}

func NewRabbitMQClient(cr command.CommandRepo) *RabbitMQClient {
	conn, err := messagebroker.GetConn()
	if err != nil {
		log.Printf("cannot connect rabbit mq client err is: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Printf("cannot get channel err is: %v", err)
	}

	return &RabbitMQClient{
		commandRepo: cr,
		conn:        conn,
		channel:     channel,
	}
}

func (r *RabbitMQClient) Run() {
	go r.ListenTriggerQueue()
}

func (r *RabbitMQClient) ListenTriggerQueue() {
	cEventProcessor := command.NewCommandEventProcessor(r.commandRepo)

	if err := r.channel.ExchangeDeclare("trigger", amqp.ExchangeFanout, true, false, false, false, nil); err != nil {
		log.Printf("cannot declare exchange err is: %v", err)
	}
	q, err := r.channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Printf("Cannot declare queue err is: %v", err)
		return
	}
	err = r.channel.QueueBind(q.Name, "", "trigger", false, nil)
	if err != nil {
		log.Printf("Cannot bind queue err is: %v", err)
		return
	}

	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Printf("Cannot consume err is: %v", err)
		return
	}

	fmt.Println("Started to listen trigger queue")
	for d := range msgs {
		log.Printf("New message received from trigger queue")
		go cEventProcessor.ProcessEvent(string(d.Body))
	}

}
