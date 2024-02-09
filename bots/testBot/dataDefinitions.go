package main

type BotConfig struct {
	ApiKey    string `xml:"ApiKey"`
	SecretKey string `xml:"SecretKey"`
	BaseURL   string `xml:"BaseURL"`
	BaseFURL  string `xml:"BaseFURL"`

	PairSymbol               string  `xml:"PairSymbol"`
	Symbol                   string  `xml:"Symbol"`
	TradeAmount              float64 `xml:"TradeAmount"`
	TakeProfitDelta          float64 `xml:"TakeProfitDelta"`
	StopLossDelta            float64 `xml:"StopLossDelta"`
	PriceDropBeforeBuy       float64 `xml:"PriceDropBeforeBuy"`
	CoolDownTimeoutBeforeBuy float64 `xml:"CoolDownTimeoutBeforeBuy"`
	Leverage                 float64 `xml:"Leverage"`
	FilePath                 string  `xml:"FilePath"`
}
