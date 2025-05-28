package shortmsg

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/THB"
	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/codes"
	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/types"
)

func WashShortMsg(ctx context.Context, req *types.WashRequest) (resp *types.WashResponse) {
	resp = &types.WashResponse{
		WashRequest: req,
		Extracted:   &types.Extracted{},
	}
	if req == nil {
		resp.Code = codes.ErrInvalidArgument
		return
	}
	if err := req.Check(); err != nil {
		resp.Code = codes.ErrInvalidArgument.New(err.Error())
		return
	}
	switch req.Currency {
	case "THB":
		resp.Extracted = THB.WashMsg(ctx, req.TypeName, req.Msg)
	default:
		resp.Code = codes.ErrUnsupportedCurrencyType
	}
	return
}
