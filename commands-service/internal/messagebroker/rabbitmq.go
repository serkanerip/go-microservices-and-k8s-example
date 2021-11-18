package messagebroker

import (
	"github.com/serkanerip/commands-service/config"
	"github.com/streadway/amqp"
)

func GetConn() (conn *amqp.Connection, err error) {
	conn, err = amqp.Dial(config.ENV.MESSAGE_BROKER_CS)
	return
}
