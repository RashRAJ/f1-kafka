package kafka

import (
	"encoding/json"
	"f1sim/internal/api"
	"fmt"
	"log"
	"strconv"

	"github.com/IBM/sarama"
)

func PublishAllData(producer sarama.SyncProducer, years []int, driver *int) {
	// 1. Publish sessions and meetings for the years
	log.Println("=== Publishing Sessions ===")
	for _, year := range years {
		PublishSessions(producer, "f1-sessions", year)
	}

	log.Println("\n=== Publishing Meetings ===")
	for _, year := range years {
		PublishMeetings(producer, "f1-meetings", year)
	}

	// 2. Fetch all session keys
	sessionKeys := fetchAllSessionKeys(years)
	log.Printf("\n=== Found %d sessions across years %v ===\n", len(sessionKeys), years)

	// 3. For each session, publish all data types to their respective topics
	topics := []string{
		"cardata", "drivers", "intervals", "laps", "locations",
		"pits", "positions", "racecontrol", "results", "grid",
		"stints", "radio", "weather",
	}

	for i, sessionKey := range sessionKeys {
		log.Printf("\n[%d/%d] Processing session %d", i+1, len(sessionKeys), sessionKey)

		for _, topic := range topics {
			kafkaTopic := "f1-" + topic
			publishForSession(producer, topic, kafkaTopic, sessionKey, driver)
		}
	}

	log.Println("\n=== All data published successfully ===")
}

func fetchAllSessionKeys(years []int) []int {
	var sessionKeys []int
	for _, year := range years {
		sessions, err := api.FetchSessions(year)
		if err != nil {
			log.Printf("Failed to fetch sessions for year %d: %v", year, err)
			continue
		}
		for _, session := range sessions {
			sessionKeys = append(sessionKeys, session.SessionKey)
		}
	}
	return sessionKeys
}

func publishForSession(producer sarama.SyncProducer, mode, topic string, sessionKey int, driver *int) {
	switch mode {
	case "cardata":
		publishCarData(producer, topic, sessionKey, driver)
	case "drivers":
		publishDrivers(producer, topic, sessionKey, driver)
	case "intervals":
		publishIntervals(producer, topic, sessionKey, driver)
	case "laps":
		publishLaps(producer, topic, sessionKey, driver)
	case "locations":
		publishLocations(producer, topic, sessionKey, driver)
	case "pits":
		publishPits(producer, topic, sessionKey)
	case "positions":
		publishPositions(producer, topic, sessionKey, driver)
	case "racecontrol":
		publishRaceControl(producer, topic, sessionKey)
	case "results":
		publishSessionResults(producer, topic, sessionKey)
	case "grid":
		publishStartingGrid(producer, topic, sessionKey)
	case "stints":
		publishStints(producer, topic, sessionKey, driver)
	case "radio":
		publishTeamRadio(producer, topic, sessionKey, driver)
	case "weather":
		publishWeather(producer, topic, sessionKey)
	}
}

func PublishSessions(producer sarama.SyncProducer, topic string, year int) {
	sessions, err := api.FetchSessions(year)
	if err != nil {
		log.Fatalf("Failed to fetch F1 sessions: %v", err)
	}

	log.Printf("Fetched %d sessions for %d", len(sessions), year)

	for i, session := range sessions {
		jsonData, err := json.Marshal(session)
		if err != nil {
			log.Printf("Error marshaling session %d: %v", i, err)
			continue
		}

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("%d", session.SessionKey)),
			Value: sarama.ByteEncoder(jsonData),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send session %d to Kafka: %v", i, err)
		} else {
			log.Printf("[%d/%d] Published: %s - %s (Session: %d) [partition=%d, offset=%d]",
				i+1, len(sessions), session.Location, session.SessionName, session.SessionKey, partition, offset)
		}
	}

	log.Println("All sessions published successfully")
}

func publishCarData(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	if driver != nil {
		log.Printf("Fetching car data for session %d, driver %d", sessionKey, *driver)
	} else {
		log.Printf("Fetching car data for session %d (all drivers)", sessionKey)
	}

	carData, err := api.FetchCarData(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch car data: %v", err)
		return
	}

	log.Printf("Fetched %d car data records", len(carData))
	publishToKafka(producer, topic, carData, "car data")
}

func publishDrivers(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	data, err := api.FetchDriverData(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch drivers: %v", err)
		return
	}
	log.Printf("Fetched %d drivers", len(data))
	publishToKafka(producer, topic, data, "drivers")
}

func publishIntervals(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	data, err := api.FetchIntervals(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch intervals: %v", err)
		return
	}
	log.Printf("Fetched %d intervals", len(data))
	publishToKafka(producer, topic, data, "intervals")
}

func publishLaps(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	data, err := api.FetchLaps(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch laps: %v", err)
		return
	}
	log.Printf("Fetched %d laps", len(data))
	publishToKafka(producer, topic, data, "laps")
}

func publishLocations(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	data, err := api.FetchLocation(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch locations: %v", err)
		return
	}
	log.Printf("Fetched %d location records", len(data))
	publishToKafka(producer, topic, data, "locations")
}

func PublishMeetings(producer sarama.SyncProducer, topic string, year int) {
	data, err := api.FetchMeetings(year, nil)
	if err != nil {
		log.Printf("Failed to fetch meetings: %v", err)
		return
	}
	log.Printf("Fetched %d meetings for %d", len(data), year)
	publishToKafka(producer, topic, data, "meetings")
}

func publishPits(producer sarama.SyncProducer, topic string, sessionKey int) {
	data, err := api.FetchPits(sessionKey)
	if err != nil {
		log.Printf("Failed to fetch pit stops: %v", err)
		return
	}
	log.Printf("Fetched %d pit stops", len(data))
	publishToKafka(producer, topic, data, "pit stops")
}

func publishPositions(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	data, err := api.FetchPositions(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch positions: %v", err)
		return
	}
	log.Printf("Fetched %d position records", len(data))
	publishToKafka(producer, topic, data, "positions")
}

func publishRaceControl(producer sarama.SyncProducer, topic string, sessionKey int) {
	data, err := api.FetchRaceControl(sessionKey)
	if err != nil {
		log.Printf("Failed to fetch race control messages: %v", err)
		return
	}
	log.Printf("Fetched %d race control messages", len(data))
	publishToKafka(producer, topic, data, "race control")
}

func publishSessionResults(producer sarama.SyncProducer, topic string, sessionKey int) {
	data, err := api.FetchSessionResults(sessionKey)
	if err != nil {
		log.Printf("Failed to fetch session results: %v", err)
		return
	}
	log.Printf("Fetched %d session results", len(data))
	publishToKafka(producer, topic, data, "session results")
}

func publishStartingGrid(producer sarama.SyncProducer, topic string, sessionKey int) {
	data, err := api.FetchStartingGrid(sessionKey)
	if err != nil {
		log.Printf("Failed to fetch starting grid: %v", err)
		return
	}
	log.Printf("Fetched %d starting grid positions", len(data))
	publishToKafka(producer, topic, data, "starting grid")
}

func publishStints(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	data, err := api.FetchStints(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch stints: %v", err)
		return
	}
	log.Printf("Fetched %d stints", len(data))
	publishToKafka(producer, topic, data, "stints")
}

func publishTeamRadio(producer sarama.SyncProducer, topic string, sessionKey int, driver *int) {
	data, err := api.FetchTeamRadio(sessionKey, driver)
	if err != nil {
		log.Printf("Failed to fetch team radio: %v", err)
		return
	}
	log.Printf("Fetched %d team radio messages", len(data))
	publishToKafka(producer, topic, data, "team radio")
}

func publishWeather(producer sarama.SyncProducer, topic string, sessionKey int) {
	data, err := api.FetchWeather(sessionKey)
	if err != nil {
		log.Printf("Failed to fetch weather: %v", err)
		return
	}
	log.Printf("Fetched %d weather records", len(data))
	publishToKafka(producer, topic, data, "weather")
}

// publishToKafka is a generic function to publish any data type to Kafka
func publishToKafka(producer sarama.SyncProducer, topic string, data interface{}, dataType string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling %s: %v", dataType, err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(strconv.FormatInt(int64(len(jsonData)), 10)),
		Value: sarama.ByteEncoder(jsonData),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send %s to Kafka: %v", dataType, err)
	}

	log.Printf("Published %s to Kafka [partition=%d, offset=%d]", dataType, partition, offset)
}
