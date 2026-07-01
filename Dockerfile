# --- Stage 1: Build the Go Backend & Fetch Migration Tools ---
FROM golang:1.26.4-alpine3.23 AS builder
WORKDIR /src

# Install git/gcc if required by any Go modules
RUN apk add --no-cache git

# Install Goose for running migrations inside the container
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Cache and install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire Go source tree and SQL directory
COPY internal/ ./internal/
COPY server/ ./server/
COPY sql/ ./sql/

# Build the Go application from the server directory
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./server

# --- Stage 2: Final Production Container ---
FROM alpine:3.19
WORKDIR /workspace

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the compiled Go binary from Stage 1
COPY --from=builder /app/server ./server

# Copy Goose binary and SQL migrations from Stage 1
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /src/sql/ /workspace/sql/

# Copy the barebones frontend directly from the host filesystem
COPY app/ /workspace/app/

# Expose web server port
EXPOSE 8080

# Run migrations first, then start the Go server
CMD ["sh", "-c", "goose -dir ./sql/schema postgres \"$DB_URL\" up && ./server"]