package fiat

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/cg917658910/fzkj-wallet/quote-service/config"
)

var (
	_hot_symbols = []string{
		"USDT",
	}
	_hot_fiats = []string{
		"CNY",
	}
	_baseSleepTime   int64 = 30
	_randomSleepTime int64 = 10
	_minSleepTime    int64 = 30
)

func init() {
	logger.Infof("FiatQuote config: %+v", config.Configs.FiatQuote)
	_baseSleepTime = max(config.Configs.FiatQuote.UpdateCacheBaseTime, _minSleepTime)
	_randomSleepTime = config.Configs.FiatQuote.UpdateCacheRandomTime
	_hot_symbols = strings.Split(config.Configs.FiatQuote.UpdateCacheHotSymbols, ",")
	_hot_fiats = strings.Split(config.Configs.FiatQuote.UpdateCacheHotFiats, ",")
	startCacheUpdater()
}

func startCacheUpdater() {
	if quoteCache == nil {
		logger.Error("quoteCache is nil, updater cannot start")
		startQuoteCache()
	}
	if !config.Configs.FiatQuote.UpdateCacheEnabled {
		logger.Infof("quote cache updater is disabled, skip starting updater")
		return
	}
	ctx := context.Background()
	go func() {
		for {
			sleepTime := rand.Int63n(_randomSleepTime) + _baseSleepTime // Random sleep time between 5 and 15 seconds
			logger.Infof("start update cache|sleepTime=%d", sleepTime)
			for _, fiat := range _hot_fiats {
				for _, symbol := range _hot_symbols {
					params := &QuoteParams{
						Symbol: symbol,
						Fiat:   fiat,
					}
					// TODO：go routine fetch
					res := fetchQuote(ctx, params)
					if res.Error != nil {
						logger.Errorf("fetchQuoteAndCache error|symbol=%s|fiat=%s|error=%v", symbol, fiat, res.Error)
						continue
					}
					logger.Infof("update cache|symbol=%s|fiat=%s|buyPrice=%v|sellPrice=%v", res.Symbol, res.Fiat, res.BuyPrice, res.SellPrice)
					quoteCache.Set(res)
				}
			}
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}()
}
