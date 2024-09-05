FROM golang:1.22.5-alpine AS builder
LABEL authors="cheglockvvv"

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /app/xernia /app/cmd/main.go

FROM alpine AS runner

COPY --from=builder /app/xernia /cartApi/xernia
COPY --from=builder /app/internal/db/migrations /cartApi/migrations

EXPOSE ${API_PORT}

CMD ["/cartApi/xernia"]