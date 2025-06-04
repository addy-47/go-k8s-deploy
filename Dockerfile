FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod .

RUN go mod tidy

COPY . .

RUN go build -o go-test-app

EXPOSE 8080

CMD ["./go-test-app"]