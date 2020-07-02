package crawler

import "github.com/spf13/pflag"

func InitialiseFlags() {

	pflag.String("backend.url", "http://backend", "base url of backend")

	pflag.String("queue.host", "localhost", "host of rabbitmq server")
	pflag.String("queue.username", "guest", "username of rabbitmq server")
	pflag.String("queue.password", "guest", "password of rabbitmq server")
	pflag.Int("queue.port", 5672, "port of rabbitmq server")
}
