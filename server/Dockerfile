FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

COPY ./server/go.mod ./

COPY ./server/go.sum ./

COPY ./server/.env ./

RUN mkdir "cmd"

RUN mkdir "internal"

COPY ./server/cmd ./cmd

COPY ./server/internal ./internal

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o webserver ./cmd/main/main.go

FROM alpine:3.19.0 AS runner

WORKDIR /root/

COPY --from=builder /app/webserver .
COPY --from=builder /app/.env .

CMD ["./webserver"]
