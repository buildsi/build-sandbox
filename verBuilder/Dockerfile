# Start from the latest golang base image
FROM ghcr.io/autamus/go:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o verBuilder .

# Start again with minimal envoirnment.
FROM spack/ubuntu-bionic:latest

RUN apt-get update && \
    apt-get install -y ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

COPY --from=builder /app/verBuilder /app/verBuilder

ENV PATH=/opt/spack/bin:$PATH

# Command to run the executable
ENTRYPOINT ["/app/verBuilder"]