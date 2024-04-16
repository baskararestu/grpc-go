# Import environment variables from .env file
include .env

# Protobuf compilation target
protoc:
	protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto

# Drop collection target
drop-collection:
	@echo "Dropping collection $(COLLECTION) from database $(DATABASE)..."
	@mongosh $(MONGODB_URI)/$(DATABASE_NAME) --eval "use $(DATABASE_NAME); db.$(COLLECTION).drop()"
	@echo "Collection $(COLLECTION) dropped from database $(DATABASE)."


# Help target
help:
	@echo "Available targets:"
	@echo "  make protoc           - Compile Protobuf files"
	@echo "  make drop-collection - Drop MongoDB collection"
	@echo "  make help             - Show this help message"

# Default target
.DEFAULT_GOAL := help
