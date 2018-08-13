FROM alpine
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR /app
COPY ./bubble-go /app/bubble-go
COPY ./config /app/config
COPY ./version /app/version
CMD ["./bubble-go"]
