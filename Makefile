.PHONY: build up down restart logs clean help build-simulator run-simulator run-listener

# Variables
DOCKER_COMPOSE = docker-compose -f Docker/docker-compose.yaml
IMAGE_NAME = f1-simulator
CONTAINER_NAME = f1-simulator
LISTENER_CONTAINER = f1-listener

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build-simulator: ## Build only the simulator Docker image
	docker build -t $(IMAGE_NAME):latest -f Docker/Dockerfile .

build: ## Build all Docker images via docker-compose
	$(DOCKER_COMPOSE) build

up: ## Start all services in detached mode
	$(DOCKER_COMPOSE) up -d

down: ## Stop and remove all containers
	$(DOCKER_COMPOSE) down

restart: down up ## Restart all services

logs: ## Show logs from all containers
	$(DOCKER_COMPOSE) logs -f

logs-app: ## Show logs from the Go app only
	$(DOCKER_COMPOSE) logs -f app

logs-kafka: ## Show logs from Kafka
	$(DOCKER_COMPOSE) logs -f kafka

logs-listener: ## Show logs from the listener service
	$(DOCKER_COMPOSE) logs -f listener

logs-producer: ## Show logs from the producer service
	$(DOCKER_COMPOSE) logs -f producer

logs-kafka-ui:
	$(DOCKER_COMPOSE) logs -f kafka-ui

ps: ## List running containers
	$(DOCKER_COMPOSE) ps

clean: ## Stop containers and remove volumes
	$(DOCKER_COMPOSE) down -v

rebuild: clean build up ## Clean rebuild and restart everything

shell-app: ## Open shell in the app container
	docker exec -it $(CONTAINER_NAME) sh

shell-listener: ## Open shell in the listener container
	docker exec -it $(LISTENER_CONTAINER) sh

run-simulator: ## Run F1 simulator locally (sends UDP to port 20777)
	cd f1sim && go run cmd/simulator/main.go

run-listener: ## Run listener locally (receives UDP and sends to Kafka)
	cd f1sim && KAFKA_BROKER=localhost:9093 go run cmd/listener/main.go

kafka-topics: ## List Kafka topics
	docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --list

kafka-describe: ## Describe the 'devices' topic
	docker exec -it kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic devices

kafka-consume: ## Consume messages from 'devices' topic
	docker exec -it kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic devices --from-beginning

kafka-ui: ## Open Kafka UI (http://localhost:8080)
	@echo "Kafka UI is available at: http://localhost:8080"
	@open http://localhost:8080 2>/dev/null || xdg-open http://localhost:8080 2>/dev/null || echo "Please open http://localhost:8080 in your browser"

dev: ## Start services and follow app logs
	$(DOCKER_COMPOSE) up --build

test: ## Run tests inside the container
	docker exec -it $(CONTAINER_NAME) go test ./...
