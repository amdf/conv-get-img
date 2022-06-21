package producer

import (
	"log"

	"github.com/Shopify/sarama"
)

var brokers = []string{
	":9095",
	":9096",
	":9097",
}

func newConfig() (config *sarama.Config) {
	config = sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	return
}

func NewSync() (p sarama.SyncProducer, err error) {
	config := newConfig()
	p, err = sarama.NewSyncProducer(brokers, config)
	if nil == err {
		log.Println("sync connected to kafka brokers")
	}

	return
}

func NewAsync() (p sarama.AsyncProducer, err error) {
	config := newConfig()
	p, err = sarama.NewAsyncProducer(brokers, config)
	if nil == err {
		log.Println("async connected to kafka brokers")
	}
	return
}

func PrepareMessage(topic string, message []byte) (msg *sarama.ProducerMessage) {
	msg = &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(message),
	}
	return
}
