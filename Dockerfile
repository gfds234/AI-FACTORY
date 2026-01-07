FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orchestrator .

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget

WORKDIR /root/
COPY --from=builder /app/orchestrator .
COPY --from=builder /app/web ./web

# Persistent directories for Railway volume mount
RUN mkdir -p /data/projects /data/artifacts

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./orchestrator"]
