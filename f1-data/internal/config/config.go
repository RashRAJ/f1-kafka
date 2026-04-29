package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type KafkaConfig struct {
	Brokers           []string `yaml:"brokers"`
	Topic             string   `yaml:"topic"`
	RetryMax          int      `yaml:"retry_max"`
	RetryBackoff      int      `yaml:"retry_backoff"`
	RequiredAcks      int      `yaml:"required_acks"`
	LingerMs          int      `yaml:"linger_ms"`
	BatchSize         int      `yaml:"batch_size"`
	BufferMemory      int      `yaml:"buffer_memory"`
	Sasl_username     string   `yaml:"sasl_username"`
	Sasl_password     string   `yaml:"sasl_password"`
	Security_protocol string   `yaml:"security_protocol"`
}

type AppConfig struct {
	DriverNumber string `yaml:"driver,omitempty"`
	YearsStr     string `yaml:"years"`
	AllData      bool   `yaml:"all_data"`
	BaseURL      string `yaml:"base_url"`
}

func LoadKafkaConfig() KafkaConfig {
	yamlData, err := os.ReadFile("kafka_config.yaml")
	if err != nil {
		log.Fatal("Error while reading Kafka config file", err)
	}

	var kafkaConfig KafkaConfig
	if err := yaml.Unmarshal(yamlData, &kafkaConfig); err != nil {
		log.Fatal("Error while parsing Kafka YAML config", err)
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
