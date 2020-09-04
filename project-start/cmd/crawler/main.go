package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/streadway/amqp"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/crawler"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/model"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/sdk"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/util"
)

func main() {

	// set up flags (viper.Get to retrieve)
	crawler.InitialiseFlags()
	// set up configuration files and parse flags
	util.InitialiseConfig("crawler")

	var (
		// query queue connection variables from config library
		username = viper.GetString("queue.username")
		password = viper.GetString("queue.password")
		host     = viper.GetString("queue.host")
		port     = viper.GetInt("queue.port")

		// query backend endpoint url from config library
		endpoint = viper.GetString("backend.url")
	)

	// create a new crawler sdk client (view internal/sdk)
	apiClient := sdk.NewClient(endpoint)

	// dial (connect) to rabbitmq
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, host, port))
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	// FailOnErrorF only runs the passed function when the deferred function is executed
	defer util.FailOnErrorF(conn.Close, "Failed to close RabbitMQ connection")
	// establish a new channel
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	defer util.FailOnErrorF(ch.Close, "Failed to close channel")

	// declare a durable queue on the current channel
	q, err := ch.QueueDeclare("pages", true, false, false, false, nil)

	// ensure we crawl a maximum of five pages concurrently
	util.FailOnError(ch.Qos(5, 0, false), "Failed to set rabbitmq qos")

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	log.Printf("data dir: %s\n", viper.GetString("crawler.data"))
	log.Println("listening to queue...")

	for {
		select {
		case msg := <-msgs:
			{
				// go func() {
				var page model.Page
				err = json.Unmarshal(msg.Body, &page)
				if err != nil {
					log.Println(err)
					// return
					continue
				}
				log.Printf("crawling %s\n", page.Url)
				response, err := crawler.Crawl(page.Url, page.CrawlID)
				if err != nil {
					log.Printf("error crawling %s: %s", page.Url, err)
					if response.StatusCode == 0 {
						// return
						continue
					}
				}
				err = apiClient.PageCallback(page, response)
				if err != nil {
					log.Println(err)
				}
				err = msg.Ack(false)
				if err != nil {
					log.Println(err)
				}
				log.Printf("processed %s\n", page.Url)
				// }()
			}
		}
	}
}
