FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o api-service cmd/api/*.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/api-service .
COPY --from=builder /app/.env .

RUN chmod +x /root/api-service

RUN ls -l /root/

CMD ["./api-service"]
