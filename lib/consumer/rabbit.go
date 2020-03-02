package consumer

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	logger "github.com/labstack/gommon/log"
	"github.com/manucorporat/try"
	"github.com/streadway/amqp"
	"github.com/tientp-floware/mgodb-stream/config"
)

var (
	log = logger.GetLogger("[partner]")
)

// Consumer hold rabbit connect
type Consumer struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	done            chan error
	consumerTag     string // Name that consumer identifies itself to the server with
	uri             string // uri of the rabbitmq server
	exchange        string // exchange that we will bind to
	exchangeType    string // topic, direct, etc...
	lastRecoverTime int64
	//track service current status
	currentStatus atomic.Value
}

const RECOVER_INTERVAL_TIME = 6 * 60

// NewConsumer returns a Consumer struct that has been initialized properly
// essentially don't touch conn, channel, or done and you can create Consumer manually
func newConsumer(consumerTag, uri, exchange, exchangeType string) *Consumer {
	name, err := os.Hostname()
	if err != nil {
		name = "ece"
	}
	consumer := &Consumer{
		consumerTag:     fmt.Sprintf("%s_%s", name, consumerTag),
		uri:             uri,
		exchange:        exchange,
		exchangeType:    exchangeType,
		done:            make(chan error),
		lastRecoverTime: time.Now().Unix(),
	}
	consumer.currentStatus.Store(true)
	return consumer
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU

}

func RunConsumer(rabbitUri, consumerTag, exchange, exchangeType, queueName, routingKey string, handler func([]byte) bool) {
	consumer := newConsumer(consumerTag, rabbitUri, exchange, exchangeType)
	if err := consumer.Connect(); err != nil {
		fmt.Printf("[%s]connect error", consumerTag)
	}
	deliveries, _ := consumer.AnnounceQueue(queueName, routingKey)
	consumer.Handle(deliveries, handler, maxParallelism(), queueName, routingKey)
}

func (c *Consumer) ReConnect(queueName, routingKey string, retryTime int) (<-chan amqp.Delivery, error) {
	c.Close()
	time.Sleep(time.Duration(config.Config.RabbitMq.Second) * time.Second)
	log.Info("Try ReConnect with times:", retryTime)

	if err := c.Connect(); err != nil {
		return nil, err
	}
	deliveries, err := c.AnnounceQueue(queueName, routingKey)
	if err != nil {
		return deliveries, errors.New("Couldn't connect")
	}
	return deliveries, nil
}

// Connect to RabbitMQ server
func (c *Consumer) Connect() error {

	var err error
	log.Info("dialing: ", c.uri)
	c.conn, err = amqp.Dial(c.uri)

	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		log.Info("closing: ", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Let Handle know it's not time to reconnect
		c.done <- errors.New("Channel Closed")
	}()

	log.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Info("got Channel, declaring Exchange ", c.exchange)
	if err = c.channel.ExchangeDeclare(
		c.exchange,     // name of the exchange
		c.exchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	return nil
}

// AnnounceQueue sets the queue that will be listened to for this
// connection...
func (c *Consumer) AnnounceQueue(queueName, routingKey string) (<-chan amqp.Delivery, error) {
	log.Info("declared Exchange, declaring Queue:", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Info(fmt.Sprintf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, routingKey))
	err = c.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		routingKey, // routingKey
		c.exchange, // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Info("Queue bound to Exchange, starting Consume consumer tag:", c.consumerTag)
	deliveries, err := c.channel.Consume(
		queue.Name,    // name
		c.consumerTag, // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}
	return deliveries, nil
}

// Close connnect rabbitmq
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
		c.channel = nil
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

// Handle consumer handle
func (c *Consumer) Handle(
	deliveries <-chan amqp.Delivery,
	fn func([]byte) bool,
	threads int,
	queue string,
	routingKey string) {
	var err error
	for {
		log.Info("Enter for busy loop with thread:", threads)
		for i := 0; i < threads; i++ {
			go func() {
				log.Info("Enter go with thread with deliveries", deliveries)
				for msg := range deliveries {
					log.Info("Enter deliver")
					ret := false
					try.This(func() {
						body := msg.Body[:]
						ret = fn(body)
					}).Finally(func() {
						if ret == true {
							msg.Ack(false)
						} else {
							msg.Reject(false)
						}
					}).Catch(func(e try.E) {
						log.Error(e)
					})
				}
			}()
		}
		// Go into reconnect loop when
		// c.done is passed non nil values
		if <-c.done != nil {
			retryTime := 1
			for {
				deliveries, err = c.ReConnect(queue, routingKey, retryTime)
				if err != nil {
					log.Error("Reconnecting Error")
					retryTime++
				} else {
					break
				}
			}
		}
		log.Info("Reconnected!!!")
	}
}
