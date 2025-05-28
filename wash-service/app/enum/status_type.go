package enum

import "github.com/cg917658910/fzkj-wallet/wash-service/lib/codes"

type HealthStatus uint8

const (
	HealthStatus_SERVING HealthStatus = iota
)

var (
	healthStatusMap = map[HealthStatus]string{
		HealthStatus_SERVING: "serving",
	}
	healthStatuStringMap = map[string]HealthStatus{}
)

func init() {
	for key, val := range healthStatusMap {
		healthStatuStringMap[val] = key
	}
}

func (tp HealthStatus) String() string {
	return healthStatusMap[tp]
}

func HealthStatusFromString(tp string) (HealthStatus, error) {
	statusType, ok := healthStatuStringMap[tp]
	if !ok {
		return HealthStatus_SERVING, codes.ErrInvalidArgument.Newf("invalid health status: %s", tp)
	}
	return statusType, nil
}
