FROM golang:1.10-alpine AS build-env
WORKDIR /go/src/github.com/server-may-cry/bubble-go
ADD . .
ADD .netrc /root/.netrc
RUN apk update && apk upgrade && \
    apk add --no-cache git make && \
    go get -u golang.org/x/vgo && \
    make build

FROM alpine
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=build-env /go/src/github.com/server-may-cry/bubble-go/bubble-go /app/server
COPY --from=build-env /go/src/github.com/server-may-cry/bubble-go/config /app/config
CMD ["/app/server"]
