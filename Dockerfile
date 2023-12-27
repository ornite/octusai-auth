# Start by building the application.
FROM golang:latest as builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies.
COPY go.mod go.sum ./

# Download dependencies.
RUN go mod download

# Copy the rest of the application's code.
COPY . .

# Build the application. This assumes that the main package is in the root directory.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Now use a smaller image to run.
FROM alpine:latest  

# Install ca-certificates in case you need HTTPS.
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage.
COPY --from=builder /app/server .

# Copy the .env file with environment variables.
COPY --from=builder /app/.env .

# Expose the port the server listens on.
EXPOSE 50051

# Run the server binary.
CMD ["./server"]
