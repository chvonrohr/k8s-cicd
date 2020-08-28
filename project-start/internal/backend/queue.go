package backend

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/model"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/util"
)

var (
	queue   amqp.Queue
	channel *amqp.Channel
)

/**
InitialiseQueue initialises a new connection to a rabbitmq server.

It is non-blocking and returns a function used to close the connection.
It expects viper and the existence of the following four viper variables: queue.username, queue.password, queue.host and queue.port.
It will exit fatally if the connection fails.
It declares the pages queue if it doesn't exist. The queue is durable and non-exclusive.
*/
func InitialiseQueue() (cancel func()) {

	// configuration variables
	var (
		username = viper.GetString("queue.username")
		password = viper.GetString("queue.password")
		host     = viper.GetString("queue.host")
		port     = viper.GetInt("queue.port")
	)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, host, port))
	util.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare("pages", true, false, false, false, nil)

	util.FailOnError(ch.Qos(5, 0, false), "Failed to set rabbitmq qos")

	queue = q
	channel = ch

	cancel = func() {
		util.FailOnErrorF(ch.Close, "Failed to close channel")
		util.FailOnErrorF(conn.Close, "Failed to close RabbitMQ connection")
	}
	return

}

/**
QueuePage pushes a page to rabbitmq.
*/
func QueuePage(page model.Page) error {
	bs, err := json.Marshal(&page)
	if err != nil {
		return err
	}
	err = channel.Publish("", queue.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bs,
		},
	)
	return err
}
