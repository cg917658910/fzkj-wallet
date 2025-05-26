package fiat

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/quote-service/lib/codes"
)

func Quote(ctx context.Context, params *QuoteParams) (result *QuoteResult, err error) {
	if params == nil {
		err = codes.ErrInvalidArgument.New("params is nil")
		return
	}
	if err = params.Check(); err != nil {
		err = codes.ErrInvalidArgument.New(err.Error())
		return
	}
	return fetchQuoteAndCache(ctx, params)
}
