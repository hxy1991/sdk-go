package delayQueue

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestDelayQueue(t *testing.T) {
	d, err := New("my_delay_queue")
	if err != nil {
		t.Fatal(err)
	}

	delayInSecondX := 3
	delayInMilliX := delayInSecondX * 1000
	err = d.Publish("text/plain", "Hello world! I am x.", delayInMilliX)
	if err != nil {
		t.Fatal(err)
	}

	delayInSecondY := 6
	delayInMilliY := delayInSecondY * 1000
	err = d.Publish("text/plain", "Hello world! I am y.", delayInMilliY)
	if err != nil {
		t.Fatal(err)
	}

	forever := make(chan bool)
	err = d.Consume(func(delivery amqp.Delivery) {
		t.Logf("consumer receive, body: %s, timestamp: %v", string(delivery.Body), delivery.Timestamp)

		err := delivery.Ack(false)
		if err != nil {
			t.Error(err)
			return
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	<-forever
}
