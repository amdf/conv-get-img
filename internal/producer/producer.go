package producer

import "github.com/Shopify/sarama"

var brokers = []string{
	"127.0.0.1:9095",
	"127.0.0.1:9096",
	"127.0.0.1:9097",
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
	return
}

func NewAsync() (p sarama.AsyncProducer, err error) {
	config := newConfig()
	p, err = sarama.NewAsyncProducer(brokers, config)
	return
}
