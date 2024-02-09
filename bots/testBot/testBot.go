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
	var lastWalletAmmout float64 = 0
	var waitTime int64 = 15

	for {
		var quantity = botConfig.TradeAmount
		walletAmmount, err := GetAvailableBalance("USDT")

		if err != nil {
			ErrorLogger.Println(err.Error())
			return
		}
		if walletAmmount > 0 {
			cPrice, err := LastPrice(botConfig.PairSymbol)
			if err != nil {
				ErrorLogger.Println(err.Error())
				return
			}
			if walletAmmount != lastWalletAmmout {
				InfoLogger.Println("Wallet Amount:", walletAmmount)
				lastWalletAmmout = walletAmmount
			}
			if quantity < 0 {
				for inc := 0.0; cPrice*inc/botConfig.Leverage <= walletAmmount; inc = inc + 0.001 {
					quantity = round(inc, 3)
				}
				quantity -= 0.001
				quantity = round(quantity, 3)
			} else {
				if cPrice*botConfig.TradeAmount/botConfig.Leverage < walletAmmount {
					quantity = botConfig.TradeAmount
				} else {
					quantity = -1
				}
			}
			if quantity > 0 {
				waitForCooldown(cPrice)
				var side = "BUY"
				var positionSide = "LONG"
				OpenPosition(botConfig.PairSymbol, quantity, botConfig.StopLossDelta, botConfig.TakeProfitDelta, side, positionSide)
				waitTime = 60
			}

		}
		fileInfo, _ := os.Stat(botConfig.FilePath)
		cSavedTime := fileInfo.ModTime()
		if cSavedTime != lastSaveTime {
			botConfig = readConfig(os.Args[2])
			lastSaveTime = cSavedTime
		}

		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}

func waitForCooldown(buyPrice float64) {
	var lastRoundPrice float64

	startTime := time.Now()
	lastRoundTime := time.Now()

	cooldownTimeout := time.Duration(botConfig.CoolDownTimeoutBeforeBuy) * time.Second

	for time.Since(startTime) < cooldownTimeout {
		lastPrice, err := LastPrice(botConfig.PairSymbol)
		if err != nil {
			ErrorLogger.Println(err.Error())
			return
		}
		priceDrop := buyPrice - lastPrice
		if priceDrop > botConfig.PriceDropBeforeBuy && lastRoundPrice < lastPrice {
			return
		}
		lastRoundPrice = lastPrice
		if time.Since(lastRoundTime) > 60*time.Second {
			lastRoundTime = time.Now()
			buyPrice = lastPrice
		}
		time.Sleep(10 * time.Second)
	}
}

func OpenPosition(pairSymbol string, quantity, StopLossDelta, TakeProfitDelta float64, side, positionSide string) bool {

	client := binance_futures_connector.NewClient(apiKey, secretKey, fbaseURL)

	err := client.CancelAllOpenOrdersService().Symbol(pairSymbol).Do(context.Background())
	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	cLastPrice, err := LastPrice(pairSymbol)
	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	//activationPrice := round(cLastPrice/priceDelta, 2)
	SLPrice := round(cLastPrice/StopLossDelta, 2)
	TPPrice := round(cLastPrice*TakeProfitDelta, 2)

	// Open new Position
	newOrder, err := client.NewCreateOrderService().Symbol(pairSymbol).Side(side).PositionSide(positionSide).
		Type("MARKET").Quantity(quantity).
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
		Side("SELL").Type("TAKE_PROFIT_MARKET").Quantity(quantity).StopPrice(TPPrice).
		Do(context.Background())

	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	// fmt.Println(binance_connector.PrettyPrint(newSellOrder))
	s, _ = json.Marshal(newSellOrder)
	InfoLogger.Println("New Take Profit Order:", string(s))
	newSellOrder, err = client.NewCreateOrderService().Symbol(pairSymbol).PositionSide(positionSide).
		Side("SELL").Type("STOP_MARKET").Quantity(quantity).StopPrice(SLPrice).
		Do(context.Background())

	if err != nil {
		ErrorLogger.Println(err.Error())
		return false
	}
	// fmt.Println(binance_connector.PrettyPrint(newSellOrder))
	s, _ = json.Marshal(newSellOrder)
	InfoLogger.Println("New Stop Lose Order:", string(s))
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
