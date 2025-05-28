package types

import (
	"errors"

	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/codes"
)

type (
	WashRequest struct {
		Msg      string
		MsgId    string
		Currency string
		TypeName string
	}
	Extracted struct {
		PayCoin float64
		PayTime int64
		Balance float64
		codes.Code
	}
	WashResponse struct {
		*WashRequest
		*Extracted
	}
)

func (wreq WashRequest) Check() error {
	if wreq.Msg == "" {
		return errors.New("invalid msg")
	}
	if wreq.Currency == "" {
		return errors.New("invalid currency")
	}
	if wreq.TypeName == "" {
		return errors.New("invalid msg type")
	}
	return nil
}
