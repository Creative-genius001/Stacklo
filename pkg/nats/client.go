package nats

import (
	"fmt"
	"time"

	"github.com/Creative-genius001/go-logger"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NATS interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error)
	QueueSubscribe(subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error)
	JetStreamPublish(subject string, data []byte) (*nats.PubAck, error)
	JetStreamQueueSubscribe(stream, subject, queue, durable string, handler nats.MsgHandler) (*nats.Subscription, error)
	Close()
}

type natsClient struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func NewNATSClient(natsURL string) (NATS, error) {
	conn, err := nats.Connect(natsURL,
		nats.Timeout(10*time.Second),
		nats.PingInterval(20*time.Second),
		nats.MaxPingsOutstanding(3),
		nats.ReconnectWait(1*time.Second),
		nats.MaxReconnects(10),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error("NATS disconnected", zap.Error(err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Error("NATS connection closed")
		}),
	)
	if err != nil {
		logger.Error("Failed to connect to NATS", zap.String("nats_url", natsURL), zap.Error(err))
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		conn.Close()
		logger.Error("Failed to get JetStream context", zap.Error(err))
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	logger.Info("Successfully connected to NATS and JetStream", zap.String("nats_url", natsURL))
	return &natsClient{conn: conn, js: js}, nil
}

func (n *natsClient) Publish(subject string, data []byte) error {
	return n.conn.Publish(subject, data)
}

func (n *natsClient) Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return n.conn.Subscribe(subject, handler)
}

func (n *natsClient) QueueSubscribe(subject, queue string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return n.conn.QueueSubscribe(subject, queue, handler)
}

func (n *natsClient) JetStreamPublish(subject string, data []byte) (*nats.PubAck, error) {
	return n.js.Publish(subject, data)
}

func (n *natsClient) JetStreamQueueSubscribe(stream, subject, queue, durable string, handler nats.MsgHandler) (*nats.Subscription, error) {
	// Ensure the stream exists (optional, but good for setup)
	// You might create streams during deployment or startup of your services
	_, err := n.js.AddStream(&nats.StreamConfig{
		Name:      stream,
		Subjects:  []string{subject},
		Retention: nats.LimitsPolicy, // Or nats.InterestPolicy, nats.WorkQueuePolicy
		MaxBytes:  1024 * 1024 * 10,  // 10MB max stream size (adjust as needed)
		MaxMsgs:   100000,            // Max 100k messages (adjust as needed)
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		logger.Error("Failed to add JetStream stream", zap.String("stream", stream), zap.Error(err))
		return nil, fmt.Errorf("failed to add JetStream stream: %w", err)
	}

	// Subscribe with Acknowledge policy (explicit ack)
	return n.js.QueueSubscribe(subject, queue, handler,
		nats.Durable(durable),
		nats.ManualAck(),             // Consumer must explicitly acknowledge messages
		nats.DeliverNew(),            // Deliver new messages only
		nats.AckWait(30*time.Second), // How long to wait for an ack before redelivering
		nats.MaxAckPending(100),      // Max unacked messages at a time
	)
}

func (n *natsClient) Close() {
	if n.conn != nil {
		n.conn.Close()
		logger.Info("NATS connection closed")
	}
}
