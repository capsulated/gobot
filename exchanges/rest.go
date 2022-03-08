package exchanges

import (
	"io/ioutil"
	"log"
	"net/http"
)

//[
//	[
//		1499040000000,      // Open time
//		"0.01634790",       // Open
//		"0.80000000",       // High
//		"0.01575800",       // Low
//		"0.01577100",       // Close
//		"148976.11427815",  // Volume
//		1499644799999,      // Close time
//		"2434.19055334",    // Quote asset volume
//		308,                // Number of trades
//		"1756.87402397",    // Taker buy base asset volume
//		"28.46694368",      // Taker buy quote asset volume
//		"17928899.62484339" // Ignore.
//	]
//]

type Candles [][]interface{}

func (s *SocketBinance) RestGet2hoursPairData() {
	resp, err := http.Get("https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=10m")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string

	// StdDev = SQRT (SUM ((CLOSE – SMA (CLOSE, N))^2, N)/N)
	sb := string(body)
	log.Printf(sb)
}
