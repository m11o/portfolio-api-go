FROM golang:latest

WORKDIR /go/src/app

COPY . /go/src/app

EXPOSE 3000
