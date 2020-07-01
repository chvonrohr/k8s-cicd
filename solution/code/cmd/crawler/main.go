package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"

	"github.com/streadway/amqp"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/crawler"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/sdk"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/util"
)

func main() {

	// set up flags (viper.Get to retrieve)
	crawler.InitialiseFlags()
	// set up configuration files and parse flags
	util.InitialiseConfig("crawler")

	var (
		username = viper.GetString("queue.username")
		password = viper.GetString("queue.password")
		host     = viper.GetString("queue.host")
		port     = viper.GetInt("queue.port")
	)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, host, port))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	defer util.FailOnErrorF(conn.Close, "Failed to close RabbitMQ connection")
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer util.FailOnErrorF(ch.Close, "Failed to close channel")

	q, err := ch.QueueDeclare("pages", true, false, false, false, nil)

	util.FailOnError(ch.Qos(5, 0, false), "Failed to set rabbitmq qos")

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	for {
		select {
		case msg := <-msgs:
			{
				var page model.Page
				err = json.Unmarshal(msg.Body, &page)
				if err != nil {
					log.Println(err)
					continue
				}
				urls, err := crawler.Crawl(page.Url)
				if err != nil {
					log.Println(err)
					continue
				}
				err = sdk.PageCallback(page, urls)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

}
