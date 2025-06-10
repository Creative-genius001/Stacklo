package natsclient

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var Conn *nats.Conn

func InitNATS() {
	url := os.Getenv("NATS_URL")
	nc, err := nats.Connect(url, nats.Timeout(5*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	Conn = nc
	log.Println("Connected to NATS")
}

func Close() {
	if Conn != nil {
		Conn.Close()
	}
}
