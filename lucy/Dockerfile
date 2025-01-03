# Step 1: Build stage
FROM golang:1.23-alpine AS builder

# Set environment variable
ENV PORT=3333
ENV CGO_ENABLED=1

# Set the working directory inside the container
WORKDIR /app

RUN apk --no-cache add gcc musl-dev sqlite-dev

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy all source files into the container
COPY . .

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN goose -dir=internal/sql/schema sqlite3 lucy_db.sqlite up

# Install sqlc and generate SQL code
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate

# Build the Go application
RUN go build -o lucy ./cmd

# Step 2: Create the final image (distroless or minimal base)
FROM alpine:latest

# Install necessary dependencies (CA certificates for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Set the working directory for the final image
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/lucy .
COPY --from=builder /app/lucy_db .

# Expose the port the application will run on
EXPOSE 3333

# Run the application when the container starts
CMD ["./lucy"]
