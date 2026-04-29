package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"f1sim/internal/telemetry"

	"github.com/IBM/sarama"
)

func main() {
	// Kafka configuration
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9093" // Default for development
	}

	// Setup Kafka producer with retry logic
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	var producer sarama.SyncProducer
	var err error

	// Retry connecting to Kafka with exponential backoff
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		producer, err = sarama.NewSyncProducer([]string{kafkaBroker}, config)
		if err == nil {
			break
		}

		waitTime := time.Duration(i+1) * time.Second
		log.Printf("Failed to connect to Kafka (attempt %d/%d): %v. Retrying in %v...",
			i+1, maxRetries, err, waitTime)
		time.Sleep(waitTime)
	}

	if err != nil {
		log.Fatalf("Failed to create Kafka producer after %d attempts: %v", maxRetries, err)
	}
	defer producer.Close()

	log.Printf("Connected to Kafka broker: %s\n", kafkaBroker)

	// Setup UDP listener
	addr := net.UDPAddr{
		Port: 20777,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()

	log.Println("Listening for F1 telemetry on UDP 20777...")

	buf := make([]byte, 2048)
	packetCount := 0

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading UDP: %v", err)
			continue
		}

		// Parse the packet
		pkt, err := parseCarTelemetryPacket(buf[:n])
		if err != nil {
			log.Printf("Error parsing packet: %v", err)
			continue
		}

		// Convert to JSON
		jsonData, err := json.Marshal(pkt)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			continue
		}

		// Send to Kafka
		msg := &sarama.ProducerMessage{
			Topic: "devices",
			Value: sarama.ByteEncoder(jsonData),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %v", err)
		} else {
			packetCount++
			// Log every 60 packets (~1 second at 60Hz)
			if packetCount%60 == 0 {
				log.Printf("Sent packet #%d to Kafka [partition=%d, offset=%d] | Frame: %d, SessionTime: %.2fs",
					packetCount, partition, offset, pkt.Header.FrameIdentifier, pkt.Header.SessionTime)
			}
		}
	}
}

func parseCarTelemetryPacket(data []byte) (*telemetry.PacketCarTelemetryData, error) {
	if len(data) < 29 {
		return nil, fmt.Errorf("packet too small: %d bytes", len(data))
	}

	pkt := &telemetry.PacketCarTelemetryData{}
	offset := 0

	// Parse header (29 bytes)
	pkt.Header.PacketFormat = binary.LittleEndian.Uint16(data[offset:])
	offset += 2
	pkt.Header.GameYear = data[offset]
	offset++
	pkt.Header.GameMajorVersion = data[offset]
	offset++
	pkt.Header.GameMinorVersion = data[offset]
	offset++
	pkt.Header.PacketVersion = data[offset]
	offset++
	pkt.Header.PacketId = data[offset]
	offset++
	pkt.Header.SessionUID = binary.LittleEndian.Uint64(data[offset:])
	offset += 8
	pkt.Header.SessionTime = math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
	offset += 4
	pkt.Header.FrameIdentifier = binary.LittleEndian.Uint32(data[offset:])
	offset += 4
	pkt.Header.OverallFrameIdentifier = binary.LittleEndian.Uint32(data[offset:])
	offset += 4
	pkt.Header.PlayerCarIndex = data[offset]
	offset++
	pkt.Header.SecondaryPlayerCarIndex = data[offset]
	offset++

	// Parse 22 car telemetry data
	for i := 0; i < 22; i++ {
		if offset+60 > len(data) {
			break // Not enough data for this car
		}

		car := &pkt.CarTelemetryData[i]

		car.Speed = binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		car.Throttle = math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
		offset += 4

		car.Steer = math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
		offset += 4

		car.Brake = math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
		offset += 4

		car.Clutch = data[offset]
		offset++

		car.Gear = int8(data[offset])
		offset++

		car.EngineRPM = binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		car.DRS = data[offset]
		offset++

		car.RevLightsPercent = data[offset]
		offset++

		car.RevLightsBitValue = binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		// Brakes temperature (4 wheels)
		for j := 0; j < 4; j++ {
			car.BrakesTemperature[j] = binary.LittleEndian.Uint16(data[offset:])
			offset += 2
		}

		// Tyres surface temperature (4 wheels)
		for j := 0; j < 4; j++ {
			car.TyresSurfaceTemperature[j] = data[offset]
			offset++
		}

		// Tyres inner temperature (4 wheels)
		for j := 0; j < 4; j++ {
			car.TyresInnerTemperature[j] = data[offset]
			offset++
		}

		// Engine temperature
		car.EngineTemperature = binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		// Tyres pressure (4 wheels)
		for j := 0; j < 4; j++ {
			car.TyresPressure[j] = math.Float32frombits(binary.LittleEndian.Uint32(data[offset:]))
			offset += 4
		}

		// Surface type (4 wheels)
		for j := 0; j < 4; j++ {
			car.SurfaceType[j] = data[offset]
			offset++
		}
	}

	// Parse footer fields
	if offset+3 <= len(data) {
		pkt.MFDPanelIndex = data[offset]
		offset++
		pkt.MFDPanelIndexSecondary = data[offset]
		offset++
		pkt.SuggestedGear = int8(data[offset])
	}

	return pkt, nil
}
