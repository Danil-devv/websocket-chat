FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

COPY .env ./

RUN mkdir "cmd"

RUN mkdir "internal"

COPY cmd ./cmd

COPY internal ./internal

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o webserver ./cmd/main/main.go

FROM alpine:3.19.0 AS runner

WORKDIR /root/

COPY --from=builder /app/webserver .
COPY --from=builder /app/.env .

CMD ["./webserver"]
