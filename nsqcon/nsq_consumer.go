package nsqcon

import (
	"time"

	"github.com/cocobao/log"
	nsq "github.com/nsqio/go-nsq"
)

var (
	nsqConsumer *nsq.Consumer
	handleChan  chan []byte
)

type ConsumerT struct{}

func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	log.Debug("receive:", msg.NSQDAddress, ", message:", string(msg.Body))
	handleChan <- msg.Body
	return nil
}

func SetupConsumer(topic string, channel string, address string, msgchan chan []byte) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second
	c, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		panic(err)
	}

	c.SetLogger(nil, 0)
	c.AddHandler(&ConsumerT{})
	// if err := c.ConnectToNSQLookupd(address); err != nil {
	// 	panic(err)
	// }
	if err := c.ConnectToNSQD(address); err != nil {
		panic(err)
	}

	nsqConsumer = c
	handleChan = msgchan
}
