# ---- Build stage ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=1 go build -o app ./cmd/server

# ---- Run stage ----
FROM alpine:3.20

RUN apk add --no-cache sqlite-libs

WORKDIR /app

# Create data directory for SQLite
RUN mkdir -p /app/data /app/storage

# Copy binary from builder
COPY --from=builder /app/app .


# Expose port
EXPOSE 8080

# Environment variables (default)
ENV PORT=8080
ENV DB_URL=/app/storage/visa.db
# Run app
CMD ["./app"]
