package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"

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

	log.Println("listening to queue...")

	for {
		select {
		case msg := <-msgs:
			{
				go func() {
					var page model.Page
					err = json.Unmarshal(msg.Body, &page)
					if err != nil {
						log.Println(err)
						return
					}
					log.Printf("crawling %s\n", page.Url)
					response, err := crawler.Crawl(page.Url)
					if err != nil {
						log.Printf("error crawling %s: %s", page.Url, err)
						if response.StatusCode == 0 {
							return
						}
					}
					err = sdk.PageCallback(page, response)
					if err != nil {
						log.Println(err)
					}
					err = msg.Ack(false)
					if err != nil {
						log.Println(err)
					}
					log.Printf("processed %s\n", page.Url)
				}()
			}
		}
	}
}
