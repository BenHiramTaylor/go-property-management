FROM golang:1.17.1

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...

CMD ["go","run","main.go"]