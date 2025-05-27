package fiat

import (
	"context"
)

func Quote(ctx context.Context, params *QuoteParams) (result *QuoteResult) {

	return fetchQuoteAndCache(ctx, params)
}
