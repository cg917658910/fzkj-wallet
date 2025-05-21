package enum

import "github.com/cg917658910/fzkj-wallet/notify-service/lib/codes"

type NotifyResultStatus uint8

const (
	NotifyResultStatusSuccessed NotifyResultStatus = iota
	NotifyResultStatusFailedInvalidParams
	NotifyResultStatusFailedMaxRetry
)

var (
	healthStatusMap = map[NotifyResultStatus]string{
		NotifyResultStatusSuccessed:           "successed",
		NotifyResultStatusFailedInvalidParams: "failed invalid params",
		NotifyResultStatusFailedMaxRetry:      "failed max retry",
	}
	healthStatuStringMap = map[string]NotifyResultStatus{}
)

func init() {
	for key, val := range healthStatusMap {
		healthStatuStringMap[val] = key
	}
}

func (tp NotifyResultStatus) String() string {
	return healthStatusMap[tp]
}

func NotifyResultStatusFromString(tp string) (NotifyResultStatus, error) {
	statusType, ok := healthStatuStringMap[tp]
	if !ok {
		return NotifyResultStatusSuccessed, codes.ErrInvalidArgument.Newf("invalid health status: %s", tp)
	}
	return statusType, nil
}
