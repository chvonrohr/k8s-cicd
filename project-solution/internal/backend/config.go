package backend

import (
	"github.com/spf13/pflag"
)

// InitialiseFlags registers a set of flags for use with the cli with our commander library.
func InitialiseFlags() {
	// server flags
	pflag.String("host", "0.0.0.0", "host to bind server to")
	pflag.Int("port", 8080, "port to bind server to")

	// persistence flags
	pflag.String("db.type", "sqlite", "database type") // sqlite or postgres
	pflag.String("db.host", "localhost", "database hostname")
	pflag.Int("db.port", 3306, "database port")
	pflag.String("db.username", "root", "database username")
	pflag.String("db.password", "", "database password")
	pflag.String("db.database", "default", "database name")

	// queue flags
	pflag.String("queue.host", "localhost", "queue host")
	pflag.String("queue.username", "guest", "queue username")
	pflag.String("queue.password", "guest", "queue password")
	pflag.Int("queue.port", 5672, "queue port")

}
