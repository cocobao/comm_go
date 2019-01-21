package nsqcon

import (
	"encoding/json"

	nsq "github.com/nsqio/go-nsq"
)

var (
	nsqprod *NsqProducerInfo
)

type NsqProducerInfo struct {
	producer *nsq.Producer
	addrs    []string
	index    int
}

func SetupProducer(addrs []string) {
	nsqprod = &NsqProducerInfo{
		addrs: addrs,
		index: 0,
	}
}

func PublishMsg(topic string, val interface{}) error {
	if nsqprod.producer == nil {
		nsqprod.index++
		if nsqprod.index >= len(nsqprod.addrs) {
			nsqprod.index = 0
		}

		config := nsq.NewConfig()
		w, err := nsq.NewProducer(nsqprod.addrs[nsqprod.index], config)
		if err != nil {
			return err
		}
		nsqprod.producer = w
	}

	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if err := nsqprod.producer.Publish(topic, data); err != nil {
		nsqprod.producer.Stop()
		nsqprod.producer = nil
		return err
	}

	return nil
}
