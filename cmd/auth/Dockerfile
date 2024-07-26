FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o auth-service cmd/auth/*.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth-service .

RUN chmod +x /root/auth-service

RUN ls -l /root/

CMD ["./auth-service"]
