package main

import (
	"context"
	"math"
	"strconv"

	binance_futures_connector "github.com/binance/binance-futures-connector-go"
)

func GetWalletAmount(symbol string) (float64, error) {
	client := binance_futures_connector.NewClient(apiKey, secretKey, baseURL)
	asset, err := client.NewBalanceService().
		Do(context.Background())

	if err != nil {
		ErrorLogger.Println(err.Error())
		return -1, err
	}

	//fmt.Println(binance_connector.PrettyPrint(asset))

	retVal, err := strconv.ParseFloat(asset[len(asset)-1].Free, 64)
	if err != nil {
		ErrorLogger.Println(err.Error())
		return retVal, err
	}

	//InfoLogger.Println("Wallet Amount:", retVal)
	return retVal, nil
}

func LastPrice(symbol string) (float64, error) {

	client := binance_futures_connector.NewClient("", "", baseURL)

	// AvgPrice
	lastPrice, err := client.NewTickerService().
		Symbol(symbol).Do(context.Background())
	if err != nil {
		ErrorLogger.Println(err.Error())
		return 0, err
	}

	fLastPrice, err := strconv.ParseFloat(lastPrice.LastPrice, 64)
	if err != nil {
		ErrorLogger.Println(err.Error())
		return 0, err
	}

	InfoLogger.Println("LastPrice:", fLastPrice)
	return fLastPrice, nil
}

// round rounds a float64 to a specified number of decimal places
func round(f float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(f*shift) / shift
}

func GetCurrentOpenOrders() {

	client := binance_futures_connector.NewClient(apiKey, secretKey, baseURL)

	// Binance Get current open orders - GET /api/v3/openOrders
	getCurrentOpenOrders, err := client.NewGetOpenOrdersService().Symbol("BTCUSDT").
		Do(context.Background())
	if err != nil {
		ErrorLogger.Println(err.Error())
		return
	}
	//fmt.Println(binance_connector.PrettyPrint(getCurrentOpenOrders))
	InfoLogger.Println("Orders", getCurrentOpenOrders)
}
