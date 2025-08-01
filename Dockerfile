FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ratelimiter ./cmd/ratelimiter/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ratelimiter .

EXPOSE 50051

ENTRYPOINT [ "./ratelimiter", "--window=100s", "--limit=5" ]