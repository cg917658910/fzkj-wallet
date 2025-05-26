package handlers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order"
	"github.com/cg917658910/fzkj-wallet/notify-service/svc"
)

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)

	// Place whatever shutdown handling you want here
	// stop notify
	ctx, cancle := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancle()
	order.OrderNotifyStop(ctx)
	errc <- terminateError
}

func SetConfig(cfg svc.Config) svc.Config {
	return cfg
}
