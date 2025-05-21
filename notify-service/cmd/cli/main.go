package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	startMetricsServer()
	server.Run(cfg)
}

func startMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Metrics + pprof running on :9090")
	go func() {
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()
}
