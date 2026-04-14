FROM golang:1.25-alpine

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /migrations
COPY *.sql .

CMD ["goose", "-dir", "/migrations", "postgres", "postgres://postgres:postgres@postgres:5432/project_tracker?sslmode=disable", "up"]
