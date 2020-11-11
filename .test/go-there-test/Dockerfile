FROM golang:1.15.3-buster AS build

WORKDIR /root

COPY ./go-there-test/ .

RUN go mod download

CMD CGO_ENABLED=0 go test -v ./...