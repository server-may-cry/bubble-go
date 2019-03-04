FROM alpine

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

ADD https://119226.selcdn.ru/bubble/bubble_all.zip /tmp/bubble_all.zip
RUN mkdir -p /app/static/bubble && unzip /tmp/bubble_all.zip -d /app/static/bubble

COPY ./bubble-go /app/bubble-go
COPY ./config /app/config

CMD ["./bubble-go"]
