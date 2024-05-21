FROM golang:1.22.3-alpine3.19

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

CMD ["/app/main"]
