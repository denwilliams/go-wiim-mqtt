FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o wiim-mqtt ./cmd/wiim-mqtt


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/wiim-mqtt .
RUN chmod +x wiim-mqtt

ENTRYPOINT ["./wiim-mqtt"]
