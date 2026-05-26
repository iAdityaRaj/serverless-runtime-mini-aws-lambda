
FROM golang:1.24-alpine

WORKDIR /app

COPY main.go .

RUN go build -o function main.go

CMD ["./function"]
