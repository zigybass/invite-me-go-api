FROM golang:1.15.0-alpine

WORKDIR /zigy

COPY . /zigy

RUN go mod download

COPY . /zigy

RUN go build /zigy/cmd/server/main.go

CMD ["./main"]
