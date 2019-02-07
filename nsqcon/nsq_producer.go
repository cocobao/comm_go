package nsqcon

import (
	"encoding/json"

	"github.com/cocobao/log"
	nsq "github.com/nsqio/go-nsq"
)

var (
	nsqprod *NsqProducerInfo
)

type NsqProducerInfo struct {
	producer []*nsq.Producer
	addrs    []string
	index    int
}

func SetupProducer(addrs []string) {
	if len(addrs) == 0 {
		return
	}
	nsqprod = &NsqProducerInfo{
		addrs:    addrs,
		index:    0,
		producer: make([]*nsq.Producer, len(addrs)),
	}
}

func PublishMsg(topic string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	index := nsqprod.index
	producer := nsqprod.producer[index]
	if producer == nil {
		config := nsq.NewConfig()
		w, err := nsq.NewProducer(nsqprod.addrs[index], config)
		if err != nil {
			return err
		}
		producer = w
		nsqprod.producer[index] = producer
		log.Debug("new nsq productor ok")
	}

	if err := producer.Publish(topic, data); err != nil {
		nsqprod.producer[index].Stop()
		nsqprod.producer[index] = nil
		// nsqprod.index++
		// if nsqprod.index >= len(nsqprod.addrs) {
		// 	nsqprod.index = 0
		// }
		return err
	}

	return nil
}
