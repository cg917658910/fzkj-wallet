package order

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/consumer"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
)

type OrderNotifyStartResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

var consumerManager *consumer.MyConsumerManager

var logger = log.DLogger()

func OrderNotifyStart(ctx context.Context) (*OrderNotifyStartResponse, error) {
	var resp *OrderNotifyStartResponse
	if consumerManager == nil {
		consumerManager = consumer.NewConsumerManager(context.Background())
		if err := consumerManager.Start(); err != nil {
			consumerManager = nil
			return nil, err
		}
	}
	resp = &OrderNotifyStartResponse{
		Code: 0,
		Msg:  "Ok",
	}
	return resp, nil
}
func OrderNotifyStop(ctx context.Context) error {
	logger.Info("Stopping Consumer Manager...")
	if consumerManager == nil {
		logger.Info("Consumer Manager stopped successfully")
		return nil
	}
	if err := consumerManager.Stop(); err != nil {
		return err
	}
	consumerManager = nil
	logger.Info("Consumer Manager stopped successfully")
	return nil
}
