package auth

import (
	"crypto/tls"
	"f1sim/config"

	"github.com/IBM/sarama"
)

// ConfigureKafkaAuth sets up authentication for Kafka producer/consumer
func ConfigureKafkaAuth(cfg *sarama.Config, kafkaConfig config.KafkaConfig) {
	// Only configure auth if credentials are provided
	if kafkaConfig.Sasl_username == "" || kafkaConfig.Security_protocol == "" {
		return
	}

	// Enable SASL authentication
	cfg.Net.SASL.Enable = true
	cfg.Net.SASL.User = kafkaConfig.Sasl_username
	cfg.Net.SASL.Password = kafkaConfig.Sasl_password
	cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	// Configure security protocol
	switch kafkaConfig.Security_protocol {
	case "SASL_SSL":
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	case "SASL_PLAINTEXT":
		cfg.Net.TLS.Enable = false
		cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext

	case "SSL":
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
}
