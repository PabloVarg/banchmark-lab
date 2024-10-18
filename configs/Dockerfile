FROM golang:1.23 AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN mkdir ./bin
RUN go build -o /usr/bin/client ./cmd/client
RUN go build -o /usr/bin/server ./cmd/server

ENTRYPOINT [ "server" ]
