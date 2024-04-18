FROM golang:1.22-alpine as builder
WORKDIR /build
COPY go.mod .
COPY src src
RUN go mod tidy
RUN go build -o child ./src

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/child .
CMD ["./child"]