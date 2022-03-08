package main

import (
	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"gobot/config"
	"gobot/exchanges"
	"gobot/workers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//
func main() {
	log := logrus.New()

	configuration := config.NewConfig()

	socketBinance := exchanges.NewSocketBinance(log, configuration)
	err := socketBinance.Connect()
	if err != nil {
		log.Fatalf("Cannot connect to web-socket: %s", err)
	}
	defer socketBinance.Close()

	socketBinance.RestGet2hoursPairData()

	processor := workers.NewProcessor(socketBinance, log)

	go processor.Reading()

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
