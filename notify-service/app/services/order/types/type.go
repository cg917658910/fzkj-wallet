package types

import (
	"github.com/IBM/sarama"
	"github.com/cg917658910/fzkj-wallet/notify-service/app/services/order/enum"
	"github.com/cg917658910/fzkj-wallet/notify-service/lib/codes"
)

type NotifyMessage struct {
	RawMsg *sarama.ConsumerMessage
	Data   struct {
		Info      map[string]any `json:"info"`
		NotifyUrl string         `json:"notify_url"`
		OrderType string         `json:"order_type"`
		DataId    string         `json:"data_id"`
	} `json:"data"`
	Platform string `json:"platform"` // 平台
}

func (d NotifyMessage) Check() error {
	// TODO: Info 需要检查是否是合法的json
	if d.Data.NotifyUrl == "" {
		// TODO: 这里需要添加一个检查，是否是合法的URL
		// 目前先简单判断是否为空
		return codes.ErrInvalidArgument.Newf("notify url is empty")
	}
	if d.Data.OrderType == "" {
		return codes.ErrInvalidArgument.Newf("order type is empty")
	}
	if d.Data.DataId == "" {
		return codes.ErrInvalidArgument.Newf("data id is empty")
	}
	return nil
}

type NotifyTask struct {
	RetryCount int
	*NotifyMessage
}

type NotifyResult struct {
	*NotifyTask
	Result      string                  `json:"result"`       //通知返回内容
	Status      enum.NotifyResultStatus `json:"status"`       //通知状态 1 成功 2 重试多次无响应 3 无效请求地址
	Msg         string                  `json:"msg"`          //通知结果描述
	RequestTime string                  `json:"request_time"` //最后通知时间
}

type NotifyRequestParams struct {
	NotifyUrl  string         `json:"notify_url"`  //通知地址
	NotifyData map[string]any `json:"notify_data"` //通知数据
}

type NotifyResponse struct {
	HTTPStatus int    `json:"http_status"` // http状态码
	Body       string `json:"body"`        // 响应体
}

type MarkMessageParams struct {
	Msg      *sarama.ConsumerMessage
	MetaData string
}
