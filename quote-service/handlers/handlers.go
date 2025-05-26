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
	result, err := fiat.Quote(ctx, params)
	if err != nil {
		return nil, err
	}
	resp = pb.FiatQuoteResponse{
		Data: &pb.FiatQuoteResult{},
	}
	resp.Code = result.Code
	resp.Msg = "success"
	if result.Error != nil {
		resp.Msg = result.Error.Error()
	}
	resp.Data.Symbol = result.Symbol
	resp.Data.Fiat = result.Fiat
	resp.Data.SellPrice = result.SellPrice
	resp.Data.BuyPrice = result.BuyPrice
	return &resp, nil
}
