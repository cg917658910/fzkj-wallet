package handlers

import (
	"context"

	"github.com/cg917658910/fzkj-wallet/quote-service/app/services/fiat"
	pb "github.com/cg917658910/fzkj-wallet/quote-service/proto"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.QuoteServiceServer {
	return quoteserviceService{}
}

type quoteserviceService struct{}

func (s quoteserviceService) HealthCheck(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	var resp pb.HealthCheckResponse
	return &resp, nil
}

func (s quoteserviceService) FiatQuote(ctx context.Context, in *pb.FiatQuoteRequest) (*pb.FiatQuoteResponse, error) {
	var resp pb.FiatQuoteResponse
	params := &fiat.QuoteParams{
		Symbol: in.Symbol,
		Fiat:   in.Fiat,
	}
	resp = pb.FiatQuoteResponse{
		Data: &pb.FiatQuoteResult{},
		Msg:  "success",
		Code: 1001,
	}
	resp.Data.Symbol = params.Symbol
	resp.Data.Fiat = params.Fiat
	result := fiat.Quote(ctx, params)
	if result != nil {
		resp.Code = result.Code
		if result.Error != nil {
			resp.Msg = result.Error.Error()
		}
		resp.Data.SellPrice = result.SellPrice
		resp.Data.BuyPrice = result.BuyPrice
	}

	return &resp, nil
}
