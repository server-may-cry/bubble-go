FROM golang:1.10-alpine

RUN apk update && apk upgrade \
    && apk add --no-cache \
        # C compiler for cgo
        gcc \
        git \
        make \
        # C 'stdlib.h'
        musl-dev \
    && go get -u golang.org/x/vgo
