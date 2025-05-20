package data

import "github.com/cg917658910/fzkj-wallet/notify-service/lib/codes"

type CallResult struct {
	RawData     string `json:"raw_data"`     //原始数据
	Response    string `json:"response"`     //通知响应返回
	Status      int32  `json:"status"`       //通知状态 1 成功 2 重试多次无响应 3 无效请求地址
	Msg         string `json:"msg"`          //通知结果描述
	RequestTime string `json:"request_time"` //最后通知时间
}

type CallRequsetParams struct {
	NotifyUrl  string                 `json:"notify_url"`  //通知地址
	NotifyData map[string]interface{} `json:"notify_data"` //通知数据
}

type OrderNotifyMessage struct {
	Data struct {
		//Info      map[string]any `json:"info"`
		NotifyUrl string `json:"notify_url"`
		OrderType string `json:"order_type"`
		DataId    string `json:"data_id"`
	} `json:"data"`
}

func (d OrderNotifyMessage) Check() error {
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
