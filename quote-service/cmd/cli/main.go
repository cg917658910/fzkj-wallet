package main

import (
	"flag"

	// This Service

	"github.com/cg917658910/fzkj-wallet/quote-service/config"
	"github.com/cg917658910/fzkj-wallet/quote-service/handlers"
	"github.com/cg917658910/fzkj-wallet/quote-service/svc/server"
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
