FRONTEND_BINARY=frontApp
DATA_BINARY=dataApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_data build_frontend
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_data: builds the data binary as a linux executable
build_data:
	@echo "Building data binary..."
	cd ../data-service && env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ${DATA_BINARY} ./cmd/api
	@echo "Done!"

## build_frontend: builds the frontend binary as a linux executable
build_frontend:
	@echo "Building frontend binary..."
	cd ../front-end && env GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ${FRONTEND_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_frontend
	@echo "Starting front end"
	cd ../front-end && ./${FRONTEND_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONTEND_BINARY}"
	@echo "Stopped front end!"
