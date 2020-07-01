package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/backend"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/util"
)

func main() {

	// set up flags (viper.Get to retrieve)
	backend.InitialiseFlags()
	// set up configuration files and parse flags
	util.InitialiseConfig("backend")

	db, err := backend.InitialisePersistence()
	if err != nil {
		panic(err)
	}

	backend.InitialiseQueue()
	// initialise http handler
	r := gin.Default()
	// set up routing
	backend.InitialiseRouter(r, db)

	// defaults to r.Run("0.0.0.0:8080")
	err = r.Run(fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port")))
	if err != nil {
		panic(err)
	}

}
