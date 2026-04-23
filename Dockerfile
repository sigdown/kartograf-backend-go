FROM golang:1.25 AS builder

WORKDIR /app

COPY go.sum go.mod ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/kartograf-api ./cmd/api

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/kartograf-api /app/kartograf-api
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/kartograf-api"]