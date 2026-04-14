FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /bin/server ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/server /usr/local/bin/server

EXPOSE 8080
CMD ["server"]
