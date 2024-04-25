FROM golang:1.22.1-alpine3.18

RUN mkdir build

# We create folder named build for our app.
WORKDIR /build

COPY go.mod go.sum ./

# Download dependencies
RUN go install github.com/cosmtrek/air@latest
RUN go mod download