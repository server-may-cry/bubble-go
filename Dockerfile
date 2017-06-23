FROM golang:1.8-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go get -u github.com/kardianos/govendor
