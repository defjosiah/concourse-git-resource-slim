.PHONY: build-check build-in build-out docker-build run-check run-in run-out

# Variables
BASE_DIR=./cmd
DOCKER_IMAGE_NAME=defjosiah/concourse-git-resource-slim

# Build the check binary
build-check:
	@echo "Building check binary..."
	cd $(BASE_DIR)/check && go build -o check

# Build the in binary
build-in:
	@echo "Building in binary..."
	cd $(BASE_DIR)/in && go build -o in

# Build the out binary
build-out:
	@echo "Building out binary..."
	cd $(BASE_DIR)/out && go build -o out

# Build the Docker image
docker-build: build-check build-in build-out
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .

# Run the check command
run-check:
	@echo "Running check command..."
	go run $(BASE_DIR)/check/main.go

# Run the in command
run-in:
	@echo "Running in command..."
	go run $(BASE_DIR)/in/main.go

# Run the out command
run-out:
	@echo "Running out command..."
	go run $(BASE_DIR)/out/main.go

# Run this if you want to build all binaries and the docker image
all: docker-build
