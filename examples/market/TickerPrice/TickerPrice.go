package main

import (
	"context"
	"fmt"

	binance_futures_connector "github.com/binance/binance-futures-connector-go"
)

func main() {
	TickerPrice()
}

func TickerPrice() {
	baseURL := "https://fapi.binance.com"

	client := binance_futures_connector.NewClient("", "", baseURL)

	// TickerPrice
	tickerPrice, err := client.NewTickerPriceService().Symbol("BTCUSDT").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(binance_futures_connector.PrettyPrint(tickerPrice))
}
