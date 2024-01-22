package binance_futures_connector

import (
	"context"
	"encoding/json"
	"net/http"
)

// Binance Symbol Price Ticker (GET /fapi/v2/ticker/price)
type TickerPrice struct {
	c      *Client
	symbol *string
}

// Symbol set symbol
func (s *TickerPrice) Symbol(symbol string) *TickerPrice {
	s.symbol = &symbol
	return s
}

// Send the request
func (s *TickerPrice) Do(ctx context.Context, opts ...RequestOption) (res *TickerPriceResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/ticker/price",
		secType:  secTypeNone,
	}
	if s.symbol != nil {
		r.setParam("symbol", *s.symbol)
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	res = new(TickerPriceResponse)
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Define TickerPrice response data
type TickerPriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
