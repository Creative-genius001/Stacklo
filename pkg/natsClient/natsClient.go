package natsclient

import (
	"os"
	"time"

	"github.com/Creative-genius001/go-logger"
	"github.com/nats-io/nats.go"
)

var Conn *nats.Conn

func InitNATS() {

	url := os.Getenv("NATS_URL")
	if url == "" {
		logger.Fatal("could not read NATS URL")
	}
	nc, err := nats.Connect(url, nats.Timeout(5*time.Second), nats.PingInterval(20*time.Second), nats.MaxPingsOutstanding(5), nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		logger.Fatal("Connection error detected: ", err)
	}))
	if err != nil {
		logger.Fatal("Failed to connect to NATS: ", err)
	}
	Conn = nc
	logger.Info("Connected to NATS")
}

func Close() {
	if Conn != nil {
		Conn.Close()
		logger.Info("Connection closed")
	}
}
