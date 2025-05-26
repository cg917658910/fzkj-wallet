package fiat

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var quoteCache *Cache
var (
	_validCacheTime = 60 // seconds
)

func init() {
	startQuoteCache()
}

func startQuoteCache() {
	quoteCache = NewCache()
}

type Cache struct {
	mu           sync.RWMutex
	quoteResults map[string]*QuoteResult
}

func NewCache() *Cache {
	return &Cache{
		quoteResults: make(map[string]*QuoteResult),
	}
}

func (c *Cache) buildKey(symbol, fiat string) string {
	return strings.ToLower(fmt.Sprintf("%s_%s_", symbol, fiat))

}

func (c *Cache) Set(q *QuoteResult) {
	if q == nil {
		return
	}
	k := c.buildKey(q.Symbol, q.Fiat)
	c.mu.Lock()
	defer c.mu.Unlock()
	q.CacheTime = time.Now()
	c.quoteResults[k] = q
}

func (c *Cache) Get(symbol, fiat string) *QuoteResult {
	k := c.buildKey(symbol, fiat)
	c.mu.RLock()
	defer c.mu.RUnlock()
	res, ok := c.quoteResults[k]
	if !ok {
		return nil
	}
	if time.Since(res.CacheTime).Seconds() > float64(_validCacheTime) {
		logger.Infof("Cache expired|symbol=%s|fiat=%s|cacheTime=%s", symbol, fiat, res.CacheTime)
		// Cache expired
		delete(c.quoteResults, k)
		return nil
	}
	return res
}
