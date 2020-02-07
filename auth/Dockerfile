FROM golang:latest

WORKDIR /app

ADD . .

RUN make

ENTRYPOINT ["./auth"]
