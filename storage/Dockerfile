FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

COPY ./storage/go.mod ./

COPY ./storage/go.sum ./

COPY ./storage/example.env ./

RUN mkdir "cmd"

RUN mkdir "internal"

COPY ./storage/cmd ./cmd

COPY ./storage/internal ./internal

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o storage_service ./cmd/main/main.go

FROM alpine:3.19.0 AS runner

WORKDIR /root/

COPY --from=builder /app/storage_service .
COPY --from=builder /app/example.env .

CMD ["./storage_service"]
