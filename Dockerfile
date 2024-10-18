# Build stage
FROM golang:1.22.2 AS build-stage

WORKDIR /app

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /goapp

# Release stage
FROM gcr.io/distroless/base-debian11 AS release-stage

WORKDIR /

# Copy executable and environment file
COPY --from=build-stage /goapp /goapp
COPY .env .env

# Expose the application port
EXPOSE 8080

# Run as non-root user
USER nonroot:nonroot

# Start the application
CMD [ "./goapp" ]
