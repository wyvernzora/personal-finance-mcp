# syntax=docker/dockerfile:1

# Stage 1: Build the Go server
FROM golang:1.24-alpine AS gobuilder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary packages
RUN apk update --no-cache && apk upgrade --no-cache \
    && apk add --no-cache git tzdata ca-certificates

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go server for the target platform
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/bin/mcp-server ./cmd/mcp-server

# Stage 2: Create a distroless image
FROM gcr.io/distroless/static-debian11 AS result

# Copy CA certificates
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the built application to the resulting container
COPY --from=gobuilder /app/bin/mcp-server /mcp-server

# Switch to an unprivileged user
USER nonroot:nonroot

# Setup the bind address and port for the app
ENV BIND_ADDRESS=0.0.0.0
ENV PORT=3000
EXPOSE 3000

# Set the default entry point to the application binary
ENTRYPOINT ["/mcp-server"]

# Set the default command to run the API server
CMD []
