package handlers

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/notify-service/app/enum"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order"
	pb "github.com/cg917658910/fzkj-wallet/notify-service/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.NotifyServiceServer {
	return notifyserviceService{}
}

type notifyserviceService struct{}

func (s notifyserviceService) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	var resp pb.HealthCheckResponse
	resp.Status = int32(enum.HealthStatus_SERVING)
	resp.Message = enum.HealthStatus_SERVING.String()
	return &resp, nil
}

func (s notifyserviceService) OrderNotifyStart(ctx context.Context, in *pb.OrderNotifyStartRequest) (*pb.OrderNotifyStartResponse, error) {
	var resp *pb.OrderNotifyStartResponse
	result, err := order.OrderNotifyStart(ctx)
	if err != nil {
		return nil, err
	}
	resp = &pb.OrderNotifyStartResponse{
		Code: result.Code,
		Msg:  result.Msg,
	}
	return resp, nil
}

func (s notifyserviceService) OrderNotifyStop(ctx context.Context, in *pb.OrderNotifyStopRequest) (*pb.OrderNotifyStopResponse, error) {
	var resp pb.OrderNotifyStopResponse
	err := order.OrderNotifyStop(ctx)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
