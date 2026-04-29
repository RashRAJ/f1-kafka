package kafka

import (
	"fmt"
	"log"
	"time"

	"f1sim/auth"
	"f1sim/config"

	"github.com/IBM/sarama"
)

// CreateTopics creates Kafka topics from configuration if they don't exist
func CreateTopics(kafkaConfig *config.KafkaConfig) error {
	if len(kafkaConfig.Topics) == 0 {
		log.Println("No topics configured for creation")
		return nil
	}

	// Setup admin client configuration
	adminConfig := sarama.NewConfig()
	adminConfig.Version = sarama.V2_6_0_0

	// Configure authentication for Confluent Cloud
	auth.ConfigureKafkaAuth(adminConfig, *kafkaConfig)

	// Create admin client
	admin, err := sarama.NewClusterAdmin(kafkaConfig.Brokers, adminConfig)
	if err != nil {
		return fmt.Errorf("failed to create admin client: %w", err)
	}
	defer admin.Close()

	// Get existing topics
	existingTopics, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	// Create each configured topic
	for _, topicConfig := range kafkaConfig.Topics {
		// Check if topic already exists
		if _, exists := existingTopics[topicConfig.Name]; exists {
			log.Printf("Topic '%s' already exists, skipping creation", topicConfig.Name)

			// Optionally validate configuration matches
			existing := existingTopics[topicConfig.Name]
			if existing.NumPartitions != topicConfig.Partitions {
				log.Printf("Warning: Topic '%s' has %d partitions, config specifies %d",
					topicConfig.Name, existing.NumPartitions, topicConfig.Partitions)
			}
			continue
		}

		// Convert config map[string]string to map[string]*string
		configEntries := make(map[string]*string)
		for key, value := range topicConfig.Config {
			v := value
			configEntries[key] = &v
		}

		// Prepare topic detail
		topicDetail := &sarama.TopicDetail{
			NumPartitions:     topicConfig.Partitions,
			ReplicationFactor: topicConfig.ReplicationFactor,
			ConfigEntries:     configEntries,
		}

		// Create topic
		log.Printf("Creating topic '%s' with %d partitions and replication factor %d",
			topicConfig.Name, topicConfig.Partitions, topicConfig.ReplicationFactor)

		err = admin.CreateTopic(topicConfig.Name, topicDetail, false)
		if err != nil {
			// Check if error is because topic already exists (race condition)
			if topicErr, ok := err.(*sarama.TopicError); ok && topicErr.Err == sarama.ErrTopicAlreadyExists {
				log.Printf("Topic '%s' was created by another process", topicConfig.Name)
				continue
			}
			return fmt.Errorf("failed to create topic '%s': %w", topicConfig.Name, err)
		}

		// Wait a bit for topic to be fully created
		time.Sleep(2 * time.Second)
		log.Printf("Successfully created topic '%s'", topicConfig.Name)
	}

	log.Println("Topic creation completed")
	return nil
}

// EnsureTopicExists verifies a topic exists, or creates it if configured
func EnsureTopicExists(topicName string, kafkaConfig *config.KafkaConfig) error {
	// Check if topic is in configuration
	for _, topicConfig := range kafkaConfig.Topics {
		if topicConfig.Name == topicName {
			return CreateTopics(kafkaConfig)
		}
	}

	log.Printf("Warning: Topic '%s' not found in configuration, relying on auto-creation", topicName)
	return nil
}
