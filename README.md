# F1 Kafka Pipeline

A data pipeline that fetches F1 telemetry and session data from the [OpenF1 API](https://api.openf1.org) and publishes it to Kafka topics — running on a GKE cluster provisioned with Terraform.

## Architecture

```
OpenF1 API ──► f1sim producer ──► Confluent Cloud / local Kafka ──► downstream consumers
                                          ▲
                              f1sim simulator (UDP telemetry)
                              f1sim listener (UDP → Kafka)
```

## Repo layout

```
.
├── f1sim/                  # Go application
│   ├── cmd/
│   │   ├── producer/       # Fetches OpenF1 data → Kafka
│   │   ├── listener/       # Receives UDP telemetry → Kafka
│   │   └── simulator/      # Generates fake F1 UDP telemetry at 60Hz
│   ├── internal/
│   │   ├── api/            # OpenF1 API client
│   │   ├── kafka/          # Sarama producer, topic admin, publishers
│   │   └── telemetry/      # UDP packet types
│   └── config/             # YAML config loaders
├── infra/                  # Terraform — GCP VPC + GKE cluster
├── Docker/                 # Dockerfiles + docker-compose
├── Makefile                # Common dev commands
└── Learnings/              # Notes and diagrams
```

## Kafka topics

| Topic | Partitions | Retention | Description |
|---|---|---|---|
| `f1-sessions` | 3 | 30 days | Session metadata |
| `f1-meetings` | 3 | 30 days | Race weekend meetings |
| `f1-cardata` | 12 | 7 days | Car telemetry per session |
| `f1-locations` | 12 | 7 days | GPS location data |
| `f1-intervals` | 12 | 7 days | Gap to leader/ahead |
| `f1-drivers` | 6 | 14 days | Driver info |
| `f1-laps` | 6 | 14 days | Lap times |
| `f1-positions` | 6 | 14 days | Race positions |
| `f1-stints` | 6 | 14 days | Tyre stint data |
| `f1-pits` | 3 | 14 days | Pit stop events |
| `f1-racecontrol` | 3 | 14 days | Flags, safety car, messages |
| `f1-results` | 3 | 30 days | Session results |
| `f1-grid` | 3 | 30 days | Starting grid |
| `f1-radio` | 6 | 14 days | Team radio |
| `f1-weather` | 3 | 14 days | Track weather |
| `devices` | 6 | 3 days | Raw UDP telemetry packets |

## Prerequisites

- Go 1.22+
- Docker + Docker Compose
- `gcloud` CLI (for GKE)
- Terraform (for infra)

## Credentials

SASL credentials for Confluent Cloud are loaded from environment variables — never committed.

```bash
cp .env.example .env
# fill in your Confluent Cloud API key and secret
```

Or export directly:

```bash
export KAFKA_SASL_USERNAME=your-api-key
export KAFKA_SASL_PASSWORD=your-api-secret
```

## Quick start (local Kafka via Docker)

```bash
# Start Kafka + Zookeeper + Kafka UI
make up

# Run the OpenF1 producer (fetches real data → local Kafka)
make run-producer   # or: cd f1sim && go run cmd/producer/main.go

# Run the UDP telemetry simulator
make run-simulator

# Run the UDP listener (receives simulator packets → Kafka)
make run-listener

# Open Kafka UI
make kafka-ui       # http://localhost:8080
```

## Producer config

Edit [f1sim/cmd/producer/app_config.yaml](f1sim/cmd/producer/app_config.yaml) to control what data is fetched:

```yaml
years: "2024"       # comma-separated or single year
all_data: true      # publish all 13 data types
driver: ""          # filter by driver number, or leave empty for all
base_url: "https://api.openf1.org/v1"
```

Kafka connection is configured in [f1sim/cmd/producer/kafka_config.yaml](f1sim/cmd/producer/kafka_config.yaml) — credentials come from env vars.

## Make targets

```
make build          Build all Docker images
make up             Start all services (detached)
make down           Stop all services
make rebuild        Clean rebuild from scratch
make logs           Tail logs from all containers
make logs-producer  Tail producer logs only
make kafka-topics   List Kafka topics
make kafka-consume  Consume from 'devices' topic
make kafka-ui       Open Kafka UI in browser
make run-simulator  Run simulator locally (UDP → port 20777)
make run-listener   Run listener locally (UDP → Kafka)
```

## Infra

GKE cluster on GCP, provisioned with Terraform. See [infra/readme.md](infra/readme.md).

```bash
cd infra
terraform init
terraform apply --auto-approve

# Get cluster credentials
gcloud container clusters get-credentials kafka-f1-cluster --zone us-central1-a --project pca-prep-441322
```
