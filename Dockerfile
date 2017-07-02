FROM golang:1.8-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache git openssl make && \
    go get -u github.com/golang/dep/cmd/dep
