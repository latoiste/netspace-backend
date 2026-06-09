# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /netspace .

# ---- Run stage ----
FROM alpine:latest

# TLS roots (needed for outbound TLS, e.g. a managed Postgres over SSL)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /netspace /app/netspace

# The app listens on $PORT (Railway/Render inject it), falling back to 8080.
EXPOSE 8080

CMD ["/app/netspace"]
