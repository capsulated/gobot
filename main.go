package main

import (
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gobot/config"
	"gobot/exchanges"
	"gobot/server"
	"gobot/telegram"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Logger
	log := logrus.New()

	// Env
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	// Config
	configuration := config.NewConfig()

	//Socket
	socketBinance := exchanges.NewSocketBinance(log, configuration)
	err = socketBinance.Connect()
	if err != nil {
		log.Fatalf("Cannot connect to web-socket: %s", err)
	}
	defer socketBinance.Close()
	//socketBinance.RestGet2hoursPairData()

	// Telega
	telega, err := telegabot.NewTelegaBot(log, configuration)
	if err != nil {
		log.Fatalf("Cannot create telega bot: %s", err)
	}

	srv := server.NewServer(log, telega, configuration)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Cannot create server: %s", err)
	}

	// NewServer
	// processor := workers.NewProcessor(socketBinance, log)

	// go processor.Reading()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT)

	for {
		select {
		case sig := <-signals:
			switch sig {
			case os.Interrupt:
				log.Warn("os.Interrupt")
				stopService(log, socketBinance)
			case os.Kill:
				log.Warn("os.Kill")
				stopService(log, socketBinance)
			case syscall.SIGINT:
				log.Warn("syscall.SIGINT")
				stopService(log, socketBinance)
			case syscall.SIGQUIT:
				log.Warn("syscall.SIGQUIT")
				stopService(log, socketBinance)
			case syscall.SIGTERM:
				log.Warn("syscall.SIGTERM")
				stopService(log, socketBinance)
			case syscall.SIGABRT:
				log.Warn("syscall.SIGABRT")
				stopService(log, socketBinance)
			}
		}
	}
}

func stopService(log *logrus.Logger, binanceSocket *exchanges.SocketBinance) {
	binanceSocket.Terminate = true
	log.Warn("Stopping service")
	// To cleanly close a connection, a client should send a close
	// frame and wait for the server to close the connection.
	err := binanceSocket.Connection.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(2 * time.Second)

	if err != nil {
		log.Warnf("Write close:", err)
	}
	binanceSocket.Close()

	log.Info("Connection to web-socket closed")
	time.Sleep(5 * time.Second)

	log.Info("Log file closed. Service finished work successful!")
	os.Exit(0)
}
