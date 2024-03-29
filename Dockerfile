FROM golang:1.17

WORKDIR /go/src/github.com/motoki317/bot-extreme
COPY . .

RUN go mod download
RUN go build -o app

CMD ["./app"]
