package producer

import "github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/errors"

// genError 用于生成错误值。
func genError(errMsg string) error {
	return errors.NewNotifyError(errors.ERROR_TYPE_PRODUCER,
		errMsg)
}
