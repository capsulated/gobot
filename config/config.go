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

type Telega struct {
	ApiToken string
	Debug    bool
}

type Server struct {
	Host string
	Port string
}

type Config struct {
	Logger        *Logger
	Telega        *Telega
	BinanceSocket *BinanceSocket
	Server        *Server
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

	telegramDebug := os.Getenv("TELEGRAM_DEBUG")
	var debug bool
	if telegramDebug == "true" {
		debug = true
	}
	telega := &Telega{
		ApiToken: os.Getenv("TELEGRAM_APITOKEN"),
		Debug:    debug,
	}

	server := &Server{
		Host: "",
		Port: "3000",
	}

	return &Config{
		Logger:        logger,
		BinanceSocket: binanceSocket,
		Telega:        telega,
		Server:        server,
	}
}
