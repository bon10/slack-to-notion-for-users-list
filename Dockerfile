FROM golang:1.16

RUN apt update && apt install git

WORKDIR /go/src/app
ADD . /go/src/app
