package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"f1sim/internal/telemetry"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:20777")
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial UDP: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down simulator...")
		cancel()
	}()

	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano())))

	var frameIdentifier uint32
	sessionUID := uint64(time.Now().Unix())

	fmt.Printf("F1 Telemetry Simulator Started\n")
	fmt.Printf("UDP Target: 127.0.0.1:20777\n")
	fmt.Printf("Session UID: %d\n", sessionUID)
	fmt.Printf("Frequency: 60Hz\n")
	fmt.Printf("Press Ctrl+C to stop\n\n")

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		frameIdentifier++
		pkt := generateTelemetryPacket(rng, sessionUID, frameIdentifier)

		buf := make([]byte, 1500)
		offset := 0

		// Serialize header (29 bytes total)
		binary.LittleEndian.PutUint16(buf[offset:], pkt.Header.PacketFormat)
		offset += 2
		buf[offset] = pkt.Header.GameYear
		offset++
		buf[offset] = pkt.Header.GameMajorVersion
		offset++
		buf[offset] = pkt.Header.GameMinorVersion
		offset++
		buf[offset] = pkt.Header.PacketVersion
		offset++
		buf[offset] = pkt.Header.PacketId
		offset++
		binary.LittleEndian.PutUint64(buf[offset:], pkt.Header.SessionUID)
		offset += 8
		binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(pkt.Header.SessionTime))
		offset += 4
		binary.LittleEndian.PutUint32(buf[offset:], pkt.Header.FrameIdentifier)
		offset += 4
		binary.LittleEndian.PutUint32(buf[offset:], pkt.Header.OverallFrameIdentifier)
		offset += 4
		buf[offset] = pkt.Header.PlayerCarIndex
		offset++
		buf[offset] = pkt.Header.SecondaryPlayerCarIndex
		offset++

		// Serialize all 22 car telemetry data
		for i := 0; i < 22; i++ {
			car := pkt.CarTelemetryData[i]

			binary.LittleEndian.PutUint16(buf[offset:], car.Speed)
			offset += 2

			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(car.Throttle))
			offset += 4

			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(car.Steer))
			offset += 4

			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(car.Brake))
			offset += 4

			buf[offset] = car.Clutch
			offset++

			buf[offset] = byte(car.Gear)
			offset++

			binary.LittleEndian.PutUint16(buf[offset:], car.EngineRPM)
			offset += 2

			buf[offset] = car.DRS
			offset++

			buf[offset] = car.RevLightsPercent
			offset++

			binary.LittleEndian.PutUint16(buf[offset:], car.RevLightsBitValue)
			offset += 2

			// Brakes temperature (4 wheels)
			for j := 0; j < 4; j++ {
				binary.LittleEndian.PutUint16(buf[offset:], car.BrakesTemperature[j])
				offset += 2
			}

			// Tyres surface temperature (4 wheels)
			for j := 0; j < 4; j++ {
				buf[offset] = car.TyresSurfaceTemperature[j]
				offset++
			}

			// Tyres inner temperature (4 wheels)
			for j := 0; j < 4; j++ {
				buf[offset] = car.TyresInnerTemperature[j]
				offset++
			}

			// Engine temperature
			binary.LittleEndian.PutUint16(buf[offset:], car.EngineTemperature)
			offset += 2

			// Tyres pressure (4 wheels)
			for j := 0; j < 4; j++ {
				binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(car.TyresPressure[j]))
				offset += 4
			}

			// Surface type (4 wheels)
			for j := 0; j < 4; j++ {
				buf[offset] = car.SurfaceType[j]
				offset++
			}
		}

		// Serialize packet footer fields
		buf[offset] = pkt.MFDPanelIndex
		offset++
		buf[offset] = pkt.MFDPanelIndexSecondary
		offset++
		buf[offset] = byte(pkt.SuggestedGear)
		offset++

		// send live packet
		conn.Write(buf[:offset])

		// Print telemetry every 60 frames (~1 second)
		if frameIdentifier%60 == 0 {
			fmt.Printf("[Frame %d] Session Time: %.2fs | ", frameIdentifier, pkt.Header.SessionTime)
			// Show telemetry for car 0 (player car)
			car := pkt.CarTelemetryData[0]
			fmt.Printf("Speed: %d km/h | Gear: %d | RPM: %d | Throttle: %.0f%% | Brake: %.0f%%\n",
				car.Speed, car.Gear, car.EngineRPM, car.Throttle*100, car.Brake*100)
		}

		time.Sleep(16 * time.Millisecond) // ~60Hz
	}
}

func generateTelemetryPacket(rng *rand.Rand, sessionUID uint64, frameIdentifier uint32) telemetry.PacketCarTelemetryData {
	pkt := telemetry.PacketCarTelemetryData{
		Header: telemetry.PacketHeader{
			PacketFormat:            2024,
			GameYear:                24,
			GameMajorVersion:        1,
			GameMinorVersion:        0,
			PacketVersion:           1,
			PacketId:                6, // Car Telemetry
			SessionUID:              sessionUID,
			SessionTime:             float32(frameIdentifier) / 60.0, // Assuming 60 Hz
			FrameIdentifier:         frameIdentifier,
			OverallFrameIdentifier:  uint32(frameIdentifier),
			PlayerCarIndex:          0,
			SecondaryPlayerCarIndex: 255, // No secondary player
		},
		MFDPanelIndex:          0,
		MFDPanelIndexSecondary: 255,
		SuggestedGear:          4,
	}

	// Generate telemetry data for all 22 cars
	for i := 0; i < 22; i++ {
		// More realistic values for active cars (first 10) vs inactive cars
		isActiveCar := i < 10

		speed := uint16(0)
		throttle := float32(0)
		brake := float32(0)
		gear := int8(-1)
		engineRPM := uint16(1000)

		if isActiveCar {
			speed = uint16(150 + rng.IntN(80))
			throttle = rng.Float32()
			brake = rng.Float32() * 0.5
			gear = int8(1 + rng.IntN(7))
			engineRPM = uint16(9000 + rng.IntN(4000))
		}

		pkt.CarTelemetryData[i] = telemetry.CarTelemetryData{
			Speed:                   speed,
			Throttle:                throttle,
			Steer:                   rng.Float32()*2 - 1, // -1 to 1
			Brake:                   brake,
			Clutch:                  0,
			Gear:                    gear,
			EngineRPM:               engineRPM,
			DRS:                     0,
			RevLightsPercent:        uint8(rng.IntN(100)),
			RevLightsBitValue:       uint16(rng.IntN(65536)),
			BrakesTemperature:       [4]uint16{uint16(200 + rng.IntN(400)), uint16(200 + rng.IntN(400)), uint16(200 + rng.IntN(400)), uint16(200 + rng.IntN(400))},
			TyresSurfaceTemperature: [4]uint8{uint8(60 + rng.IntN(60)), uint8(60 + rng.IntN(60)), uint8(60 + rng.IntN(60)), uint8(60 + rng.IntN(60))},
			TyresInnerTemperature:   [4]uint8{uint8(70 + rng.IntN(50)), uint8(70 + rng.IntN(50)), uint8(70 + rng.IntN(50)), uint8(70 + rng.IntN(50))},
			EngineTemperature:       uint16(80 + rng.IntN(40)),
			TyresPressure:           [4]float32{20.0 + rng.Float32()*5, 20.0 + rng.Float32()*5, 20.0 + rng.Float32()*5, 20.0 + rng.Float32()*5},
			SurfaceType:             [4]uint8{0, 0, 0, 0}, // 0 = Tarmac
		}
	}

	return pkt
}
