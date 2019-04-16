package kafkaeventbus

type KafkaConfig struct {
	Server         string
	ConsumerConfig ConsumerConfig
	ProducerConfig ProducerConfig
}

type ConsumerConfig struct {
	Group        string
	Offset       string
	CommitPeriod int
	Version      string
}

type ProducerConfig struct {
	Version string
}
