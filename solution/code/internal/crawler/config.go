package crawler

import "github.com/spf13/pflag"

// InitialiseFlags registers a set of flags for use with the cli with our commander library.
func InitialiseFlags() {

	pflag.String("backend.url", "http://backend", "base url of backend")

	pflag.String("queue.host", "localhost", "host of rabbitmq server")
	pflag.String("queue.username", "guest", "username of rabbitmq server")
	pflag.String("queue.password", "guest", "password of rabbitmq server")
	pflag.Int("queue.port", 5672, "port of rabbitmq server")

	pflag.String("crawler.data", "/var/data", "data directory for page responses")
	pflag.Bool("crawler.dump", true, "dump page responses to data directory")
}
