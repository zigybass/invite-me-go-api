FROM golang:1.15.0-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o server .

CMD ["/app/server"]