FROM golang:1.23.2-alpine3.20 AS builder

RUN apk add --no-cache git

WORKDIR /usr/src/app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN go build -o ./bin/app ./cmd/main.go

FROM alpine:3.20 AS runner

COPY --from=builder /usr/src/app/bin/app /app

COPY .env .env

COPY migrations /migrations
COPY logs /logs


RUN mkdir -p /logs /migrations

# Запускаем приложение
<<<<<<< HEAD
CMD ["/app"]
=======
CMD ["/app"]
>>>>>>> 1bc302a1131fc1a4f3331f4ed2fd15212c2d44d7
