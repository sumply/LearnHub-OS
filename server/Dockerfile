FROM golang:1.26.0-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/server/server.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/app .

COPY /db/migrations ./migrations

EXPOSE 8000

ENTRYPOINT [ "./app" ]

CMD ["serve"]
