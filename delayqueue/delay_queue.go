package delayQueue

import (
	"fmt"
	"github.com/hxy1991/sdk-go/log"
	rabbitMQ "github.com/hxy1991/sdk-go/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"strconv"
	"sync"
	"time"
)

var initConnection sync.Once

var producerConn *rabbitMQ.Connection
var consumerConn *rabbitMQ.Connection

type DelayQueue struct {
	name string
}

func New(name string) (*DelayQueue, error) {
	initConnection.Do(func() {
		_producerConn, err := rabbitMQ.New()
		if err != nil {
			log.Error(err)
			return
		}
		producerConn = _producerConn

		log.With().Info("init shared producer connection of all delayQueues successfully")

		_consumerConn, err := rabbitMQ.New()
		if err != nil {
			log.With().Error(err)
			return
		}
		consumerConn = _consumerConn

		log.With().Info("init shared consumer connection of all delayQueues successfully")
	})

	if producerConn == nil {
		return nil, fmt.Errorf("init producer connection fail")
	}

	if consumerConn == nil {
		return nil, fmt.Errorf("init consumer connection fail")
	}

	d := &DelayQueue{
		name: name,
	}

	err := d.exchangeDeclare()
	if err != nil {
		return nil, err
	}

	err = d.producerQueueDeclare()
	if err != nil {
		return nil, err
	}

	err = d.consumerQueueDeclare()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d DelayQueue) exchangeDeclare() error {
	ch, err := producerConn.Channel()
	if err != nil {
		return err
	}
	defer d.closeChannel(ch)

	err = rabbitMQ.ExchangeDeclareDefault(ch, d.getExchange())
	if err != nil {
		return err
	}
	return nil
}

func (d DelayQueue) producerQueueDeclare() error {
	ch, err := producerConn.Channel()
	if err != nil {
		return err
	}
	defer d.closeChannel(ch)

	_, err = rabbitMQ.QueueDeclareDefault(ch, d.getQueueNameForProducer(), d.getExchange())
	if err != nil {
		return err
	}
	return nil
}

func (d DelayQueue) consumerQueueDeclare() error {
	ch, err := consumerConn.Channel()
	if err != nil {
		return err
	}
	defer d.closeChannel(ch)

	q, err := rabbitMQ.QueueDeclareDefault(ch, d.getQueueNameForConsumer(), "")
	if err != nil {
		return err
	}

	err = rabbitMQ.QueueBindDefault(ch, q.Name, d.getExchange())
	if err != nil {
		return err
	}
	return nil
}

func (d DelayQueue) Publish(contentType, body string, delayInMilli int) error {
	ch, err := producerConn.Channel()
	if err != nil {
		return err
	}
	defer d.closeChannel(ch)

	// RabbitMQ默认提供了一个Exchange，名字是空字符串，类型是Direct，绑定到所有的Queue。 每一个Queue和这个无名Exchange之间的Binding Key是Queue的名字。
	return ch.Publish(
		"",
		d.getQueueNameForProducer(),
		false,
		false,
		amqp.Publishing{
			// 持久化消息
			DeliveryMode: 2,
			ContentType:  contentType,
			Body:         []byte(body),
			Expiration:   strconv.Itoa(delayInMilli),
			Timestamp:    time.Now(),
		})
}

func (d DelayQueue) Consume(callback func(amqp.Delivery)) error {
	ch, err := consumerConn.Channel()
	if err != nil {
		return err
	}

	name := d.getQueueNameForConsumer()
	delivery, err := rabbitMQ.ConsumeDefault(ch, name)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.With("queueName", name).Error(e)
			}
		}()

		for d := range delivery {
			callback(d)
		}
	}()

	log.With("queueName", name).Info("consumer is waiting")

	return nil
}

func (d DelayQueue) closeChannel(ch *amqp.Channel) {
	func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			log.With().Error(err)
		}
	}(ch)
}

func (d DelayQueue) getQueueNameForProducer() string {
	return d.name + "_producer"
}

func (d DelayQueue) getQueueNameForConsumer() string {
	return d.name + "_consumer"
}

func (d DelayQueue) getExchange() string {
	delayExChange := d.name + ".delay"
	return delayExChange
}
