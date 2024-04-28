FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

COPY ./services/chat/go.mod ./

COPY ./services/chat/go.sum ./

COPY ./services/chat/example.env ./

RUN mkdir "cmd"

RUN mkdir "internal"

COPY ./services/chat/cmd ./cmd

COPY ./services/chat/internal ./internal

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o chat_service ./cmd/main/main.go

FROM alpine:3.19.0 AS runner

WORKDIR /root/

COPY --from=builder /app/chat_service .
COPY --from=builder /app/example.env .

CMD ["./chat_service"]
