# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build service binary (passed via build arg)
ARG SERVICE
RUN go build -o /${SERVICE} ./${SERVICE}

# ---- Runtime stage ----
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
ARG SERVICE
COPY --from=builder /${SERVICE} /app/${SERVICE}

# Run service
CMD ["sh", "-c", "./${SERVICE}"]
