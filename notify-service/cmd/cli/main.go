package main

import (
	"flag"

	// This Service
	"github.com/cg917658910/fzkj-wallet/notify-service/config"
	"github.com/cg917658910/fzkj-wallet/notify-service/handlers"
	"github.com/cg917658910/fzkj-wallet/notify-service/svc/server"
)

func main() {
	// Setup flags
	flag.Parse()
	// Setup config
	config.Setup()

	cfg := server.DefaultConfig

	cfg = handlers.SetConfig(cfg)

	server.Run(cfg)
}
