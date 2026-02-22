FROM golang:1.26.0-alpine3.23 AS builder

ENV APP=/data/app
WORKDIR $APP

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/server/server.go


FROM alpine:3.21

WORKDIR /app

COPY --from=builder /data/app/app .

EXPOSE 8000

CMD ["./app", "serve"]
