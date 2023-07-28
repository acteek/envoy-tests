FROM golang:1.20 AS builder

WORKDIR /src
COPY go.mod .
COPY go.sum .
COPY ./ ./
RUN go mod download
RUN go build -ldflags '-w -s' -o app .
CMD ["./app"]


