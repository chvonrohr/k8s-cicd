package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/crawler"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/util"
)

var (
	username, password, host, port string
)

func main() {

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, host, port))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer util.FailOnErrorF(conn.Close, "Failed to close RabbitMQ connection")
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer util.FailOnErrorF(ch.Close, "Failed to close channel")

	q, err := ch.QueueDeclare("pages", true, false, false, false, nil)

	util.FailOnError(ch.Qos(5, 5, false), "Failed to set rabbitmq qos")

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	for {
		select {
		case msg := <-msgs:
			{
				// todo: consume msg
				var page model.Page
				json.Unmarshal(msg.Body, &page)
				crawler.Crawl(page.Url)
			}
		}
	}

}
