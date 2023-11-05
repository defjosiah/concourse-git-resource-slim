# Start from a golang base image
FROM golang:1.18 as builder

# Set the Current Working Directory inside the container
WORKDIR /go/src/your-package-name

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

# Build the binaries using go build. This will generate the `check`, `in`, and `out` binaries.
RUN CGO_ENABLED=0 GOOS=linux go build -o /opt/resource/check ./cmd/check/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /opt/resource/in ./cmd/in/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /opt/resource/out ./cmd/out/main.go

# Use a small image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /opt/resource/check /opt/resource/check
COPY --from=builder /opt/resource/in /opt/resource/in
COPY --from=builder /opt/resource/out /opt/resource/out

# At runtime, nothing needs to be executed, the container just needs to exist
# for Concourse to run the binaries
CMD ["tail", "-f", "/dev/null"]
