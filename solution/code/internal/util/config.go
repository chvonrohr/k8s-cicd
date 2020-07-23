package util

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitialiseConfig(name string) {

	// look for env variables in the format "LETSBOOT_PORT=1338"
	viper.SetEnvPrefix("LETSBOOT")

	// look for config files with name name.yml, name.toml, name.json...
	viper.SetConfigName(name)

	// ... in these folders
	viper.AddConfigPath("/etc/letsboot")
	viper.AddConfigPath("$HOME/.letsboot")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".") // working directory

	// parse flags from process arg list
	pflag.Parse()

	// bind parsed flags to config library
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	// check for environment variables now
	viper.AutomaticEnv()

	// try to find and read config file now
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// watch config file for changes
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

}
