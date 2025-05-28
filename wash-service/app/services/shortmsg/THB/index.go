package THB

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/codes"
	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/types"
)

func WashMsg(ctx context.Context, type_name, msg string) (res *types.Extracted) {
	res = &types.Extracted{}
	switch type_name {
	case "SCB", "SCB读取":
		res = ExtractSCB(msg)
	case "SCB通知":
		res = ExtractSCBNotify(msg)
	case "SCB流水":
		res = ExtractSCBWater(msg)
	case "KTB":
		res = ExtractKTB(msg)
	case "KTBLine":
		res = ExtractKTB(msg)
	case "KTB通知":
		res = ExtractKTBNotice(msg)
	case "KTB流水":
		res = ExtractKTBWater(msg)
	case "KBANK通知":
		res = ExtractKBANKNotify(msg)
	case "KBANK读取":
		res = ExtractKBANKRead(msg)
	case "KBANK流水":
		res = ExtractKBANKRead(msg)
	case "BBL流水":
		res = ExtractBBLWater(msg)
	case "BBL":
		res = ExtractBBL(msg)
	case "BAAC":
		res = ExtractBAAC(msg)
	case "TTB", "TTB读取", "TTB通知":
		res = ExtractTTB(msg)
	default:
		res.Code = codes.ErrUnsupportedMessageType
	}
	return
}
