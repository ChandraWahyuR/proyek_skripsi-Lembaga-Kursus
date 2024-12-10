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
COPY static /static
COPY docs /docs

 # Salin file .env ke dalam container
COPY Keys_GCS-Image.json Keys_GCS-Image.json  
# Salin file credentials JSON

# Assets html
COPY --from=build-stage /app/assets /assets
# Expose the application port
EXPOSE 8080

# Start the application
CMD [ "./goapp" ]
