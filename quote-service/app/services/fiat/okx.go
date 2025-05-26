package fiat

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OKXProvider = "okx"
)

func fetchOkxQuote(ctx context.Context, asset, fiat, side string) (quoteResp *QuoteResponse) {
	quoteResp = &QuoteResponse{
		Platform: OKXProvider, Asset: asset, Fiat: fiat, Side: side,
	}
	url := fmt.Sprintf("https://www.okx.com/priapi/v3/b2c/deposit/quotedPrice?side=%s&quoteCurrency=%s&baseCurrency=%s", side, fiat, asset)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		quoteResp.Error = err
		return
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	type Asset struct {
		Price string `json:"price"`
	}
	var r struct {
		Code         int      `json:"code"`
		Data         []*Asset `json:"data"`
		ErrorMessage string   `json:"error_message"`
		ErrorCode    string   `json:"error_code"`
	}
	if err := json.Unmarshal(raw, &r); err != nil {
		quoteResp.Error = fmt.Errorf("parse error: %v", err)
		return
	}
	if len(r.Data) == 0 {
		quoteResp.Error = fmt.Errorf("Quote assets empty: %s", r.ErrorMessage)
		return
	}
	fmt.Println("r.Data[0]: ", r.Data[0])
	fmt.Sscanf(r.Data[0].Price, "%f", &quoteResp.Price)
	return
}
