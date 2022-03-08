package exchanges

import (
	"crypto/tls"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gobot/config"
	"net/url"
)

type SocketBinance struct {
	IsConnected bool
	Log         *logrus.Logger
	Dialer      websocket.Dialer
	WebAddress  string
	Connection  *websocket.Conn
	Tls         bool
	Terminate   bool
}

type Request struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

type Response struct {
	Result         string `json:"result"`
	Connect        string `json:"connect"`
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	LastName       string `json:"last_name"`
	CurrencySymbol string `json:"currency_symbol"`
	Demo           string `json:"demo"`
	Error          string `json:"error"`
}

func NewSocketBinance(log *logrus.Logger, config *config.Config) *SocketBinance {
	binanceClient := &SocketBinance{
		IsConnected: false,
		Log:         log,
	}

	binanceClient.Dialer = websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	var host = config.BinanceSocket.Address + ":" + config.BinanceSocket.Port
	u := url.URL{
		Scheme: "wss",
		Host:   host,
		Path:   config.BinanceSocket.Path,
	}

	binanceClient.WebAddress = u.String()
	binanceClient.Tls = true

	binanceClient.Terminate = false
	return binanceClient
}

func (s *SocketBinance) Connect() error {
	s.Log.Info("Dialing to socket...")
	err := s.Dial()
	if err != nil {
		s.Log.Errorf("Dial to socket %s err: %s", s.WebAddress, err)
		return err
	}
	s.Log.Info("OK dialing to socket successful!")

	//s.Log.Info("Handshaking with socket...")
	//err = s.Handshake()
	//if err != nil {
	//	s.Log.Error("Handshake to socket err: ", err)
	//	return err
	//}
	//s.Log.Info("OK handshaking with socket successful!")

	s.Log.Info("Subscribe btcusdt@kline_15m to binance...")
	err = s.Subscribe()
	if err != nil {
		s.Log.Error("Subscribe btcusdt@kline_15m err: ", err)
		return err
	}
	s.Log.Info("OK subscribing to btcusdt@kline_15m!")
	s.IsConnected = true
	return nil
}

func (s *SocketBinance) Dial() error {
	if !s.Tls {
		s.Dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	connection, _, err := s.Dialer.Dial(s.WebAddress, nil)
	if err != nil {
		return err
	}
	s.Connection = connection
	return nil
}

func (s *SocketBinance) Subscribe() error {
	params := []string{"btcusdt@kline_15m"}
	request := Request{
		Method: "SUBSCRIBE",
		Params: params,
		ID:     1,
	}

	err := s.Connection.WriteJSON(request)
	if err != nil {
		return err
	}

	response := Response{}
	err = s.Connection.ReadJSON(&response)
	if err != nil {
		return err
	}

	if response.Connect != "success" && response.Error == "0" {
		return errors.New(response.Connect)
	}

	return nil
}

func (s *SocketBinance) Close() error {
	return s.Connection.Close()
}
