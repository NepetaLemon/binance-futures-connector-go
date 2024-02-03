package binance_futures_connector

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// Binance New Order endpoint (POST /fapi/v1/order)
//
// CreateOrderService create order
type CreateOrderService struct {
	c                *Client
	symbol           string
	side             string
	positionSide     *string
	orderType        string
	timeInForce      *string
	quantity         *float64
	reduceOnly       *string
	price            *float64
	newClientOrderID *string
	stopPrice        *float64
	closePosition    *string
	activationPrice  *float64
	callbackRate     *float64
	workingType      *string
	priceProtect     *string
	newOrderRespType *string
}

// Symbol set symbol
func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

// Side set side
func (s *CreateOrderService) Side(side string) *CreateOrderService {
	s.side = side
	return s
}

// Type set positionSide
func (s *CreateOrderService) PositionSide(positionSide string) *CreateOrderService {
	s.positionSide = &positionSide
	return s
}

// Type set type
func (s *CreateOrderService) Type(orderType string) *CreateOrderService {
	s.orderType = orderType
	return s
}

// TimeInForce set timeInForce
func (s *CreateOrderService) TimeInForce(timeInForce string) *CreateOrderService {
	s.timeInForce = &timeInForce
	return s
}

// Quantity set quantity
func (s *CreateOrderService) Quantity(quantity float64) *CreateOrderService {
	s.quantity = &quantity
	return s
}

// ReduceOnly set reduceOnly
func (s *CreateOrderService) ReduceOnly(ReduceOnly string) *CreateOrderService {
	s.reduceOnly = &ReduceOnly
	return s
}

// Price set price
func (s *CreateOrderService) Price(price float64) *CreateOrderService {
	s.price = &price
	return s
}

// NewClientOrderId set newClientOrderId
func (s *CreateOrderService) NewClientOrderId(newClientOrderId string) *CreateOrderService {
	s.newClientOrderID = &newClientOrderId
	return s
}

// StopPrice set stopPrice
func (s *CreateOrderService) StopPrice(stopPrice float64) *CreateOrderService {
	s.stopPrice = &stopPrice
	return s
}

// ClosePosition set closePosition
func (s *CreateOrderService) ClosePosition(closePosition string) *CreateOrderService {
	s.closePosition = &closePosition
	return s
}

// ActivationPrice set activationPrice
func (s *CreateOrderService) ActivationPrice(activationPrice float64) *CreateOrderService {
	s.activationPrice = &activationPrice
	return s
}

// CallbackRate set callbackRate
func (s *CreateOrderService) CallbackRate(callbackRate float64) *CreateOrderService {
	s.callbackRate = &callbackRate
	return s
}

// IcebergQuantity set icebergQuantity
func (s *CreateOrderService) WorkingType(workingType string) *CreateOrderService {
	s.workingType = &workingType
	return s
}

// PriceProtect set priceProtect
func (s *CreateOrderService) PriceProtect(priceProtect string) *CreateOrderService {
	s.priceProtect = &priceProtect
	return s
}

// NewOrderRespType set newOrderRespType
func (s *CreateOrderService) NewOrderRespType(newOrderRespType string) *CreateOrderService {
	s.newOrderRespType = &newOrderRespType
	return s
}

const (
	ACK    = 1
	RESULT = 2
)

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (res interface{}, err error) {
	respType := RESULT
	r := &request{
		method:   http.MethodPost,
		endpoint: "/fapi/v1/order",
		secType:  secTypeSigned,
	}
	r.setParam("symbol", s.symbol)
	r.setParam("side", s.side)
	r.setParam("type", s.orderType)
	switch s.orderType {
	case "MARKET":
		respType = RESULT
	case "LIMIT":
		respType = RESULT
	}
	if s.positionSide != nil {
		r.setParam("positionSide", *s.positionSide)
	}
	if s.timeInForce != nil {
		r.setParam("timeInForce", *s.timeInForce)
	}
	if s.quantity != nil {
		r.setParam("quantity", strconv.FormatFloat(*s.quantity, 'f', -1, 64))
	}
	if s.reduceOnly != nil {
		r.setParam("reduceOnly", *s.reduceOnly)
	}
	if s.price != nil {
		r.setParam("price", *s.price)
	}
	if s.newClientOrderID != nil {
		r.setParam("newClientOrderId", *s.newClientOrderID)
	}
	if s.stopPrice != nil {
		r.setParam("stopPrice", *s.stopPrice)
	}
	if s.closePosition != nil {
		r.setParam("closePosition", *s.closePosition)
	}
	if s.activationPrice != nil {
		r.setParam("activationPrice", *s.activationPrice)
	}
	if s.callbackRate != nil {
		r.setParam("callbackRate", *s.callbackRate)
	}
	if s.workingType != nil {
		r.setParam("workingType", *s.workingType)
	}
	if s.newOrderRespType != nil {
		r.setParam("newOrderRespType", *s.newOrderRespType)
		switch *s.newOrderRespType {
		case "ACK":
			respType = ACK
		case "RESULT":
			respType = RESULT
		}
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}
	switch respType {
	case RESULT: //case ACK - NOT HANDLED
		res = new(CreateOrderResponse)
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type CreateOrderResponse struct {
	Symbol                  string `json:"symbol"`
	OrderID                 int64  `json:"orderId"`
	ClientOrderID           string `json:"clientOrderId"`
	Price                   string `json:"price"`
	OrigQuantity            string `json:"origQty"`
	ExecutedQuantity        string `json:"executedQty"`
	CumQuote                string `json:"cumQuote"`
	ReduceOnly              bool   `json:"reduceOnly"`
	Status                  string `json:"status"`
	StopPrice               string `json:"stopPrice"` // please ignore when order type is TRAILING_STOP_MARKET
	TimeInForce             string `json:"timeInForce"`
	Type                    string `json:"type"`
	OrigType                string `json:"origType"`
	Side                    string `json:"side"`
	UpdateTime              int64  `json:"updateTime"`
	GoodTillDate            int64  `json:"goodTillDate"`
	WorkingType             string `json:"workingType"` //order pre-set auot cancel time for TIF GTD order
	ActivatePrice           string `json:"activatePrice"`
	PriceRate               string `json:"priceRate"`
	AvgPrice                string `json:"avgPrice"`
	PositionSide            string `json:"positionSide"`
	ClosePosition           bool   `json:"closePosition"`           // if Close-All
	PriceProtect            bool   `json:"priceProtect"`            // if conditional order trigger is protected
	PriceMatch              string `json:"priceMatch"`              //price match mode
	SelfTradePreventionMode string `json:"selfTradePreventionMode"` //self trading preventation mode
	RateLimitOrder10s       string `json:"rateLimitOrder10s,omitempty"`
	RateLimitOrder1m        string `json:"rateLimitOrder1m,omitempty"`
}

// Binance New Order endpoint (POST /fapi/v1/order)
//
// CreateOrderService create order
type FuturesAccountBalanceService struct {
	c *Client
}

// Do send request
func (s *FuturesAccountBalanceService) Do(ctx context.Context, opts ...RequestOption) (res []*FuturesAccountBalanceResponse, err error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v2/balance",
		secType:  secTypeSigned,
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return []*FuturesAccountBalanceResponse{}, err
	}

	res = make([]*FuturesAccountBalanceResponse, 0)
	err = json.Unmarshal(data, &res)
	if err != nil {
		return []*FuturesAccountBalanceResponse{}, err
	}
	return res, nil
}

// Balance define user balance of your account
type FuturesAccountBalanceResponse struct {
	AccountAlias       string `json:"accountAlias"`
	Asset              string `json:"asset"`
	Balance            string `json:"balance"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
	AvailableBalance   string `json:"availableBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
}
