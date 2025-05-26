package fiat

import (
	"time"

	"github.com/cg917658910/fzkj-wallet/quote-service/lib/codes"
)

type (
	QuoteParams struct {
		Symbol string `json:"symbol"`
		Fiat   string `json:"fiat"`
	}

	QuoteResult struct {
		*QuoteParams
		SellPrice float32
		BuyPrice  float32
		ErrorMsg  string
		SellCode  int
		BuyCode   int
		Code      int32
		Error     error
		CacheTime time.Time
	}
	QuoteResponse struct {
		Platform string  `json:"platform"`
		Asset    string  `json:"asset"`
		Fiat     string  `json:"fiat"`
		Side     string  `json:"side"`
		Price    float64 `json:"price"`
		Time     string  `json:"time"`
		Error    error   `json:"error,omitempty"`
		ErrorMsg string  `json:"error_msg,omitempty"`
	}
)

func (q QuoteParams) Check() error {
	if q.Symbol == "" {
		return codes.ErrInvalidArgument
	}
	if q.Fiat == "" {
		return codes.ErrInvalidArgument
	}

	return nil
}

const (
	notFoundQuote    = "no valid quotes found"
	ErrNotFoundQuote = 1002
)
