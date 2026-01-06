FROM golang:alpine AS builder
WORKDIR /app
# Install build dependencies
RUN apk add --no-cache gcc musl-dev vips-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Enable CGO for libvips binding
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /analyzer ./cmd/api/

FROM alpine:latest
# Install runtime dependencies
RUN apk add --no-cache vips
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
COPY --from=builder /analyzer /analyzer
USER nonroot:nonroot
EXPOSE 8080

ENV GIN_MODE=release

ENTRYPOINT ["/analyzer"]