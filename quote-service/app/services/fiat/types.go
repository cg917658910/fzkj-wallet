package fiat

type (
	QuoteParams struct {
		Symbol string `json:"symbol"`
		Fiat   string `json:"fiat"`
	}

	QuoteResult struct {
		*QuoteParams
		SellPrice float32
		BuyPrice  float32
	}
)
