FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

ENV GOPROXY=https://goproxy.io,direct


RUN go mod download

COPY . ./

WORKDIR /app/src/cmd

RUN go build -v -o /app/server

WORKDIR /app

CMD ["./server"]
