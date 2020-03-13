FROM golang:1.14

WORKDIR /go/src/github.com/motoki317/bot-extreme
COPY . .

RUN go mod install
RUN go build -o app

CMD ["app"]
