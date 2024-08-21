FROM golang:alpine
LABEL authors="cheglockvvv"

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o app cmd/main.go

EXPOSE 8080

CMD ["./app"]

ENTRYPOINT ["top", "-b"]