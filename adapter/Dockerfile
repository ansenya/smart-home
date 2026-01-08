FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /build

FROM alpine:latest

RUN apk add curl

COPY --from=builder /build /app/build

WORKDIR /app

EXPOSE 8080

CMD ["/app/build"]
