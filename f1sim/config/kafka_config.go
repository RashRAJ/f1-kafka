package config

type TopicConfig struct {
	Name              string            `yaml:"name"`
	Partitions        int32             `yaml:"partitions"`
	ReplicationFactor int16             `yaml:"replication_factor"`
	Config            map[string]string `yaml:"config"`
}

type KafkaConfig struct {
	Brokers           []string       `yaml:"brokers"`
	Topic             string         `yaml:"topic"`
	RetryMax          int            `yaml:"retry_max"`
	RetryBackoff      int            `yaml:"retry_backoff"`
	RequiredAcks      int            `yaml:"required_acks"`
	LingerMs          int            `yaml:"linger_ms"`
	BatchSize         int            `yaml:"batch_size"`
	BufferMemory      int            `yaml:"buffer_memory"`
	Sasl_username     string         `yaml:"sasl_username"`
	Sasl_password     string         `yaml:"sasl_password"`
	Security_protocol string         `yaml:"security_protocol"`
	Topics            []TopicConfig  `yaml:"topics"`
}
