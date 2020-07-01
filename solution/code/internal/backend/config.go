package backend

import "github.com/spf13/pflag"

func InitialiseFlags() {
	// define our flags
	pflag.String("host", "0.0.0.0", "host to bind server to")
	pflag.Int("port", 8080, "port to bind server to")
}
