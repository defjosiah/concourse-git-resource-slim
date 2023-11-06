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
	docker build -t $(DOCKER_IMAGE_NAME):latest .

docker-push: docker-build
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE_NAME):latest

# Run the check command
run-check:
	@go run $(BASE_DIR)/check/main.go

# Run the in command
run-in:
	@go run $(BASE_DIR)/in/main.go ./tmp/in

# Run the out command
run-out:
	@go run $(BASE_DIR)/out/main.go

# Run this if you want to build all binaries and the docker image
all: docker-build
