LASTFM_BINARY=lastFmApp


## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## build_broker: builds the broker binary as a linux executable
build_lastfm:
	@echo "Building lastfm binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${LASTFM_BINARY} ./
	@echo "Done!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_lastfm
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"