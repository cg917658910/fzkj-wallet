package fiat

import (
	"context"
	"errors"
	"sync"
)

func fetchQuoteAndCache(ctx context.Context, params *QuoteParams) (res *QuoteResult, err error) {
	if err := params.Check(); err != nil {
		return nil, err
	}
	// 查询缓存
	cacheResult := quoteCache.Get(params.Symbol, params.Fiat)
	if cacheResult != nil {
		logger.Infof("fetchQuote use cache|asset=%s|fiat=%s|buyPrice=%v|sellPrice=%v", params.Symbol, params.Fiat, cacheResult.BuyPrice, cacheResult.SellPrice)
		return cacheResult, nil
	}
	res, err = fetchQuote(ctx, params)
	if err != nil {
		return nil, err
	}
	// set cache
	quoteCache.Set(res)
	return res, nil
}

func fetchQuote(ctx context.Context, params *QuoteParams) (res *QuoteResult, err error) {
	if params == nil {
		err = errors.New("params is nil")
		return
	}
	var (
		wg     sync.WaitGroup
		symbol = params.Symbol
		fiat   = params.Fiat
	)
	results := make(chan *QuoteResponse, 4)
	for _, side := range []string{"BUY", "SELL"} {
		wg.Add(2)
		go func(symbol, side string) {
			defer wg.Done()
			logger.Infof("send fetchBinanceQuote request|asset=%s|fiat=%s|side=%s", symbol, fiat, side)
			result := fetchBinanceQuote(ctx, symbol, fiat, side)
			if result.Error != nil {
				result.ErrorMsg = result.Error.Error()
			}
			results <- result
			logger.Infof("fetchBinanceQuote result|asset=%s|fiat=%s|side=%s|errorMsg=%s", symbol, fiat, side, result.ErrorMsg)

		}(symbol, side)
		// okx
		go func(symbol, side string) {
			defer wg.Done()
			logger.Infof("send fetchOkxQuote request|asset=%s|fiat=%s|side=%s", symbol, fiat, side)
			result := fetchOkxQuote(ctx, symbol, fiat, side)
			if result.Error != nil {
				result.ErrorMsg = result.Error.Error()
			}
			results <- result
			logger.Infof("fetchOkxQuote result|asset=%s|fiat=%s|side=%s|errorMsg=%s", symbol, fiat, side, result.ErrorMsg)

		}(symbol, side)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var buyQuotes []*QuoteResponse
	var sellQuotes []*QuoteResponse
	for q := range results {
		if q.Error != nil {
			logger.Errorf("Fetch quote failed|platform=%s|error=%s", q.Platform, q.Error)
			continue
		}
		if q.Side == "BUY" {
			buyQuotes = append(buyQuotes, q)
			continue
		}
		sellQuotes = append(sellQuotes, q)
	}
	res = &QuoteResult{}
	res.QuoteParams = params
	res.BuyCode = 1001 // not quote
	res.SellCode = 1001
	if len(buyQuotes) > 0 {
		res.BuyCode = 0
		res.BuyPrice = float32(buyQuotes[0].Price)
	}
	if len(sellQuotes) > 0 {
		res.SellCode = 0
		res.SellPrice = float32(sellQuotes[0].Price)
	}

	return
}
