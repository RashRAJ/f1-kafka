# F1 2024 Telemetry Simulator

A UDP-based simulator that generates realistic F1 2024 telemetry packets for testing and development.

## Overview

This simulator sends F1 2024 telemetry data packets over UDP at 60Hz, mimicking the real F1 game telemetry output. It generates realistic data for all 22 cars on the grid.

## Project Structure

```
f1sim/
├── cmd/
│   ├── simulator/     # Telemetry packet generator
│   └── listener/      # UDP packet receiver
├── internal/
│   └── packets/       # F1 2024 packet structures
└── README.md
```

## Prerequisites

- Go 1.22 or later
- UDP port 20777 available

## Installation

```bash
cd f1sim
go build ./cmd/simulator
go build ./cmd/listener
```

## Usage

### Running the Simulator

Start the telemetry simulator to send packets to `127.0.0.1:20777`:

```bash
./simulator
```

**Output:**
```
F1 Telemetry Simulator Started
UDP Target: 127.0.0.1:20777
Session UID: 1763397300
Frequency: 60Hz
Press Ctrl+C to stop

[Frame 60] Session Time: 1.00s | Speed: 164 km/h | Gear: 2 | RPM: 10402 | Throttle: 34% | Brake: 42%
[Frame 120] Session Time: 2.00s | Speed: 180 km/h | Gear: 3 | RPM: 9945 | Throttle: 76% | Brake: 10%
```

The simulator displays telemetry for the player car (car index 0) every second.

### Running the Listener

In a separate terminal, start the listener to receive and display packets:

```bash
./listener
```

### Stopping

Press `Ctrl+C` in either terminal to gracefully shut down the simulator or listener.

## Telemetry Data

The simulator generates complete F1 2024 telemetry packets (packet ID: 6) with:

### Packet Header
- Packet format (2024)
- Game version
- Session UID
- Frame identifier
- Session time
- Player car index

### Car Telemetry (for all 22 cars)
- **Speed** (km/h)
- **Throttle** (0-100%)
- **Steering** (-1 to 1)
- **Brake** (0-100%)
- **Clutch**
- **Gear** (-1 to 7)
- **Engine RPM**
- **DRS status**
- **Rev lights**
- **Brake temperatures** (all 4 wheels)
- **Tyre surface temperatures** (all 4 wheels)
- **Tyre inner temperatures** (all 4 wheels)
- **Engine temperature**
- **Tyre pressures** (all 4 wheels)
- **Surface types** (all 4 wheels)

### Packet Footer
- MFD panel index
- Suggested gear

## Configuration

### Changing the UDP Target

Edit `cmd/simulator/main.go`:

```go
addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:20777")
```

Change the IP address and port as needed.

### Adjusting Transmission Rate

The default is 60Hz (16ms sleep). To change it, edit:

```go
time.Sleep(16 * time.Millisecond) // ~60Hz
```

## Data Characteristics

- **Active cars** (0-9): Generate realistic racing telemetry
  - Speed: 150-230 km/h
  - RPM: 9,000-13,000
  - Gears: 1-7
  - Variable throttle/brake

- **Inactive cars** (10-21): Generate minimal/stationary telemetry
  - Speed: 0 km/h
  - Gear: -1 (neutral)
  - RPM: 1,000 (idle)

## Use Cases

- Testing F1 telemetry consumers
- Kafka/streaming pipeline development
- Dashboard and visualization development
- Race analytics development
- Learning F1 2024 UDP protocol

## Packet Format

The simulator follows the official F1 2024 UDP specification. Total packet size is approximately 1,350 bytes:

- Header: 33 bytes
- Car telemetry: 22 cars × 60 bytes = 1,320 bytes
- Footer: 3 bytes

## License

MIT


