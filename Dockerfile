FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .env .

COPY . .

RUN make build

FROM ubuntu:latest

COPY --from=builder /app /app

COPY --from=builder /app/.env /app/.env

CMD ["/app/build/main"]
