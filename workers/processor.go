package workers

import (
	"github.com/sirupsen/logrus"
	"gobot/exchanges"
	"time"
)

type Processor struct {
	Log           *logrus.Logger
	SocketBinance *exchanges.SocketBinance
}

type CandleResponse struct {
	EventType string  `json:"e"`
	EventTime int     `json:"E"`
	Symbol    string  `json:"s"`
	Candle    *Candle `json:"k"`
}

type Candle struct {
	KlineStartTime           int64  `json:"t"`
	KlineCloseTime           int64  `json:"T"`
	Symbol                   string `json:"s"`
	Interval                 string `json:"i"`
	FirstTradeID             int    `json:"f"`
	LastTradeID              int    `json:"L"`
	OpenPrice                string `json:"o"`
	ClosePrice               string `json:"c"`
	HighPrice                string `json:"h"`
	LowPrice                 string `json:"l"`
	BaseAssetVolume          string `json:"v"`
	NumberOfTrades           int    `json:"n"`
	IsThisKlineClosed        bool   `json:"x"`
	QuoteAssetVolume         string `json:"q"`
	TakerBuyBaseAssetVolume  string `json:"V"`
	TakerBuyQuoteAssetVolume string `json:"Q"`
	Ignore                   string `json:"B"`
}

func NewProcessor(socketBinance *exchanges.SocketBinance, log *logrus.Logger) *Processor {
	return &Processor{
		Log:           log,
		SocketBinance: socketBinance,
	}
}

func (p *Processor) Reading() {
	defer p.SocketBinance.Close()
	for {
		response := CandleResponse{}
		if p.SocketBinance.IsConnected && !p.SocketBinance.Terminate {
			p.SocketBinance.Connection.SetReadDeadline(time.Now().Add(time.Minute))
			err := p.SocketBinance.Connection.ReadJSON(&response)
			//_, message, err := p.SocketBinance.Connection.ReadMessage()
			if err != nil {
				p.Log.Errorf("Error read JSON from web-socket: %s", err.Error())
				p.SocketBinance.IsConnected = false
			} else {
				p.ParseRequest(&response)
			}
		} else if !p.SocketBinance.Terminate {
			p.SocketBinance.Connect()
		}
	}
}

func (p *Processor) ParseRequest(response *CandleResponse) {
	var isClosed string
	if response.Candle.IsThisKlineClosed {
		isClosed = "yep"
	} else {
		isClosed = "nope"
	}
	p.Log.Infof("[%s] Candle: %s-%s, OpenPrice: %s, ClosePrice: %s, isClosed: %s",
		time.Now().Format("15:04:05"),
		// todo check and fix
		time.Unix(response.Candle.KlineStartTime, 100).Format("15:04:05"),
		time.Unix(response.Candle.KlineCloseTime, 100).Format("15:04:05"),
		response.Candle.OpenPrice,
		response.Candle.ClosePrice,
		isClosed,
	)
}
