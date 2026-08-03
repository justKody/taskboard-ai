# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /taskboard-go-api ./cmd/main.go

# Runtime stage
FROM alpine:3.24

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& adduser -D -H -u 10001 appuser

COPY --from=builder /taskboard-go-api /app/taskboard-go-api

USER appuser

EXPOSE 8080

ENV PORT=8080

CMD ["/app/taskboard-go-api"]
