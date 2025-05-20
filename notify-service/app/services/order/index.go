package order

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/scheduler"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/log"
)

type OrderNotifyStartResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}

var cgScheduler scheduler.Scheduler

var logger = log.DLogger()

func OrderNotifyStart(ctx context.Context) (*OrderNotifyStartResponse, error) {
	var resp *OrderNotifyStartResponse
	if cgScheduler == nil {
		scheduler := scheduler.NewScheduler()
		if err := scheduler.Init(); err != nil {
			return nil, err
		}
		if err := scheduler.Start(); err != nil {
			return nil, err
		}
		cgScheduler = scheduler
	}
	resp = &OrderNotifyStartResponse{
		Code: 0,
		Msg:  "Ok",
	}
	return resp, nil
}
func OrderNotifyStop(ctx context.Context) error {
	logger.Info("Stopping Notify Scheduler...")
	if cgScheduler == nil {
		logger.Info("Notify Scheduler stopped successfully")
		return nil
	}
	if err := cgScheduler.Stop(); err != nil {
		return err
	}
	cgScheduler = nil
	logger.Info("Notify Scheduler stopped successfully")
	return nil
}
