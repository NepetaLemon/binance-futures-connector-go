package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"time"

	binance_futures_connector "github.com/binance/binance-futures-connector-go"
)

var apiKey = ""
var secretKey = ""
var fbaseURL = ""
var lastSaveTime time.Time
var botConfig BotConfig

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func main() {
	if len(os.Args) > 1 && os.Args[2] != "" {
		initLogger()
		botConfig = readConfig(os.Args[2])
		apiKey = botConfig.ApiKey
		secretKey = botConfig.SecretKey
		fbaseURL = botConfig.BaseFURL

		//GetCurrentOpenOrders()
		StartBot()

	} else {
		ErrorLogger.Println("Missing config argument")
	}
}

func StartBot() {

	// var lastPrice float64 = 0
	var lastWalletAmmout float64 = 0

	for {
		walletAmmount, err := GetWalletAmount("USDT")

		if err != nil {
			ErrorLogger.Println(err.Error())
			return
		}
		if walletAmmount > 0 {
			//cPrice, err := LastPrice(botConfig.PairSymbol)
			if err != nil {
				ErrorLogger.Println(err.Error())
				return
			}
			if walletAmmount != lastWalletAmmout {
				InfoLogger.Println("Wallet Amount:", walletAmmount)
				lastWalletAmmout = walletAmmount
			}

			var side = "BUY"
			var positionSide = "LONG"
			OpenPosition(botConfig.PairSymbol, botConfig.TradeAmount, botConfig.ProfitPriceDelta, side, positionSide)
			// lastPrice = cPrice

		}
		fileInfo, _ := os.Stat(botConfig.FilePath)
		cSavedTime := fileInfo.ModTime()
		if cSavedTime != lastSaveTime {
			botConfig = readConfig(os.Args[2])
			lastSaveTime = cSavedTime
		}

		time.Sleep(60 * time.Second)
	}
}

func OpenPosition(pairSymbol string, quantity, priceDelta float64, side, positionSide string) bool {

	client := binance_futures_connector.NewClient(apiKey, secretKey, fbaseURL)
	// Create new order

	// newOrder, err := client.NewCreateOrderService().Symbol(pairSymbol).Side(side).PositionSide(positionSide).Type("MARKET").
	// 	Quantity(quantity).
	// 	Do(context.Background())
	cLastPrice, err := LastPrice(pairSymbol)
	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	activationPrice := round(cLastPrice/priceDelta, 2)
	SLPrice := round(cLastPrice/priceDelta, 2)
	TPPrice := round(cLastPrice*priceDelta, 2)

	newOrder, err := client.NewCreateOrderService().Symbol(pairSymbol).PositionSide(positionSide).
		Side("BUY").Type("TRAILING_STOP_MARKET").Quantity(quantity).ActivationPrice(activationPrice).CallbackRate(0.1).
		Do(context.Background())
	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	s, _ := json.Marshal(newOrder)
	InfoLogger.Println("New Open Position:", string(s))
	//fmt.Println(binance_connector.PrettyPrint(newOrder))
	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}

	newSellOrder, err := client.NewCreateOrderService().Symbol(pairSymbol).PositionSide(positionSide).
		Side("SELL").Type("TAKE_PROFIT_MARKET").Quantity(quantity).StopPrice(SLPrice).TimeInForce("GTC").ClosePosition("true").
		Do(context.Background())

	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	// fmt.Println(binance_connector.PrettyPrint(newSellOrder))
	s, _ = json.Marshal(newSellOrder)
	InfoLogger.Println("New Stop Lose Order:", string(s))
	newSellOrder, err = client.NewCreateOrderService().Symbol(pairSymbol).PositionSide(positionSide).
		Side("SELL").Type("STOP_MARKET").Quantity(quantity).StopPrice(TPPrice).ClosePosition("true").
		Do(context.Background())

	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	// fmt.Println(binance_connector.PrettyPrint(newSellOrder))
	s, _ = json.Marshal(newSellOrder)
	InfoLogger.Println("New Take Profit Order:", string(s))
	return true
}

func readConfig(filepath string) BotConfig {
	var config BotConfig
	file, err := os.Open(filepath)
	if err != nil {
		ErrorLogger.Println(err.Error())
		return config
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		ErrorLogger.Println(err.Error())
		return config
	}
	fileInfo, _ := os.Stat(filepath)
	lastSaveTime = fileInfo.ModTime()
	config.FilePath = filepath

	InfoLogger.Println("Reload Configuration")
	return config
}
func initLogger() {
	file, err := os.OpenFile("flogs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
