package handlers

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg"
	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/codes"
	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/types"
	pb "github.com/cg917658910/fzkj-wallet/wash-service/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.WashServiceServer {
	return washserviceService{}
}

type washserviceService struct{}

func (s washserviceService) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	var resp pb.HealthCheckResponse
	return &resp, nil
}

func (s washserviceService) WashShortMsg(ctx context.Context, in *pb.WashShortMsgRequest) (*pb.WashShortMsgResponse, error) {
	var resp pb.WashShortMsgResponse
	req := &types.WashRequest{
		MsgId:    in.MsgId,
		Msg:      in.Msg,
		Currency: in.Currency,
		TypeName: in.TypeName,
	}
	res := shortmsg.WashShortMsg(ctx, req)
	if res == nil {
		resp.Code = int32(codes.ErrInternal.Code)
		resp.Msg = "WashShortMsg return nil"
		return &resp, nil
	}
	resp.Code = int32(res.Code.Code)
	resp.Msg = res.Code.Msg
	if res.Extracted != nil {
		resp.Data = &pb.WashShortMsgResult{
			PayTime: res.PayTime,
			PayCoin: res.PayCoin,
			Balance: res.Balance,
		}
	}
	return &resp, nil
}
