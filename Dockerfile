FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/swagger-ui ./swagger-ui

RUN mkdir logs

EXPOSE 8080

CMD ["./main"]