package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadKafkaConfig() KafkaConfig {
	yamlData, err := os.ReadFile("kafka_config.yaml")
	if err != nil {
		log.Fatal("Error while reading Kafka config file", err)
	}

	var kafkaConfig KafkaConfig
	if err := yaml.Unmarshal(yamlData, &kafkaConfig); err != nil {
		log.Fatal("Error while parsing Kafka YAML config", err)
	}

	if v := os.Getenv("KAFKA_SASL_USERNAME"); v != "" {
		kafkaConfig.Sasl_username = v
	}
	if v := os.Getenv("KAFKA_SASL_PASSWORD"); v != "" {
		kafkaConfig.Sasl_password = v
	}
	if kafkaConfig.Sasl_username == "" || kafkaConfig.Sasl_password == "" {
		log.Fatal("KAFKA_SASL_USERNAME and KAFKA_SASL_PASSWORD must be set")
	}

	return kafkaConfig
}

func LoadAppConfig() AppConfig {
	yamlData, err := os.ReadFile("app_config.yaml")
	if err != nil {
		log.Fatal("Error while reading App config file", err)
	}

	var appConfig AppConfig
	if err := yaml.Unmarshal(yamlData, &appConfig); err != nil {
		log.Fatal("Error while parsing App YAML config", err)
	}

	return appConfig
}
