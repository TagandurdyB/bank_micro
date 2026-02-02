package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQClient(url string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Create the queue (If Durable: True, messages will not be deleted even if RabbitMQ closes)
	_, err = ch.QueueDeclare(
		"deposit_queue", // Queue name
		true,            // Durable
		false,           // Auto-delete
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)

	return &RabbitMQClient{Conn: conn, Channel: ch}, err
}

func (r *RabbitMQClient) Publish(queueName string, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return r.Channel.PublishWithContext(
		context.Background(),
		"",        // Exchange (Default)
		queueName, // Routing Key
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent, // It saves the message to disk.
		},
	)
}

func (r *RabbitMQClient) Consume() {

}

func (r *RabbitMQClient) Close() {
	r.Channel.Close()
	r.Conn.Close()
}
