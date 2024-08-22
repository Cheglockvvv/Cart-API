FROM golang:1.23-alpine
LABEL authors="cheglockvvv"

WORKDIR /usr/local/src

RUN apk --no-cache add bash

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build

