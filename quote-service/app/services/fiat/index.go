package fiat

import "context"

func Quote(ctx context.Context, params *QuoteParams) (result *QuoteResult, err error) {
	result = &QuoteResult{
		QuoteParams: params,
		SellPrice:   123.00,
		BuyPrice:    124.00,
	}
	return
}
