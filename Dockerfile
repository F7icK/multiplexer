# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -v -o /go/bin/share ./cmd/multiplexer

FROM alpine
WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=builder /go/bin/share .

EXPOSE 8081

ENTRYPOINT ["/app/share"]
