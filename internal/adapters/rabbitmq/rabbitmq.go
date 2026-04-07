package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/nlsnnn/berezhok/internal/lib/logger/sl"
)

const (
	ExchangeNotifications = "notifications"
)

type Client struct {
	url    string
	conn   *amqp.Connection
	ch     *amqp.Channel
	mu     sync.Mutex
	logger *slog.Logger
}

func New(url string, logger *slog.Logger) (*Client, error) {
	c := &Client{url: url, logger: logger}
	if err := c.connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("rabbitmq dial: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		err = conn.Close()
		c.failOnError(err, "rabbitmq open channel")
		return fmt.Errorf("rabbitmq open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		ExchangeNotifications,
		amqp.ExchangeTopic,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		err = ch.Close()
		c.failOnError(err, "rabbitmq declare exchange")
		err = conn.Close()
		c.failOnError(err, "rabbitmq declare exchange")
		return fmt.Errorf("rabbitmq declare exchange: %w", err)
	}

	// Publisher confirms - брокер подтверждает, что принял сообщение
	if err = ch.Confirm(false); err != nil {
		err = ch.Close()
		c.failOnError(err, "rabbitmq confirm mode")
		err = conn.Close()
		c.failOnError(err, "rabbitmq confirm mode")
		return fmt.Errorf("rabbitmq confirm mode: %w", err)
	}

	c.conn = conn
	c.ch = ch
	return nil
}

// reconnect пытается переподключиться с экспоненциальным backoff
func (c *Client) reconnect() {
	backoff := time.Second
	for {
		c.logger.Warn("rabbitmq reconnecting...", "backoff", backoff)
		time.Sleep(backoff)

		if err := c.connect(); err == nil {
			c.logger.Info("rabbitmq reconnected")
			return
		}

		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func (c *Client) Publish(ctx context.Context, routingKey string, body []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.ch.PublishWithContext(ctx,
		ExchangeNotifications,
		routingKey,
		true, // mandatory — вернёт ошибку, если нет подходящей очереди
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // сообщение переживёт рестарт брокера
			Body:         body,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		// Канал после ошибки публикации становится невалидным — переподключаемся
		go c.reconnect()
		return fmt.Errorf("rabbitmq publish: %w", err)
	}

	// Ждём подтверждения от брокера
	confirmed, ok := <-c.ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	if !ok || !confirmed.Ack {
		return fmt.Errorf("rabbitmq: broker did not confirm message delivery")
	}

	return nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.ch.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

func (c *Client) failOnError(err error, msg string) {
	if err != nil {
		c.logger.Error(msg, sl.Err(err))
	}
}
