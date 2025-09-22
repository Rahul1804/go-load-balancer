FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go build -o load-balancer ./cmd/load-balancer

EXPOSE 8080

CMD ["./load-balancer"]
