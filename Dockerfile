FROM alpine
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates && update-ca-certificates
COPY ./bubble-go /app/bubble-go
COPY ./config /app/config
CMD ["/app/bubble-go"]
