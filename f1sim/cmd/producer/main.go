package main

import (
	"log"

	"f1sim/config"
	"f1sim/internal/kafka"
)

func main() {
	// Load Kafka configuration
	kafkaConfig := config.LoadKafkaConfig()

	// Create topics from configuration before starting producer
	log.Println("Ensuring Kafka topics are created from configuration...")
	if err := kafka.CreateTopics(&kafkaConfig); err != nil {
		log.Fatalf("Failed to create topics: %v", err)
	}

	// Create producer with the loaded config
	kafkaProducer := kafka.NewProducer(&kafkaConfig)
	defer kafkaProducer.Close()

	appConfig := config.LoadAppConfig()
	driver := kafka.ParseDriverNumber(appConfig.DriverNumber)
	years := kafka.ParseYears(appConfig.YearsStr)

	if appConfig.AllData {
		log.Printf("Fetching ALL data types for years %v", years)
		kafka.PublishAllData(kafkaProducer, years, driver)
	} else {
		log.Println("Set all_data=true in app_config.yaml to fetch all data types")
	}
}
