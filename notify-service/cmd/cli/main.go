package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	// This Service
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order"
	"github.com/cg917658910/fzkj-wallet/notify-service/config"
	"github.com/cg917658910/fzkj-wallet/notify-service/handlers"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/cache"
	"github.com/cg917658910/fzkj-wallet/notify-service/svc/server"
)

func main() {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*30)
	defer cancle()
	// Setup flags
	flag.Parse()
	// Setup config
	config.Setup()

	cfg := server.DefaultConfig

	cfg = handlers.SetConfig(cfg)
	cache.SetupRedis(ctx)
	order.OrderNotifyStart(ctx)
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
