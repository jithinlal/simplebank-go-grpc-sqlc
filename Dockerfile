# build step
FROM golang:1.20.12-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run step
FROM alpine:3.13
WORKDIR /app

COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./db/migrations

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
