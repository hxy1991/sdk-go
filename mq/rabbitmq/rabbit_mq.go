package rabbitMQ

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn *amqp.Connection
}

func New() (*Connection, error) {
	// 建立链接
	_conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	return &Connection{
		conn: _conn,
	}, nil
}

func (c *Connection) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Connection) Channel() (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func ExchangeDeclareDefault(ch *amqp.Channel, name string) error {
	err := ch.ExchangeDeclare(
		name,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func QueueDeclareDefault(ch *amqp.Channel, name, deadLetterExChange string) (amqp.Queue, error) {
	var arg amqp.Table
	if deadLetterExChange != "" {
		arg = amqp.Table{
			// 当消息过期时把消息发送到 logs 这个 exchange
			"x-dead-letter-exchange": deadLetterExChange,
		}
	}
	q, err := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		arg,
	)
	return q, err
}

func QueueBindDefault(ch *amqp.Channel, name, exchange string) error {
	return ch.QueueBind(
		name,
		"",
		exchange,
		false,
		nil)
}

func ConsumeDefault(ch *amqp.Channel, name string) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}
