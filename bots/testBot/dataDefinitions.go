package main

type BotConfig struct {
	ApiKey    string `xml:"ApiKey"`
	SecretKey string `xml:"SecretKey"`
	BaseURL   string `xml:"BaseURL"`
	BaseFURL  string `xml:"BaseFURL"`

	PairSymbol       string  `xml:"PairSymbol"`
	Symbol           string  `xml:"Symbol"`
	TradeAmount      float64 `xml:"TradeAmount"`
	ProfitPriceDelta float64 `xml:"ProfitPriceDelta"`
	FilePath         string  `xml:"FilePath"`
}
