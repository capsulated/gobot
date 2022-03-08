package config

import "os"

type Logger struct {
	Path  string
	Level string
	Json  bool
}

type BinanceSocket struct {
	ApiKey    string
	SecretKey string
	Address   string
	Port      string
	Path      string
}

type Config struct {
	Logger        *Logger
	BinanceSocket *BinanceSocket
}

// NewConfig todo add all from env
func NewConfig() *Config {
	logger := &Logger{
		Path:  "./",
		Level: "debug",
		Json:  true,
	}

	binanceSocket := &BinanceSocket{
		ApiKey:    os.Getenv("API_KEY"),
		SecretKey: os.Getenv("SECRET_KEY"),
		Address:   "stream.binance.com",
		Port:      "9443",
		Path:      "/ws",
	}

	return &Config{
		Logger:        logger,
		BinanceSocket: binanceSocket,
	}
}
