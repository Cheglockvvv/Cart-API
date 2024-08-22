FROM golang:1.23-alpine AS builder
LABEL authors="cheglockvvv"

WORKDIR /usr/local/src

RUN apk --no-cache add bash

# dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# build
COPY . .
RUN go build -o ./bin/app app/cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /

CMD ["/app"]