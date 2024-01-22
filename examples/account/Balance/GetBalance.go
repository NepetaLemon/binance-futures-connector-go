package main

import (
	"context"
	"fmt"

	binance_futures_connector "github.com/binance/binance-futures-connector-go"
)

func main() {
	NewOrder()
}

func NewOrder() {
	apiKey := ""
	secretKey := ""
	baseURL := "https://fapi.binance.com"

	client := binance_futures_connector.NewClient(apiKey, secretKey, baseURL)

	// Binance New Order endpoint - POST /fapi/v2/balance
	balance, err := client.NewBalanceService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binance_futures_connector.PrettyPrint(balance))
}
