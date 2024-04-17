# Import environment variables from .env file
include .env

# Protobuf compilation target
protoc:
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto

.PHONY: start-mongodb
## start-mongodb: starts mongodb instance used for the app
start-mongodb:
	@ docker-compose up -d mongodb
	@ echo "Waiting for MongoDB to start..."
	@ until docker exec $(MONGODB_DATABASE_CONTAINER_NAME) mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1; do \
		echo "MongoDB not ready, sleeping for 5 seconds..."; \
		sleep 5; \
	done
	@ echo "MongoDB is up and running."

.PHONY: stop-mongodb
## stop-mongodb: stops mongodb instance used for the app
stop-mongodb:
	@ docker-compose stop mongodb

.PHONY: start-test-mongodb
## start-test-mongodb: starts mongodb instance used for integration tests
start-test-mongodb:
	@ docker-compose up -d mongodb_test
	@ echo "Waiting for Test MongoDB to start..."
	@ until docker exec $(MONGODB_TEST_DATABASE_CONTAINER_NAME) mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1; do \
		echo "Test MongoDB not ready, sleeping for 5 seconds..."; \
		sleep 5; \
	done
	@ echo "Test MongoDB is up and running."

.PHONY: stop-test-mongodb
## stop-test-mongodb: stops mongodb instance used for integration tests
stop-test-mongodb:
	@ docker-compose stop mongodb_test
	@ docker rm ${MONGODB_TEST_DATABASE_CONTAINER_NAME}

.PHONY: test
## test: runs both unit and integration tests
test: test
	@ go test -v ./...

# Help target
help:
	@echo "Available targets:"
	@echo "  make protoc           - Compile Protobuf files"
	@echo "  make drop-collection - Drop MongoDB collection"
	@echo "  make help             - Show this help message"

.PHONY: run
## run: runs the gRPC server
run: start-mongodb
	@ go run cmd/main.go

# Default target
.DEFAULT_GOAL := help
