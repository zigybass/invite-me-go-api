FROM golang:1.15.0-alpine

WORKDIR /app

COPY . /app

RUN go build -o server .

CMD ["/app/server"]