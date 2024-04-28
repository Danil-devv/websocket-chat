FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

COPY ./services/server/go.mod ./

COPY ./services/server/go.sum ./

COPY ./services/server/example.env ./

RUN mkdir "cmd"

RUN mkdir "internal"

COPY ./services/server/cmd ./cmd

COPY ./services/server/internal ./internal

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o chat_service ./cmd/main/main.go

FROM alpine:3.19.0 AS runner

WORKDIR /root/

COPY --from=builder /app/chat_service .
COPY --from=builder /app/example.env .

CMD ["./chat_service"]
