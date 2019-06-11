package main

import (
	"flag"

	"github.com/tinywell/fabclient/examples/client"
)

var (
	configFile string
)

func main() {
	flag.StringVar(&configFile, "config", "./config,yaml", "config path")
	flag.Parse()
	client.Run(client.Cfg{ConfigFile: configFile})
}
