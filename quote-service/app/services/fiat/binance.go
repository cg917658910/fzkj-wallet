package fiat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func fetchBinanceQuote(ctx context.Context, asset, fiat, side string) (quoteResp *QuoteResponse) {
	quoteResp = &QuoteResponse{
		Platform: "binance", Asset: asset, Fiat: fiat, Side: side,
	}
	url := "https://p2p.binance.com/bapi/c2c/v2/public/c2c/adv/quoted-price"

	body := map[string]any{
		"assets":       []string{asset},
		"fiatCurrency": fiat,
		"tradeType":    side,
		"fromUserRole": "USER",
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		quoteResp.Error = err
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	type Asset struct {
		ReferencePrice float64 `json:"referencePrice"`
	}
	var r struct {
		Code    string   `json:"code"`
		Data    []*Asset `json:"data"`
		Message string   `json:"message"`
		Success bool     `json:"success"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		quoteResp.Error = fmt.Errorf("parse error: %v", err)
		return
	}
	if !r.Success {
		quoteResp.Error = fmt.Errorf("Quote Binance error: %s", r.Message)
		return
	}
	if len(r.Data) == 0 {
		quoteResp.Error = fmt.Errorf("Quote Binance assets empty: %s", r.Message)
		return
	}
	quoteResp.Price = r.Data[0].ReferencePrice
	return
}
