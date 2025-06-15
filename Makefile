BUILD_DIR := ./bin
API_NAME := api

# Apps
gen:
	sqlc generate

build: gen
	@echo "Building api..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(API_NAME) ./cmd/api

api: build
	GOMANIA_CONNECTION_STRING="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" $(BUILD_DIR)/$(API_NAME) $(ARGS)

debug: build
	GOMANIA_CONNECTION_STRING="postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable" dlv exec $(BUILD_DIR)/$(API_NAME) --listen=:2345 --headless=true --api-version=2 --accept-multiclient -- $(ARGS)

# Docker
docker-up:
	@echo "Building docker image..."
	docker compose up -d

docker-down:
	@echo "Stopping docker containers..."
	docker compose down

# database
db-up:
	@echo "Running migrations..."
	DATABASE_URL="postgres://postgres@127.0.0.1:5430/postgres?sslmode=disable" dbmate -d ./data/sql/migrations up

db-drop:
	@echo "Dropping database..."
	DATABASE_URL="postgres://postgres@127.0.0.1:5430/postgres?sslmode=disable" dbmate -d ./data/sql/migrations drop

db-seed:
	@echo "seeding database..."
	GOMANIA_CONNECTION_STRING="postgres://postgres@127.0.0.1:5430/postgres?sslmode=disable" go run ./cmd/tools/seed/main.go

# Testing
test-api:
	@echo "Testing API endpoints..."
	./scripts/test_api.sh

setup: docker-up db-up db-seed
	@echo "Complete setup finished! Server ready to start."

.PHONY: gen build api debug docker-up docker-down db-init db-migrate db-seed db-reset db-up db-drop test-api setup dev
