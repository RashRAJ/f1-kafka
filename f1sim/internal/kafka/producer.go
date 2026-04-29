package kafka

import (
	"log"
	"time"

	"f1sim/auth"
	"f1sim/config"

	"github.com/IBM/sarama"
)

func NewProducer(kafkaConfig *config.KafkaConfig) sarama.SyncProducer {
	// Setup Kafka producer configuration
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Successes = true
	producerConfig.Producer.RequiredAcks = sarama.RequiredAcks(kafkaConfig.RequiredAcks)
	producerConfig.Producer.Retry.Max = kafkaConfig.RetryMax
	producerConfig.Producer.Retry.Backoff = time.Duration(kafkaConfig.RetryBackoff) * time.Millisecond
	producerConfig.Producer.Flush.Frequency = time.Duration(kafkaConfig.LingerMs) * time.Millisecond
	producerConfig.Producer.Flush.MaxMessages = kafkaConfig.BatchSize
	producerConfig.ChannelBufferSize = kafkaConfig.BufferMemory

	// Configure authentication for Confluent Cloud
	auth.ConfigureKafkaAuth(producerConfig, *kafkaConfig)

	// Enable verbose logging for debugging
	producerConfig.Producer.Return.Errors = true
	producerConfig.Version = sarama.V2_6_0_0 // Confluent Cloud compatibility

	producerInstance, err := sarama.NewSyncProducer(kafkaConfig.Brokers, producerConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	log.Printf("Connected to Kafka brokers: %v", kafkaConfig.Brokers)
	return producerInstance
}
