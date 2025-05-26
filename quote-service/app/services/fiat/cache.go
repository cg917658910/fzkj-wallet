package fiat

import (
	"fmt"
	"strings"
	"sync"
)

var quoteCache *Cache

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
	return res
}
